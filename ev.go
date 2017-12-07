package ev

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/evio"
)

type Ev struct {
	Server                         evio.Server
	Conn                           Conn
	Host                           string
	DataHandler                    DataHandlerFunc
	TickHandler                    TickHandlerFunc
	UnexpectedDisconnectionHandler UnexpectedDisconnectionFunc
	shutdownFlag                   bool
	closeAfterRespond              bool
	lock                           sync.Mutex
	queue                          chan *Queue
}

type Queue struct {
	ID       int
	ConnInfo *ConnInfo
}

type Conn map[int]*ConnInfo

type ConnInfo struct {
	ID          int
	Input       []byte
	Output      []byte
	Info        evio.Info
	InputStream evio.InputStream
}

type DataHandlerFunc func(connInfo *ConnInfo)
type TickHandlerFunc func() (delay time.Duration)
type UnexpectedDisconnectionFunc func(connInfo *ConnInfo)

const (
	DefaultTickDelayDuration   = time.Millisecond * 100
	InternalServerErrorMessage = "Internal server error: %s"
	BadRequestErrorMessage     = "Bad request: %s"
	ContentTypeTextPlain       = "text/plain"
	Separator                  = "\r\n\r\n"
)

var (
	ErrorMalformedBody   = errors.New("Malformed request body")
	ErrorMalformedHeader = errors.New("Malformed request header")

	ContentLengthRegexp = regexp.MustCompile("Content-Length: ([0-9]+)")
	HostPortOnlyRegexp  = regexp.MustCompile("^:[0-9]{1,5}$")
)

func New(host string, dataHandler DataHandlerFunc, tickHandler TickHandlerFunc, unexpectedDisconnectionHandler UnexpectedDisconnectionFunc) *Ev {
	return &Ev{
		Conn:                           Conn{},
		Host:                           host,
		DataHandler:                    dataHandler,
		TickHandler:                    tickHandler,
		UnexpectedDisconnectionHandler: unexpectedDisconnectionHandler,
		queue: make(chan *Queue, math.MaxUint16),
	}
}

func (e *Ev) Shutdown() {
	e.shutdownFlag = true
}

func (e *Ev) Listen() (err error) {
	var events evio.Events

	events.Serving = serving(e)
	events.Opened = opened(e)
	events.Closed = closed(e)
	events.Data = data(e)

	if e.TickHandler != nil {
		events.Tick = tick(e)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(e)
	}

	if HostPortOnlyRegexp.MatchString(e.Host) {
		e.Host = "0.0.0.0" + e.Host
	}

	return evio.Serve(events, "tcp://"+e.Host)
}

func worker(e *Ev) {
	for {
		q := <-e.queue

		// connInfo.Output will be filled here
		e.DataHandler(q.ConnInfo)

		ok := e.Server.Wake(q.ID)
		if !ok {
			if e.UnexpectedDisconnectionHandler != nil {
				// client already disconnected
				// do something
				e.UnexpectedDisconnectionHandler(q.ConnInfo)
			}
		}
	}
}

func serving(e *Ev) func(server evio.Server) (action evio.Action) {
	return func(srv evio.Server) (action evio.Action) {
		e.Server = srv
		return
	}
}

func opened(e *Ev) func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
	return func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
		// log.Printf("Connected user with id %d and address %s\n", id, info.RemoteAddr.String())

		// add new connection to map
		e.lock.Lock()
		e.Conn[id] = &ConnInfo{
			ID:   id,
			Info: info,
		}
		e.lock.Unlock()

		return
	}
}

func closed(e *Ev) func(id int, err error) (action evio.Action) {
	return func(id int, err error) (action evio.Action) {
		// log.Printf("Disconnected user with id %d and address %s\n", id, e.Conn[id].RemoteAddress.String())

		// remove the current connection from map
		e.lock.Lock()
		delete(e.Conn, id)
		e.lock.Unlock()

		return
	}
}

func data(e *Ev) func(id int, in []byte) (out []byte, action evio.Action) {
	return func(id int, in []byte) (out []byte, action evio.Action) {
		// handle wake up call
		if in == nil {
			out = e.Conn[id].Output
			action = evio.Close

			return
		}

		connInfo := e.Conn[id]

		data := connInfo.InputStream.Begin(in)

		finished, err := isStreamFinished(id, data)
		if err != nil {
			status := http.StatusInternalServerError
			body := []byte(fmt.Sprintf(InternalServerErrorMessage, err))

			if err == ErrorMalformedHeader {
				status = http.StatusBadRequest
				body = []byte(fmt.Sprintf(BadRequestErrorMessage, err))
			}
			if err == ErrorMalformedBody {
				status = http.StatusBadRequest
				body = []byte(fmt.Sprintf(BadRequestErrorMessage, err))
			}

			out = NewRawHTTPResponse(status, ContentTypeTextPlain, body)
			action = evio.Close
			return
		}

		if finished {
			connInfo.Input = data
			e.queue <- &Queue{
				ID:       id,
				ConnInfo: connInfo,
			}
		}

		connInfo.InputStream.End(data)

		return
	}
}

// Find out how to know if the stream is finished or not
func isStreamFinished(id int, data []byte) (finished bool, err error) {
	res := strings.SplitN(string(data), Separator, 2)
	if len(res) < 2 {
		// header probably not completely sent, wait for next tick
		return
	}

	// find content length
	contentLength := 0
	matchesBytes := ContentLengthRegexp.FindSubmatch([]byte(res[0]))
	if len(matchesBytes) == 2 {
		contentLength, err = strconv.Atoi(string(matchesBytes[1]))
	} else {
		// if no Content-Length, return finished (or error?)
		// err = ErrorMalformedHeader
		finished = true
		return
	}
	if err != nil {
		return
	}
	// if contentLength is zero, return finished
	if contentLength == 0 {
		finished = true
		return
	}

	if len(res[1]) == contentLength {
		finished = true
		return
	}

	return
}

func tick(e *Ev) func() (delay time.Duration, action evio.Action) {
	return func() (delay time.Duration, action evio.Action) {
		if e.shutdownFlag {
			action = evio.Shutdown
		}

		delay = time.Duration(DefaultTickDelayDuration)
		if e.TickHandler != nil {
			delay = e.TickHandler()
		}

		return
	}
}

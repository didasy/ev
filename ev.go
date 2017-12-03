package ev

import (
	"net"
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
	ShutdownFlag                   bool
}

type Conn map[int]*ConnInfo

type ConnInfo struct {
	ID            int
	Output        []byte
	LocalAddress  net.Addr
	RemoteAddress net.Addr
}

type DataHandlerFunc func(in []byte, connInfo *ConnInfo)
type TickHandlerFunc func() (delay time.Duration)
type UnexpectedDisconnectionFunc func(in []byte, connInfo *ConnInfo)

const (
	DefaultTickDelayDuration = time.Millisecond * 100
)

func New(host string, dataHandler DataHandlerFunc, tickHandler TickHandlerFunc, unexpectedDisconnectionHandler UnexpectedDisconnectionFunc) *Ev {
	return &Ev{
		Conn:                           Conn{},
		Host:                           host,
		DataHandler:                    dataHandler,
		TickHandler:                    tickHandler,
		UnexpectedDisconnectionHandler: unexpectedDisconnectionHandler,
	}
}

func (e *Ev) Shutdown() {
	e.ShutdownFlag = true
}

func (e *Ev) Listen() (err error) {
	var events evio.Events

	events.Serving = serving(e)
	events.Opened = opened(e)
	events.Data = data(e)
	events.Tick = tick(e)

	return evio.Serve(events, e.Host)
}

func serving(e *Ev) func(server evio.Server) (action evio.Action) {
	return func(srv evio.Server) (action evio.Action) {
		e.Server = srv
		return
	}
}

func opened(e *Ev) func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
	return func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
		e.Conn[id] = &ConnInfo{
			ID:            id,
			LocalAddress:  info.LocalAddr,
			RemoteAddress: info.RemoteAddr,
		}

		return
	}
}

func data(e *Ev) func(id int, in []byte) (out []byte, action evio.Action) {
	return func(id int, in []byte) (out []byte, action evio.Action) {
		// handle wake up call
		if in == nil {
			out = e.Conn[id].Output
			return
		}

		// or handle input
		go func() {
			e.DataHandler(in, e.Conn[id])

			ok := e.Server.Wake(id)
			if !ok {
				// client already disconnected
				// do something
				e.UnexpectedDisconnectionHandler(in, e.Conn[id])
			}
		}()

		return
	}
}

func tick(e *Ev) func() (delay time.Duration, action evio.Action) {
	return func() (delay time.Duration, action evio.Action) {
		if e.ShutdownFlag {
			action = evio.Shutdown
		}

		delay = time.Duration(DefaultTickDelayDuration)
		if e.TickHandler != nil {
			delay = e.TickHandler()
		}

		return
	}
}

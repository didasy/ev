package ev_test

import (
	"testing"
	"time"

	"github.com/JesusIslam/ev"
	"github.com/stretchr/testify/assert"
)

const (
	Host         = "tcp://localhost:5000"
	InputText    = "test"
	DefaultDelay = time.Millisecond * 100
)

var (
	evObj *ev.Ev
)

func dataHandler(in []byte, connInfo *ev.ConnInfo) {
	connInfo.Output = in
}

func tickHandler() (delay time.Duration) {
	delay = time.Duration(DefaultDelay)

	return
}

func TestNew(t *testing.T) {
	evObj = ev.New(Host, dataHandler, tickHandler)
	assert.NotNil(t, evObj)
}

func TestDataHandler(t *testing.T) {
	in := []byte(InputText)
	connInfo := &ev.ConnInfo{}
	evObj.DataHandler(in, connInfo)

	assert.Equal(t, in, connInfo.Output)
}

func TestTickHandler(t *testing.T) {
	delay := evObj.TickHandler()

	assert.Equal(t, DefaultDelay, delay)
}

func TestListen(t *testing.T) {
	go func() {
		err := evObj.Listen()
		assert.Nil(t, err)
	}()

	// wait for listener to be initialized first
	<-time.After(time.Second)
	// shut it down so .Listen will stop blocking
	evObj.Shutdown()
}

// TODO: testing listener with TCP client

package main

import (
	"log"
	"time"

	"github.com/JesusIslam/ev"
)

func main() {
	host := "tcp://localhost:5000"
	e := ev.New(host, dataHandler, nil, nil)
	e.SetCloseAfterRespond(true)
	log.Println("Listening at:", host)
	err := e.Listen()
	if err != nil {
		panic(err)
	}
}

func dataHandler(in []byte, connInfo *ev.ConnInfo) {
	_, body, err := ev.GetHTTPRequest(in)
	if err != nil {
		panic(err)
	}

	// emulate work
	<-time.After(time.Second)

	connInfo.Output = ev.NewRawHTTPResponse(200, "application/json", body)
}

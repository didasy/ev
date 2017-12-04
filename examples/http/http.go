package main

import (
	"log"

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

func dataHandler(connInfo *ev.ConnInfo) {
	_, body, err := ev.GetHTTPRequest(connInfo.Input)
	if err != nil {
		panic(err)
	}

	connInfo.Output = ev.NewRawHTTPResponse(200, "application/json", body)
}

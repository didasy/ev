package main

import (
	"log"

	"github.com/JesusIslam/ev"
)

func main() {
	host := "tcp://localhost:5000"
	e := ev.New(host, dataHandler, nil)
	log.Println("Listening at port: ", host)
	err := e.Listen()
	if err != nil {
		panic(err)
	}
}

func dataHandler(in []byte, connInfo *ev.ConnInfo) {
	log.Print("Input: ", string(in))

	connInfo.Output = []byte("PONG")
}

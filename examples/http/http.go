package main

import (
	"io/ioutil"
	"log"

	"github.com/JesusIslam/ev"
)

func main() {
	host := ":5000"
	e := ev.New(host, dataHandler, nil, nil)
	log.Println("Listening at", host)
	err := e.Listen()
	if err != nil {
		panic(err)
	}
}

func dataHandler(connInfo *ev.ConnInfo) {
	req, err := ev.GetHTTPRequest(connInfo.Input)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	connInfo.Output = ev.NewRawHTTPResponse(200, req.Header.Get("Content-Type"), body)
}

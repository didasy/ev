package ev

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const (
	MinimumHTTPStringFormat = "HTTP/1.1 %d %s\r\nDate: %s\r\nServer: Ev\r\nAccept-Ranges: bytes\r\nConnection: close\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s"
)

var (
	HTTPStringFormat = MinimumHTTPStringFormat
)

func NewRawHTTPResponse(status int, contentType string, body []byte) (res []byte) {
	statusText := http.StatusText(status)
	date := time.Now().Format(time.RFC1123)
	res = []byte(fmt.Sprintf(HTTPStringFormat, status, statusText, date, contentType, len(body), string(body)))

	return
}

func SetHTTPStringFormat(format string) {
	HTTPStringFormat = format
}

func GetHTTPRequest(req []byte) (httpReq *http.Request, err error) {
	reader := bytes.NewReader(req)
	br := bufio.NewReader(reader)
	httpReq, err = http.ReadRequest(br)
	if err != nil {
		return
	}

	return
}

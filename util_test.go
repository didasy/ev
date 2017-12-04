package ev_test

import (
	"net/http"
	"testing"

	"github.com/JesusIslam/ev"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultContentType = "text/plain"

	RawHTTPRequest       = "POST / HTTP/1.1\r\nHost: www.example.com\r\nConnection: keep-alive\r\nContent-Length: 4\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36\r\nCache-Control: no-cache\r\nOrigin: chrome-extension://fhbjgbiflinjbdggehcddcbncdddomop\r\nPostman-Token: cd135163-b2c6-fc69-9889-8ad6f1358e2b\r\nContent-Type: text/plain\r\nAccept: */*\r\nAccept-Encoding: gzip, deflate, br\r\nAccept-Language: en-US,en;q=0.9,id;q=0.8\r\n\r\ntest"
	Method               = "POST"
	Protocol             = "HTTP/1.1"
	HostName             = "www.example.com"
	ContentTypeHeader    = "text/plain"
	ContentTypeHeaderKey = "Content-Type"

	RawHTTPResponse = "HTTP/1.1 200 OK\r\nDate: Mon, 04 Dec 2017 18:12:30 WIB\r\nServer: Ev\r\nAccept-Ranges: bytes\r\nConnection: close\r\nContent-Type: text/plain\r\nContent-Length: 4\r\n\r\ntest"

	DefaultHTTPStringFormat = "format"
)

func TestNewRawHTTPResponse(t *testing.T) {
	res := ev.NewRawHTTPResponse(http.StatusOK, DefaultContentType, []byte(InputText))
	assert.Equal(t, len(RawHTTPResponse), len(res))
}

func TestGetHTTPRequest(t *testing.T) {
	req, body, err := ev.GetHTTPRequest([]byte(RawHTTPRequest))
	assert.Nil(t, err)
	assert.Equal(t, InputText, string(body))
	// assert the request's properties
	assert.Equal(t, Method, req.Method)
	assert.Equal(t, Protocol, req.Proto)
	assert.Equal(t, HostName, req.Host)
	assert.Equal(t, ContentTypeHeader, req.Header.Get(ContentTypeHeaderKey))
}

func TestSetHTTPStringFormat(t *testing.T) {
	ev.SetHTTPStringFormat(DefaultHTTPStringFormat)
	assert.Equal(t, DefaultHTTPStringFormat, ev.HTTPStringFormat)
}

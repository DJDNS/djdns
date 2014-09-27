package app

import (
	"net"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jcelliott/turnpike"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func assertError(t *testing.T, expected string, err error) {
	if expected == "" {
		assert.NoError(t, err)
	} else {
		if assert.Error(t, err) {
			assert.Equal(t, expected, err.Error())
		}
	}
}

func setupRouter() (url string, closer func()) {
	websocket_server := httptest.NewServer(turnpike.NewServer().Handler)
	closer = func() {
		websocket_server.CloseClientConnections()
		websocket_server.Close()
	}
	url = strings.Replace(websocket_server.URL, "http", "deje", 1) + "/root"
	return
}

// Fulfills the dns.ResponseWriter interface with stub functions.
type DummyResponseWriter struct{}

func (drw DummyResponseWriter) RemoteAddr() net.Addr {
	return &net.UDPAddr{}
}
func (drw DummyResponseWriter) WriteMsg(*dns.Msg) error {
	return nil
}
func (drw DummyResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}
func (drw DummyResponseWriter) Close() error {
	return nil
}
func (drw DummyResponseWriter) TsigStatus() error {
	return nil
}
func (drw DummyResponseWriter) TsigTimersOnly(bool) {
	return
}
func (drw DummyResponseWriter) Hijack() {
	return
}

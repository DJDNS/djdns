package app

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jcelliott/turnpike"
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

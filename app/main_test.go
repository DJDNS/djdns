package app

import (
	"bytes"
	"testing"

	"github.com/DJDNS/djdns/server"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestSetupServer(t *testing.T) {
	buf := new(bytes.Buffer)
	conf := ServerConfig{
		RootAlias:   "deje://",    // Uses DEJE scheme, but guarantees failure
		HostAddress: "9.9.9.9:13", // Will also fail, but differently
		ErrorWriter: buf,
	}

	logger, s := setupServer(conf)

	query := new(dns.Msg)
	query.SetQuestion("apple.sauce.", dns.TypeMX)
	s.ServeDNS(DummyResponseWriter{}, query) // Trigger DPG logging
	logger.Printf("Some numbers: %d, %d, %d", 6, -20, 192145)

	expectedOutput := "djdns: No network logging: Error connecting to websocket server: websocket.Dial ws:///ws: dial tcp: missing address" +
		"\nclient 'deje://': Could not open client" +
		"\ndjdns: Error connecting to websocket server: websocket.Dial ws:///ws: dial tcp: missing address" +
		"\ndjdns: Some numbers: 6, -20, 192145\n"
	assert.Equal(t, expectedOutput, buf.String())
	assert.Equal(t, logger, s.Logger)
	assert.Equal(t, "deje://", s.PageGetter.(server.AliasPageGetter).Aliases["<ROOT>"])

	assertError(t, "listen udp 9.9.9.9:13: bind: cannot assign requested address", s.Run(conf.HostAddress))
}

func TestMain(t *testing.T) {
	tests := []struct {
		Argv           []string
		ExpectedOutput string
	}{
		// Bad argument structure
		{
			[]string{"--addr"},
			"djdns: --addr requires argument\n", // slight difference in behavior between test and real
		},
		// Unhostable address
		{
			[]string{"--addr", "9.9.9.9:13", "--root", "model/demo.json"},
			"djdns: No network logging: URL does not start with 'deje://': 'model/demo.json'" +
				"\ndjdns: Starting server on 9.9.9.9:13" +
				"\ndjdns: <ROOT> is 'model/demo.json'" +
				"\ndjdns: listen udp 9.9.9.9:13: bind: cannot assign requested address\n",
		},
		// TODO: Test success, when we can do so without starting an unkillable goroutine
	}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		Main(test.Argv, false, buf)
		assert.Equal(t, test.ExpectedOutput, buf.String())
	}
}

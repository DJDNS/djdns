package app

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPeerWriter(t *testing.T) {
	realHostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("Failed to get hostname: %v", err)
	}

	ws_server_url, ws_closer := setupRouter()
	defer ws_closer()

	// Set this back to normal afterward
	defer func() {
		getHostnameShouldFail = false
	}()

	tests := []struct {
		Conf             ServerConfig
		GetHostnameFails bool

		ExpectedHostname string
		ShouldHaveClient bool
		ExpectedOutput   string
	}{
		// Empty config
		{
			ServerConfig{}, false,
			realHostname,
			false,
			"djdns: No network logging: URL does not start with 'deje://': ''\n",
		},
		// Override the display name
		{
			ServerConfig{DisplayName: "archibald"}, false,
			"archibald",
			false,
			"djdns: No network logging: URL does not start with 'deje://': ''\n",
		},
		// Succesful client creation
		{
			ServerConfig{RootAlias: ws_server_url}, false,
			realHostname,
			true,
			"",
		},
		// Failure to get system hostname
		{
			ServerConfig{RootAlias: ws_server_url}, true,
			"",
			true,
			"djdns: Hostname detection failed: Error for testing purposes\n",
		},
	}
	for _, test := range tests {
		// Capture output
		buf := new(bytes.Buffer)
		test.Conf.ErrorWriter = buf
		getHostnameShouldFail = test.GetHostnameFails

		// Compare resulting object and output
		pw := NewPeerWriter(test.Conf)
		assert.Equal(t, buf, pw.RealWriter)
		assert.Equal(t, test.ExpectedHostname, pw.Hostname)
		assert.Equal(t, test.ShouldHaveClient, pw.Client != nil, "pw.Client was actually %#v", pw.Client)
		assert.Equal(t, test.ExpectedOutput, buf.String())
	}
}

// Used as a mock for DEJE client
type StoringPublisher struct {
	Stored     []interface{}
	ShouldFail bool
}

func (sp *StoringPublisher) Publish(data interface{}) error {
	if sp.ShouldFail {
		return errors.New("Arbitrary failure to publish")
	} else {
		sp.Stored = append(sp.Stored, data)
		return nil
	}
}

func TestPeerWriter_Write(t *testing.T) {
	realHostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("Failed to get hostname: %v", err)
	}

	tests := []struct {
		Conf    ServerConfig
		LogLine string

		PublishShouldFail bool
		ExpectedPublished []interface{}
		ExpectedPrinted   string
	}{
		// Empty config
		{
			ServerConfig{},
			"Example written data",
			false,
			[]interface{}{
				map[string]interface{}{
					"type":  "log",
					"value": "Example written data",
					"host":  realHostname,
				},
			},
			"djdns: No network logging: URL does not start with 'deje://': ''" +
				"\nExample written data",
		},
		// Multiple lines
		{
			ServerConfig{},
			"Line 1\nLine 2",
			false,
			[]interface{}{
				map[string]interface{}{
					"type":  "log",
					"value": "Line 1\nLine 2",
					"host":  realHostname,
				},
			},
			"djdns: No network logging: URL does not start with 'deje://': ''" +
				"\nLine 1" +
				"\nLine 2",
		},
		// Failed publish
		{
			ServerConfig{},
			"Foobar",
			true,
			[]interface{}{},
			"djdns: No network logging: URL does not start with 'deje://': ''" +
				"\ndjdns: Arbitrary failure to publish" +
				"\nFoobar",
		},
		// No client
		{
			ServerConfig{},
			"Bazquux",
			false,
			nil, // Signal not to set client
			"djdns: No network logging: URL does not start with 'deje://': ''" +
				"\nBazquux",
		},
	}

	for _, test := range tests {
		// Capture output
		buf := new(bytes.Buffer)
		test.Conf.ErrorWriter = buf
		pw := NewPeerWriter(test.Conf)
		publisher := StoringPublisher{
			make([]interface{}, 0),
			test.PublishShouldFail,
		}
		if test.ExpectedPublished != nil {
			pw.Client = &publisher
		} else {
			pw.Client = nil
			publisher.Stored = nil
		}

		pw.Write([]byte(test.LogLine))

		assert.Equal(t, test.ExpectedPublished, publisher.Stored)
		assert.Equal(t, test.ExpectedPrinted, buf.String())
	}
}

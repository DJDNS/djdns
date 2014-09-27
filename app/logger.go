package app

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/DJDNS/go-deje"
)

// Tees a write, so that it is both written to the
// underlying writer, and to the WAMP network.
type PeerWriter struct {
	RealWriter io.Writer
	Hostname   string
	Client     Publisher // DEJE client or compatible API
}

func NewPeerWriter(conf ServerConfig) PeerWriter {
	pw := PeerWriter{conf.ErrorWriter, conf.DisplayName, nil}
	logger := pw.GetLogger()
	var err error

	pw.Client, err = getLoggingClient(conf.RootAlias)
	if err != nil {
		// It's fine for peer_writer_client to be nil.
		// We just need to make sure the error gets printed non-fatally.
		logger.Printf("No network logging: %v\n", err)
		pw.Client = nil // ensure this, trust no one
	}
	if pw.Hostname == "" {
		pw.Hostname, err = getHostname()
		if err != nil {
			logger.Printf("Hostname detection failed: %v\n", err)
			pw.Hostname = ""
		}
	}

	return pw
}

func (pw PeerWriter) GetLogger() *log.Logger {
	return log.New(pw.RealWriter, "djdns: ", 0)
}

func (pw PeerWriter) Write(p []byte) (n int, err error) {
	if err = pw.writeToNetwork(p); err != nil {
		// Don't crash, just fall back to simpler logger
		pw.GetLogger().Println(err)
	}
	return pw.RealWriter.Write(p)
}

func (pw PeerWriter) writeToNetwork(p []byte) error {
	// Exit early if no client
	if pw.Client == nil {
		return nil
	}

	// Construct and publish message object
	data := map[string]interface{}{
		"type":  "log",
		"value": string(p),
	}
	if pw.Hostname != "" {
		data["host"] = pw.Hostname
	}
	return pw.Client.Publish(data)
}

// ------------------------------------------------

type Publisher interface {
	Publish(interface{}) error
}

func getLoggingClient(url string) (*deje.Client, error) {
	router, topic, err := deje.GetRouterAndTopic(url)
	if err != nil {
		return nil, err
	}
	client := deje.NewClient(topic)
	return &client, client.Connect(router)
}

var getHostnameShouldFail = false

func getHostname() (string, error) {
	if getHostnameShouldFail {
		return "", errors.New("Error for testing purposes")
	} else {
		return os.Hostname()
	}
}

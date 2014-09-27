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
	var logger = log.New(conf.ErrorWriter, "djdns: ", 0)

	peer_writer_client, err := getLoggingClient(conf.RootAlias)
	if err != nil {
		// It's fine for peer_writer_client to be nil.
		// We just need to make sure the error gets printed non-fatally.
		logger.Printf("No network logging: %v\n", err)
		peer_writer_client = nil // ensure this, trust no one
	}
	hostname := conf.DisplayName
	if hostname == "" {
		hostname, err = getHostname()
		if err != nil {
			logger.Printf("Hostname detection failed: %v\n", err)
			hostname = ""
		}
	}
	return PeerWriter{conf.ErrorWriter, hostname, peer_writer_client}
}

func (pl PeerWriter) Write(p []byte) (n int, err error) {
	if err = pl.writeToNetwork(p); err != nil {
		// Don't crash, just fall back to simpler logger
		var logger = log.New(pl.RealWriter, "djdns: ", 0)
		logger.Println(err)
	}
	return pl.RealWriter.Write(p)
}

func (pl PeerWriter) writeToNetwork(p []byte) error {
	// Exit early if no client
	if pl.Client == nil {
		return nil
	}

	// Construct and publish message object
	data := map[string]interface{}{
		"type":  "log",
		"value": string(p),
	}
	if pl.Hostname != "" {
		data["host"] = pl.Hostname
	}
	return pl.Client.Publish(data)
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

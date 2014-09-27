package app

import (
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
	Client     *deje.Client
}

func NewPeerWriter(conf ServerConfig) PeerWriter {
	peer_writer_client, err := getLoggingClient(conf.RootAlias)
	if err != nil {
		// It's fine for peer_writer_client to be nil.
		// We just need to make sure the error gets printed non-fatally.
		log.Printf("No network logging: %v\n", err)
		peer_writer_client = nil // ensure this, trust no one
	}
	hostname := conf.DisplayName
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			log.Printf("Hostname detection failed: %v\n", err)
			hostname = ""
		}
	}
	return PeerWriter{os.Stderr, hostname, peer_writer_client}
}

func (pl PeerWriter) Write(p []byte) (n int, err error) {
	if err = pl.writeToNetwork(p); err != nil {
		// Don't crash, just fall back to simpler logger
		log.Println(err)
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

func getLoggingClient(url string) (*deje.Client, error) {
	router, topic, err := deje.GetRouterAndTopic(url)
	if err != nil {
		return nil, err
	}
	client := deje.NewClient(topic)
	return &client, client.Connect(router)
}

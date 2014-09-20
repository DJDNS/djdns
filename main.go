package main

//import "flag"
import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/DJDNS/djdns/server"
	"github.com/DJDNS/go-deje"
)

var root_alias = flag.String("root", "deje://localhost:8080/root", "Target URL to serve as <ROOT>")
var display_name = flag.String("display-name", "", "Hostname to provide in network log messages")

type PeerWriter struct {
	RealWriter io.Writer
	Hostname   string
	Client     *deje.Client
}

func (pl PeerWriter) Write(p []byte) (n int, err error) {
	data := map[string]interface{}{
		"type":  "log",
		"value": string(p),
	}
	if pl.Hostname != "" {
		data["host"] = pl.Hostname
	}
	pl.Client.Publish(data)
	return pl.RealWriter.Write(p)
}

func getLoggingClient(url string) (*deje.Client, error) {
	router, topic, err := deje.GetRouterAndTopic(url)
	if err != nil {
		return nil, err
	}
	client := deje.NewClient(topic)
	return &client, client.Connect(router)
}

func makePeerWriter(url string) (PeerWriter, error) {
	peer_writer_client, err := getLoggingClient(url)
	if err != nil {
		return PeerWriter{}, err
	}
	hostname := *display_name
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			log.Printf("Hostname detection failed: %v\n", err)
			hostname = ""
		}
	}
	return PeerWriter{os.Stderr, hostname, peer_writer_client}, nil
}

func main() {
	flag.Parse()
	addr := "0.0.0.0:9953"

	var log_writer io.Writer
	var err error
	log_writer, err = makePeerWriter(*root_alias)
	if err != nil {
		log.Printf("No network logging: %v\n", err)
		log_writer = os.Stderr
	}
	logger := log.New(log_writer, "djdns: ", 0)

	spgc := server.NewStandardPGConfig(log_writer)
	spgc.Alias.Aliases = map[string]string{
		"<ROOT>": *root_alias,
	}

	s := server.NewServer(spgc.Alias)
	s.Logger = logger

	logger.Printf("Starting server on %s", addr)
	logger.Printf("<ROOT> is '%s'", *root_alias)
	err = s.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

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

type PeerWriter struct {
	RealWriter io.Writer
	Client     *deje.Client
}

func (pl PeerWriter) Write(p []byte) (n int, err error) {
	pl.Client.Publish(map[string]interface{}{
		"type":  "log",
		"value": string(p),
	})
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

func main() {
	flag.Parse()
	addr := "0.0.0.0:9953"

	peer_writer_client, err := getLoggingClient(*root_alias)
	if err != nil {
		log.Fatal(err)
	}
	log_writer := PeerWriter{os.Stderr, peer_writer_client}
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

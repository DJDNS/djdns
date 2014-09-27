package main

//import "flag"
import (
	"io"
	"log"
	"os"

	"github.com/DJDNS/djdns/server"
	"github.com/DJDNS/go-deje"
	"github.com/docopt/docopt-go"
)

var version = "djdns 0.0.12"
var usage = `djdns

Usage: djdns [options]

Options:
    --root=<root>         Target URL to serve as <ROOT> [default: deje://localhost:8080/root]
    --addr=<address>      host:port to expose DNS on [default: 0.0.0.0:9953]
    --display-name=<name> Hostname to provide in network log messages.
    -h --help             Show this message.
    --version             Print the version number.
`

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

func makePeerWriter(url, display_name string) (PeerWriter, error) {
	peer_writer_client, err := getLoggingClient(url)
	if err != nil {
		return PeerWriter{}, err
	}
	hostname := display_name
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			log.Printf("Hostname detection failed: %v\n", err)
			hostname = ""
		}
	}
	return PeerWriter{os.Stderr, hostname, peer_writer_client}, nil
}

func getShellArg(arguments map[string]interface{}, key string) string {
	arg := arguments[key]
	if arg == nil {
		return ""
	} else {
		return arg.(string)
	}
}

func main() {
	arguments, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		log.Fatal(err)
	}

	root_alias := getShellArg(arguments, "--root")
	display_name := getShellArg(arguments, "--display-name")
	addr := getShellArg(arguments, "--addr")

	var log_writer io.Writer
	log_writer, err = makePeerWriter(root_alias, display_name)
	if err != nil {
		log.Printf("No network logging: %v\n", err)
		log_writer = os.Stderr
	}
	logger := log.New(log_writer, "djdns: ", 0)

	spgc := server.NewStandardPGConfig(log_writer)
	spgc.Alias.Aliases = map[string]string{
		"<ROOT>": root_alias,
	}

	s := server.NewServer(spgc.Alias)
	s.Logger = logger

	logger.Printf("Starting server on %s", addr)
	logger.Printf("<ROOT> is '%s'", root_alias)
	err = s.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

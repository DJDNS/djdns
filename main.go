package main

//import "flag"
import (
	"flag"
	"log"
	"os"

	"github.com/DJDNS/djdns/server"
)

var root_alias = flag.String("root", "deje://localhost:8080/root", "Target URL to serve as <ROOT>")

func main() {
	flag.Parse()
	addr := "127.0.0.1:9953"

	logger := log.New(os.Stderr, "djdns: ", 0)

	spgc := server.NewStandardPGConfig(os.Stderr)
	spgc.Alias.Aliases = map[string]string{
		"<ROOT>": *root_alias,
	}

	s := server.NewServer(spgc.Alias)
	s.Logger = logger

	logger.Printf("Starting server on %s", addr)
	logger.Printf("<ROOT> is '%s'", *root_alias)
	err := s.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

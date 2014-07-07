package main

//import "flag"
import (
	"github.com/campadrenalin/djdns/server"
	"log"
	"os"
)

func main() {
	root_alias := "./model/demo.json"
	addr := "127.0.0.1:9953"

	logger := log.New(os.Stderr, "djdns: ", 0)

	spgc := server.NewStandardPGConfig()
	spgc.Alias.Aliases = map[string]string{
		"<ROOT>": root_alias,
	}

	s := server.NewServer(spgc.Alias)
	s.Logger = logger

	logger.Printf("Starting server on %s", addr)
	err := s.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

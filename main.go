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
	aliases := map[string]string{
		"<ROOT>": root_alias,
	}

	logger := log.New(os.Stderr, "djdns: ", 0)

	s := server.NewServer()
	s.Logger = logger
	s_aliases := s.PageGetter.(server.AliasPageGetter).Aliases
	for k, v := range aliases {
		s_aliases[k] = v
	}

	logger.Printf("Starting server on %s", addr)
	err := s.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

package app

import (
	"io"
	"log"

	"github.com/DJDNS/djdns/server"
)

func Main(argv []string) {
	conf, err := Parse(argv)
	if err != nil {
		log.Fatal(err)
	}

	var log_writer io.Writer = NewPeerWriter(conf)
	logger := log.New(log_writer, "djdns: ", 0)

	spgc := server.NewStandardPGConfig(log_writer)
	spgc.Alias.Aliases = map[string]string{
		"<ROOT>": conf.RootAlias,
	}

	s := server.NewServer(spgc.Alias)
	s.Logger = logger

	logger.Printf("Starting server on %s", conf.HostAddress)
	logger.Printf("<ROOT> is '%s'", conf.RootAlias)
	err = s.Run(conf.HostAddress)
	if err != nil {
		log.Fatal(err)
	}
}

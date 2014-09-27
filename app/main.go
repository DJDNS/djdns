package app

import (
	"io"
	"log"

	"github.com/DJDNS/djdns/server"
)

func setupServer(conf ServerConfig) (*log.Logger, server.DjdnsServer) {
	pw := NewPeerWriter(conf)
	logger := pw.GetLogger()

	spgc := server.NewStandardPGConfig(pw)
	spgc.Alias.Aliases = map[string]string{
		"<ROOT>": conf.RootAlias,
	}

	s := server.NewServer(spgc.Alias)
	s.Logger = logger

	return logger, s
}

func Main(argv []string, exit bool, debug_writer io.Writer) {
	conf, err := Parse(argv, exit)
	if err != nil {
		log.New(debug_writer, "djdns: ", 0).Println(err)
		return
	}

	conf.ErrorWriter = debug_writer
	logger, s := setupServer(conf)

	logger.Printf("Starting server on %s", conf.HostAddress)
	logger.Printf("<ROOT> is '%s'", conf.RootAlias)
	if err := s.Run(conf.HostAddress); err != nil {
		logger.Println(err)
	}
}

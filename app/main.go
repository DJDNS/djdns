package app

import (
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

func Main(argv []string, exit bool) error {
	conf, err := Parse(argv, exit)
	if err != nil {
		return err
	}

	logger, s := setupServer(conf)

	logger.Printf("Starting server on %s", conf.HostAddress)
	logger.Printf("<ROOT> is '%s'", conf.RootAlias)
	return s.Run(conf.HostAddress)
}

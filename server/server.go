package server

import "github.com/campadrenalin/djdns/model"

type DjdnsServer struct {
	Port int
	Root model.Page
}

func NewServer() DjdnsServer {
	var server DjdnsServer
	server.Port = 9953
	return server
}

func (ds *DjdnsServer) GetRecords(q string) []model.Record {
	branch := ds.Root.GetBranchForQuery(q)
	return branch.Records
}

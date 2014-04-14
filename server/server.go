package server

import "github.com/campadrenalin/djdns/model"

type DjdnsServer struct {
	Port int
	Root model.Page
}

// Initialize a DjdnsServer with default values.
//
// This does not start service - you still need to call
// DjdnsServer.Run(), possibly as a goroutine.
func NewServer() DjdnsServer {
	var server DjdnsServer
	server.Port = 9953
	return server
}

// Returns nil slice, if no such branch.
func (ds *DjdnsServer) GetRecords(q string) []model.Record {
	branch := ds.Root.GetBranchForQuery(q)
	if branch == nil {
		return nil
	} else {
		return branch.Records
	}
}

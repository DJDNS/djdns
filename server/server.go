package server

import (
	"github.com/campadrenalin/djdns/model"
	"github.com/miekg/dns"
)

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

// Construct a response for a single DNS request.
func (ds *DjdnsServer) Handle(query *dns.Msg) (*dns.Msg, error) {
	response := new(dns.Msg)
	response.Question = query.Question
	if len(query.Question) > 0 {
		// Ignore secondary questions
		question := query.Question[0]
		records := ds.GetRecords(question.Name)
		response.Answer = make([]dns.RR, len(records))
		for i, record := range records {
			answer, err := record.ToDns()
			if err != nil {
				return nil, err
			}
			response.Answer[i] = answer
		}
		response.Ns = make([]dns.RR, 0)
		response.Extra = make([]dns.RR, 0)
	}

	return response, nil
}

func (ds *DjdnsServer) ServeDNS(rw dns.ResponseWriter, r *dns.Msg) {
	// TODO: Handle errors
	response, _ := ds.Handle(r)
	/*
	   if err != nil {
	       response = new(dns.Msg)
	       response.SetRcode(r, dns.RcodeNameError)
	   }
	*/
	// TODO: Handle errors here too
	_ = rw.WriteMsg(response)
}

func (ds *DjdnsServer) Run(addr string) error {
	server := new(dns.Server)
	server.Addr = addr
	server.Net = "udp"
	server.Handler = ds
	return server.ListenAndServe()
}

func (ds *DjdnsServer) Close() {
}

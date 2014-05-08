package server

import (
	"github.com/campadrenalin/djdns/model"
	"github.com/miekg/dns"
	"log"
)

type DjdnsServer struct {
	Port       int
	PageGetter PageGetter
	Logger     *log.Logger
}

// Initialize a DjdnsServer with default values.
//
// This does not start service - you still need to call
// DjdnsServer.Run(), possibly as a goroutine.
func NewServer() DjdnsServer {
	return DjdnsServer{
		Port:       9953,
		PageGetter: NewAliasPageGetter(NewFilePageGetter()),
	}
}

func (ds *DjdnsServer) GetRecords(q string) ([]model.Record, error) {
	page, err := ds.PageGetter.GetPage("<ROOT>", nil)
	if err != nil {
		return nil, err
	}
	branch := page.Data.GetBranchForQuery(q)
	if branch == nil {
		return nil, nil
	} else {
		return branch.Records, nil
	}
}

// Construct a response for a single DNS request.
func (ds *DjdnsServer) Handle(query *dns.Msg) (*dns.Msg, error) {
	response := new(dns.Msg)
	response.MsgHdr.Id = query.MsgHdr.Id
	response.Question = query.Question
	if len(query.Question) > 0 {
		// Ignore secondary questions
		question := query.Question[0]
		records, err := ds.GetRecords(question.Name)
		if err != nil {
			return nil, err
		}
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
	response, err := ds.Handle(r)
	if err != nil {
		if ds.Logger != nil {
			ds.Logger.Print(err)
		}
		response = new(dns.Msg)
		response.SetRcode(r, dns.RcodeServerFailure)
	}
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

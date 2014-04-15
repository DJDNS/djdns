package model

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
)

type Record struct {
	DomainName string      `json:"domain_name"`
	Rttl       int         `json:"rttl"`
	Rclass     string      `json:"rclass"`
	Rtype      string      `json:"rtype"`
	Rdata      interface{} `json:"rdata"`
}

// Fill in data that was left implied via defaults
func (r *Record) Normalize() {
	if r.Rtype == "" {
		r.Rtype = "A"
	}
}

func (r *Record) RdataString() (string, error) {
	var rdata string
	var ok bool
	var err error

	switch r.Rtype {
	case "A":
	case "MX":
		rdata, ok = r.Rdata.(string)
	default:
		return rdata, errors.New("Unknown Rtype")
	}

	if !ok {
		err = errors.New("Rdata was wrong type for Rtype")
	}
	return rdata, err
}

// Export the data as dns.RR
func (r *Record) ToDns() (dns.RR, error) {
	rdata, err := r.RdataString()
	if err != nil {
		return nil, err
	}
	rr_string := fmt.Sprintf(
		"%s %d %s %s %s",
		r.DomainName,
		r.Rttl,
		r.Rclass,
		r.Rtype,
		rdata,
	)
	return dns.NewRR(rr_string)
}

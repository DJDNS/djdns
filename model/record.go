package model

type Record struct {
	DomainName string      `json:"domain_name"`
	Rtype      string      `json:"rtype"`
	Rdata      interface{} `json:"rdata"`
}

// Fill in data that was left implied via defaults
func (sr *Record) Normalize() {
	if sr.Rtype == "" {
		sr.Rtype = "A"
	}
}

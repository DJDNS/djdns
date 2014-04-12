package model

//import "github.com/miekg/dns"
import djutil "github.com/campadrenalin/go-deje/util"

type Page struct {
	Meta     PageMeta
	Branches []Branch
}

type PageMeta struct {
	Authority string
	Contact   string
	About     string
	Policy    string
}

type Branch struct {
	Selector string
	Targets  []string
	Records  []Record
}

type Record struct {
	DomainName string      `json:"domain_name"`
	Rtype      string      `json:"rtype"`
	Rdata      interface{} `json:"rdata"`
}

func (sp *Page) LoadFrom(data interface{}) error {
	err := djutil.CloneMarshal(data, sp)
	if err != nil {
		return err
	}

	for b := range sp.Branches {
		sp.Branches[b].Normalize()
	}
	return nil
}

// Set default values for missing data
func (sb *Branch) Normalize() {
	for r := range sb.Records {
		sb.Records[r].Normalize()
	}
}

// Set default values for missing data
func (sr *Record) Normalize() {
	if sr.Rtype == "" {
		sr.Rtype = "A"
	}
}

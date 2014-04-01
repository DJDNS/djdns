package main

//import "github.com/miekg/dns"
import djutil "github.com/campadrenalin/go-deje/util"

type SerializedPage struct {
	Meta     PageMeta
	Branches []SerializedBranch
}

type PageMeta struct {
	Authority string
	Contact   string
	About     string
	Policy    string
}

type SerializedBranch struct {
	Selector string
	Targets  []string
	Records  []SerializedRecord
}

type SerializedRecord struct {
	DomainName string      `json:"domain_name"`
	Rtype      string      `json:"rtype"`
	Rdata      interface{} `json:"rdata"`
}

func (sp *SerializedPage) LoadFrom(data interface{}) error {
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
func (sb *SerializedBranch) Normalize() {
	for r := range sb.Records {
		sb.Records[r].Normalize()
	}
}

// Set default values for missing data
func (sr *SerializedRecord) Normalize() {
	if sr.Rtype == "" {
		sr.Rtype = "A"
	}
}

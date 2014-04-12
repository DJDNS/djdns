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

func (sp *Page) LoadFrom(data interface{}) error {
	err := djutil.CloneMarshal(data, sp)
	if err != nil {
		return err
	}

	sp.Normalize()
	return nil
}

// Fill in data that was left implied via defaults
func (sp *Page) Normalize() {
	for b := range sp.Branches {
		sp.Branches[b].Normalize()
	}
}

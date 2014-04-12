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

// Does not forward errors from Branch.Matches (which
// would error on bad regexes, for example).
//
// Result will either be a valid Branch pointer or nil,
// depending on whether a matching branch exists. If
// there are multiple branches that could have matched,
// you will always get the first one in the array. This
// is by design, and why branch order matters.
func (sp *Page) GetBranchForQuery(query string) *Branch {
	for b := range sp.Branches {
		// TODO: Handle errors
		matched, _ := sp.Branches[b].Matches(query)
		if matched {
			return &sp.Branches[b]
		}
	}
	return nil
}

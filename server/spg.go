package server

import (
	"fmt"
	"net/url"
)

// Hands off to another PageGetter based on URL scheme
type SchemePageGetter struct {
	Children map[string]PageGetter
}

func NewSchemePageGetter() SchemePageGetter {
	return SchemePageGetter{
		Children: make(map[string]PageGetter),
	}
}

func (spg SchemePageGetter) GetPage(urlstr string, ab Aborter) (Page, error) {
	urlobj, err := url.Parse(urlstr)
	if err != nil {
		return Page{}, err
	}
	scheme := urlobj.Scheme
	child, ok := spg.Children[scheme]
	if ok {
		return child.GetPage(urlstr, ab)
	} else {
		return Page{}, fmt.Errorf(
			"No PageGetter registered for scheme '%s'",
			scheme,
		)
	}
}

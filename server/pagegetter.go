package server

import (
	"time"

	"github.com/campadrenalin/djdns/model"
)

type Aborter <-chan time.Time

type PageGetter interface {
	GetPage(url string, timeout Aborter) (Page, error)
}

type Page struct {
	Url     string
	Expired bool
	Data    model.Page
}

type FilePageGetter struct {
	Directory string
}

// Create a FilePageGetter with default properties
func NewFilePageGetter() FilePageGetter {
	return FilePageGetter{
		".",
	}
}

func (fpg FilePageGetter) GetPage(url string, timeout Aborter) (Page, error) {
	json, err := model.GetJSONFromFile(url)
	if err != nil {
		return Page{}, err
	}

	page := Page{
		Url:     url,
		Expired: false,
	}
	page.Data.LoadFrom(json)
	return page, nil
}

type StandardPGConfig struct {
	Alias  AliasPageGetter
	File   FilePageGetter
	Scheme SchemePageGetter
}

func NewStandardPGConfig() (spgc StandardPGConfig) {
	spgc.File = NewFilePageGetter()
	spgc.Scheme = NewSchemePageGetter()
	spgc.Scheme.Children["file"] = spgc.File
	spgc.Scheme.Children[""] = spgc.File
	spgc.Alias = NewAliasPageGetter(spgc.Scheme)

	return spgc
}

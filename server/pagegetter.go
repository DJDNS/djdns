package server

import (
	"github.com/campadrenalin/djdns/model"
)

type Aborter chan interface{}

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

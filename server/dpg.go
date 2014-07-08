package server

/*
import (
	"fmt"
	"net/url"
)
*/

// Retrieves DEJE documents, and uses their content
type DejePageGetter struct {
}

func NewDejePageGetter() DejePageGetter {
	return DejePageGetter{}
}

func (pg DejePageGetter) getRouterUrl(string) string {
	return "Not a real value"
}

func (dpg DejePageGetter) GetPage(urlstr string, ab Aborter) (Page, error) {
	return Page{}, nil
}

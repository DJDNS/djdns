package server

import (
	"net/url"
	"strings"
)

// Retrieves DEJE documents, and uses their content
type DejePageGetter struct {
}

func NewDejePageGetter() DejePageGetter {
	return DejePageGetter{}
}

func (pg DejePageGetter) getRouterUrl(deje_url string) (string, error) {
	// Handle schemeless URIs
	valid_prefixes := []string{"http://", "https://", "ws://", "//"}
	has_valid_prefix := false
	for _, prefix := range valid_prefixes {
		if strings.HasPrefix(deje_url, prefix) {
			has_valid_prefix = true
		}
	}
	if !has_valid_prefix {
		deje_url = "ws://" + deje_url
	}

	// Parse and fix up
	url_obj, err := url.Parse(deje_url)
	if err != nil {
		return "", err
	}
	url_obj.Scheme = "ws"
	url_obj.Path = "/ws"
	return url_obj.String(), nil
}

func (dpg DejePageGetter) GetPage(urlstr string, ab Aborter) (Page, error) {
	return Page{}, nil
}

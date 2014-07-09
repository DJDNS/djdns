package server

import (
	"io"
	"log"
	"net/url"
	"strings"

	deje "github.com/campadrenalin/go-deje"
	deje_doc "github.com/campadrenalin/go-deje/document"
)

// Retrieves DEJE documents, and uses their content
type DejePageGetter struct {
	clients map[string]*deje.SimpleClient
}

func NewDejePageGetter() DejePageGetter {
	return DejePageGetter{
		clients: make(map[string]*deje.SimpleClient),
	}
}

// Handle schemeless URIs
func (pg DejePageGetter) unbareUrl(deje_url string) string {
	valid_prefixes := []string{"http://", "https://", "ws://", "deje://", "//"}
	has_valid_prefix := false
	for _, prefix := range valid_prefixes {
		if strings.HasPrefix(deje_url, prefix) {
			has_valid_prefix = true
		}
	}
	if !has_valid_prefix {
		deje_url = "ws://" + deje_url
	}
	return deje_url
}

func (pg DejePageGetter) getRouterUrl(deje_url string) (string, error) {
	url_obj, err := url.Parse(pg.unbareUrl(deje_url))
	if err != nil {
		return "", err
	}
	url_obj.Scheme = "ws"
	url_obj.Path = "/ws"
	return url_obj.String(), nil
}

func (pg DejePageGetter) getTopic(deje_url string) (string, error) {
	url_obj, err := url.Parse(pg.unbareUrl(deje_url))
	if err != nil {
		return "", err
	}
	url_obj.Scheme = "deje"
	if url_obj.Path == "" {
		url_obj.Path = "/"
	}
	return url_obj.String(), nil
}

func (pg DejePageGetter) getDoc(deje_url string, w io.Writer) (*deje_doc.Document, error) {
	client, ok := pg.clients[deje_url]
	if !ok {
		ws_url, err := pg.getRouterUrl(deje_url)
		if err != nil {
			return nil, err
		}

		topic, err := pg.getTopic(deje_url)
		if err != nil {
			return nil, err
		}

		var logger *log.Logger
		if w != nil {
			logger = log.New(w, "client '"+deje_url+"': ", 0)
		}

		client = deje.NewSimpleClient(topic, logger)
		err = client.Connect(ws_url)
		if err != nil {
			return nil, err
		}

		pg.clients[deje_url] = client
	}
	return client.GetDoc(), nil
}

func (dpg DejePageGetter) GetPage(urlstr string, ab Aborter) (Page, error) {
	return Page{}, nil
}

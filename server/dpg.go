package server

import (
	"errors"
	"io"
	"log"
	"net/url"
	"strings"

	deje "github.com/campadrenalin/go-deje"
	deje_doc "github.com/campadrenalin/go-deje/document"
	deje_state "github.com/campadrenalin/go-deje/state"
)

type dejeClientData struct {
	Client  *deje.SimpleClient
	Waiting bool
	Waiter  chan struct{}
}

func NewDCD(c *deje.SimpleClient) *dejeClientData {
	dcd := &dejeClientData{
		Client:  c,
		Waiting: true,
		Waiter:  make(chan struct{}),
	}
	// TODO: Use better condition for "synced"
	c.SetPrimitiveCallback(func(p deje_state.Primitive) {
		if dcd.Waiting {
			close(dcd.Waiter)
			dcd.Waiting = false
		}
	})
	return dcd
}

// Retrieves DEJE documents, and uses their content
type DejePageGetter struct {
	clients map[string]*dejeClientData
	writer  io.Writer
}

func NewDejePageGetter(w io.Writer) DejePageGetter {
	return DejePageGetter{
		clients: make(map[string]*dejeClientData),
		writer:  w,
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

func (pg DejePageGetter) getDoc(deje_url string) (*deje_doc.Document, error) {
	dcd, ok := pg.clients[deje_url]
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
		if pg.writer != nil {
			logger = log.New(pg.writer, "client '"+deje_url+"': ", 0)
		}

		client := deje.NewSimpleClient(topic, logger)
		err = client.Connect(ws_url)
		if err != nil {
			return nil, err
		}

		dcd = NewDCD(client)
		pg.clients[deje_url] = dcd
	}
	return dcd.Client.GetDoc(), nil
}

func (pg DejePageGetter) GetPage(urlstr string, ab Aborter) (Page, error) {
	doc, err := pg.getDoc(urlstr)
	if err != nil {
		return Page{}, err
	}

	// We know this exists - we just ensured it with getDoc
	dcd := pg.clients[urlstr]

	select {
	case <-ab:
		return Page{}, errors.New("DEJE sync timed out")
	case <-dcd.Waiter:
		page := Page{
			Url: urlstr,
		}
		err = page.Data.LoadFrom(doc.State.Export())
		return page, err
	}
}

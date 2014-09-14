package server

import (
	"errors"
	"io"
	"log"

	deje "github.com/DJDNS/go-deje"
	deje_doc "github.com/DJDNS/go-deje/document"
	deje_state "github.com/DJDNS/go-deje/state"
)

type dejeClientData struct {
	Client  *deje.SimpleClient
	Waiting bool
	Waiter  chan struct{}
	LastTip string
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

		var tipstr string
		if c.Tip != nil {
			tipstr = c.Tip.Hash()
		} else {
			tipstr = "nil"
		}
		if tipstr != dcd.LastTip {
			c.Log("Tip is " + tipstr)
			dcd.LastTip = tipstr
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

func (pg DejePageGetter) getDoc(deje_url string) (*deje_doc.Document, error) {
	dcd, ok := pg.clients[deje_url]
	if !ok {
		var logger *log.Logger
		if pg.writer != nil {
			logger = log.New(pg.writer, "client '"+deje_url+"': ", 0)
		}

		client, err := deje.Open(deje_url, logger, nil)
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

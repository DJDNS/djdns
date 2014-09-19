package server

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DJDNS/go-deje"
	"github.com/jcelliott/turnpike"
)

type UrlTest struct {
	Input  string
	Output string
}

func setup_deje_env(t *testing.T) (string, string, func(), *deje.SimpleClient, deje.Client) {
	tp := turnpike.NewServer()
	router := httptest.NewServer(tp.Handler)
	deje_url := strings.Replace(router.URL, "http", "deje", 1) + "/"
	router_url, topic, err := deje.GetRouterAndTopic(deje_url)
	if err != nil {
		t.Fatal(err)
	}
	closer := func() {
		router.CloseClientConnections()
		router.Close()
	}

	clever, err := deje.Open(deje_url, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	dumb := deje.NewClient(topic)
	err = dumb.Connect(router_url)
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(50 * time.Millisecond)

	return router_url, topic, closer, clever, dumb
}

func TestDejePageGetter_getDoc(t *testing.T) {
	pg := NewDejePageGetter(nil)
	_, topic, closer, _, _ := setup_deje_env(t)
	defer closer()

	doc, err := pg.getDoc(topic)
	if err != nil {
		t.Fatal(err)
	}
	doc2, err := pg.getDoc(topic)
	if err != nil {
		t.Fatal(err)
	}

	if doc2 != doc {
		t.Fatalf("Did not return same pointer for both documents - %v vs %v", doc, doc2)
	}
}

func TestDejePageGetter_getDoc_logging(t *testing.T) {
	buf := new(bytes.Buffer)
	pg := NewDejePageGetter(buf)
	_, topic, closer, _, dumb := setup_deje_env(t)
	defer closer()

	// Use fresh log, so hanger-on client has no access
	buf = new(bytes.Buffer)
	pg.writer = buf

	// Create client, wait for it to connect
	_, err := pg.getDoc(topic)
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(50 * time.Millisecond)

	// Put something in the log, and wait for it to broadcast
	if err = dumb.Publish("foo"); err != nil {
		t.Fatal(err)
	}
	<-time.After(50 * time.Millisecond)

	heading := "client '" + topic + "': "
	expected_log := heading + "Tip is nil\n" +
		heading + "Non-{} message\n"
	got_log := buf.String()
	if got_log != expected_log {
		t.Fatalf("Log content was wrong.\nexp: '%s'\ngot: '%s'", expected_log, got_log)
	}

}

func TestDejePageGetter_GetPage_Abort(t *testing.T) {
	pg := NewDejePageGetter(nil)
	_, topic, closer, _, _ := setup_deje_env(t)
	defer closer()

	ab := time.After(0)
	_, err := pg.GetPage(topic, ab)
	if err == nil {
		t.Fatal("Should have failed due to timeout")
	}

	expected_error := "DEJE sync timed out"
	got_error := err.Error()
	if got_error != expected_error {
		t.Fatalf("Expected '%s', got '%s'", expected_error, got_error)
	}
}

func TestDejePageGetter_GetPage(t *testing.T) {
	buf := new(bytes.Buffer)
	pg := NewDejePageGetter(buf)
	_, topic, closer, clever, _ := setup_deje_env(t)
	defer closer()

	// Set up event
	doc := clever.GetDoc()
	event := doc.NewEvent("SET")
	event.Arguments["path"] = []interface{}{}
	event.Arguments["value"] = map[string]interface{}{
		"meta": map[string]interface{}{
			"authority": "Example Authority",
		},
	}
	event.Register()
	clever.Promote(event)

	ab := time.After(50 * time.Millisecond)
	page, err := pg.GetPage(topic, ab)
	if err != nil {
		t.Fatal(err)
	}

	got_auth := page.Data.Meta.Authority
	exp_auth := "Example Authority"
	if got_auth != exp_auth {
		t.Error("Data was not synced from remote host")
		t.Error(page)
	}
}

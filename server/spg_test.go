package server

import (
	"errors"
	"reflect"
	"testing"
)

type DummyPageGetter struct {
	LastCallUrl     string
	LastCallAborter Aborter
	PageData        Page
	ErrorData       error
}

func (dpg *DummyPageGetter) GetPage(url string, ab Aborter) (Page, error) {
	dpg.LastCallUrl = url
	dpg.LastCallAborter = ab
	return dpg.PageData, dpg.ErrorData
}

func setup() SchemePageGetter {
	pg := NewSchemePageGetter()

	dummy1 := DummyPageGetter{
		PageData: Page{
			Url:     "eggs",
			Expired: true,
		},
	}
	dummy1.PageData.Data.Meta.Policy = "Stinky"
	pg.Children["dummy1"] = &dummy1

	dummy2 := DummyPageGetter{
		PageData:  Page{},
		ErrorData: errors.New("DUMMY ERROR"),
	}
	pg.Children["dummy2"] = &dummy2

	return pg
}

func TestSchemePageGetter_GetPage(t *testing.T) {
	pg := setup()
	ab := make(Aborter)
	url := "dummy1://whatever"
	page, err := pg.GetPage(url, ab)
	if err != nil {
		t.Fatal(err)
	}
	dpg := pg.Children["dummy1"].(*DummyPageGetter)
	expected_page := dpg.PageData
	if !reflect.DeepEqual(page, expected_page) {
		t.Error("page != expected_page")
		t.Log(page)
		t.Log(expected_page)
	}

	if dpg.LastCallUrl != url {
		t.Error("Called with wrong URL")
		t.Logf("Expected: %s", url)
		t.Logf("Got: %s", dpg.LastCallUrl)
	}
	if dpg.LastCallAborter != ab {
		t.Error("Called with wrong Aborter")
		t.Logf("Expected: %r", ab)
		t.Logf("Got: %r", dpg.LastCallAborter)
	}
}

func TestSchemePageGetter_GetPage_ParseFail(t *testing.T) {
	pg := setup()
	_, err := pg.GetPage("%", nil)
	if err == nil {
		t.Fatal("Should have failed on unparseable URL")
	}
}

func TestSchemePageGetter_GetPage_SubFails(t *testing.T) {
	pg := setup()
	var ab Aborter // Heck, let's try this with nil!
	url := "dummy2://whatever"
	_, err := pg.GetPage(url, ab)

	expected := "DUMMY ERROR"
	if err == nil {
		t.Fatal("Should have failed, but didn't")
	}
	if err.Error() != expected {
		t.Fatalf("Expected %s, got %s", expected, err)
	}

	dpg := pg.Children["dummy2"].(*DummyPageGetter)
	if dpg.LastCallUrl != url {
		t.Error("Called with wrong URL")
		t.Logf("Expected: %s", url)
		t.Logf("Got: %s", dpg.LastCallUrl)
	}
	if dpg.LastCallAborter != ab {
		t.Error("Called with wrong Aborter")
		t.Logf("Expected: %r", ab)
		t.Logf("Got: %r", dpg.LastCallAborter)
	}
}

func TestSchemePageGetter_GetPage_UnregisteredScheme(t *testing.T) {
	pg := setup()
	_, err := pg.GetPage("foo://", nil)

	expected := "No PageGetter registered for scheme 'foo'"
	if err == nil {
		t.Fatal("Should have failed, but didn't")
	}
	if err.Error() != expected {
		t.Fatalf("Expected %s, got %s", expected, err)
	}
}

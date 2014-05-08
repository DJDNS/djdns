package server

import (
	"reflect"
	"testing"
)

func setup_apg() AliasPageGetter {
	dpg := DummyPageGetter{
		PageData: Page{
			Url:     "eggs",
			Expired: true,
		},
	}
	dpg.PageData.Data.Meta.Policy = "Stinky"

	pg := NewAliasPageGetter(&dpg)
	pg.Aliases["<FOO>"] = "bar"
	return pg
}

func check_apg(t *testing.T, apg AliasPageGetter, url, exp_url string, ab Aborter) {
	page, err := apg.GetPage(url, ab)
	if err != nil {
		t.Fatal(err)
	}
	dpg := apg.Child.(*DummyPageGetter)
	expected_page := dpg.PageData
	if !reflect.DeepEqual(page, expected_page) {
		t.Error("page != expected_page")
		t.Log(page)
		t.Log(expected_page)
	}

	if dpg.LastCallUrl != exp_url {
		t.Error("Called with wrong URL")
		t.Logf("Expected: %s", exp_url)
		t.Logf("Got: %s", dpg.LastCallUrl)
	}
	if dpg.LastCallAborter != ab {
		t.Error("Called with wrong Aborter")
		t.Logf("Expected: %r", ab)
		t.Logf("Got: %r", dpg.LastCallAborter)
	}
}

func TestAliasPageGetter_GetPage_HasAlias(t *testing.T) {
	pg := setup_apg()
	ab := make(chan interface{})
	url := "<FOO>"
	check_apg(t, pg, url, "bar", ab)
}

func TestAliasPageGetter_GetPage_NoAlias(t *testing.T) {
	pg := setup_apg()
	ab := make(chan interface{})
	url := "Bazzerific"
	check_apg(t, pg, url, url, ab)
}

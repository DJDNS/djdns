package server

import "testing"

type UrlTest struct {
	Input  string
	Output string
}

func TestDejePageGetter_getRouterUrl(t *testing.T) {
	pg := NewDejePageGetter()
	tests := []UrlTest{
		UrlTest{"http://foo/bar", "ws://foo/ws"},
		UrlTest{"http://foo:8080", "ws://foo:8080/ws"},
		UrlTest{"foo.bar.baz:8080", "ws://foo.bar.baz:8080/ws"},
	}
	for _, test := range tests {
		got := pg.getRouterUrl(test.Input)
		if got != test.Output {
			t.Fatalf("Bad result:\ngot: '%s'\nexp: '%s'", got, test.Output)
		}
	}
}

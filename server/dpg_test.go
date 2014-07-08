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
		UrlTest{"foo.bar.baz", "ws://foo.bar.baz/ws"},
		UrlTest{"foo.bar.baz:8080", "ws://foo.bar.baz:8080/ws"},
		UrlTest{"//foo.bar.baz:8080", "ws://foo.bar.baz:8080/ws"},
		UrlTest{"%", "<error>: parse ws://%: hexadecimal escape in host"},
	}
	for _, test := range tests {
		got, err := pg.getRouterUrl(test.Input)
		if err != nil {
			got = "<error>: " + err.Error()
		}
		if got != test.Output {
			t.Errorf("Bad result:\ninput: %s\ngot: '%s'\nexp: '%s'",
				test.Input,
				got,
				test.Output,
			)
		}
	}
}

func TestDejePageGetter_getTopic(t *testing.T) {
	pg := NewDejePageGetter()
	tests := []UrlTest{
		UrlTest{"http://foo/bar", "deje://foo/bar"},
		UrlTest{"http://foo:8080", "deje://foo:8080/"},
		UrlTest{"foo.bar.baz", "deje://foo.bar.baz/"},
		UrlTest{"foo.bar.baz:8080", "deje://foo.bar.baz:8080/"},
		UrlTest{"//foo.bar.baz:8080", "deje://foo.bar.baz:8080/"},
		UrlTest{"%", "<error>: parse ws://%: hexadecimal escape in host"},
	}
	for _, test := range tests {
		got, err := pg.getTopic(test.Input)
		if err != nil {
			got = "<error>: " + err.Error()
		}
		if got != test.Output {
			t.Errorf("Bad result:\ninput: %s\ngot: '%s'\nexp: '%s'",
				test.Input,
				got,
				test.Output,
			)
		}
	}
}

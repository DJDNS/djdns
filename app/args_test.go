package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		Args           []string
		ExpectedConfig ServerConfig
		ExpectedError  string
	}{
		// No arguments, should get default values
		{
			[]string{},
			ServerConfig{
				HostAddress: "0.0.0.0:9953",
				DisplayName: "",
				RootAlias:   "deje://localhost:8080/root",
				ErrorWriter: os.Stderr,
			},
			"",
		},
		// Override multiple arguments
		{
			[]string{"--addr", "foo", "--display-name=bar"},
			ServerConfig{
				HostAddress: "foo",
				DisplayName: "bar",
				RootAlias:   "deje://localhost:8080/root",
				ErrorWriter: os.Stderr,
			},
			"",
		},
		// Error scenario - argument name, but no value
		{
			[]string{"--addr"},
			ServerConfig{},
			"--addr requires argument",
		},
	}
	for _, test := range tests {
		conf, err := Parse(test.Args, false)

		assertError(t, test.ExpectedError, err)
		assert.Equal(t, test.ExpectedConfig, conf)
	}
}

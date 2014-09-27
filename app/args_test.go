package app

import (
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
		if test.ExpectedError == "" {
			assert.NoError(t, err, "ServerConfig failed to instantiate")
		} else {
			if assert.Error(t, err) {
				assert.Equal(t, test.ExpectedError, err.Error())
			}
		}

		assert.Equal(t, test.ExpectedConfig, conf)
	}
}

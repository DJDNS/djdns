package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertError(t *testing.T, expected string, err error) {
	if expected == "" {
		assert.NoError(t, err)
	} else {
		if assert.Error(t, err) {
			assert.Equal(t, expected, err.Error())
		}
	}
}

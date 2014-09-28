package server

import (
	"errors"
	"testing"
	"time"

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

// Useful for testing timeouts
type SlowPageGetter time.Duration

func (pg SlowPageGetter) GetPage(url string, ab Aborter) (Page, error) {
	select {
	case <-time.After(time.Duration(pg)):
		return Page{}, errors.New("Sorry dear, what were you saying?")
	case <-ab:
		return Page{}, errors.New("Ran out of time")
	}
}

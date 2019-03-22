package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapOverridesToStructFields(t *testing.T) {
	assert.Equal(t, "", mapOverridesToStructFields("something-random"))
}

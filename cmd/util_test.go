package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleName(t *testing.T) {
	file := "../fixtures/double-name.yml"
	assert := assert.New(t)
	_, err := getKubeResourcesManifest(file)
	assert.NotNil(err)
}

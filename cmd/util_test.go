package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleNameV1(t *testing.T) {
	file := "../fixtures/double-name_v1.yml"
	assert := assert.New(t)
	_, err := getKubeResourcesManifest(file)
	assert.NotNil(err)
}

func TestUnknownResourceV1(t *testing.T) {
	file := "../fixtures/unknown_type_v1.yml"
	assert := assert.New(t)
	_, err := getKubeResourcesManifest(file)

	assert.Nil(err)
}

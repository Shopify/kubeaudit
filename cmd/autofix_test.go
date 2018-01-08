package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFix(t *testing.T) {
	file := "../fixtures/autofix.yml"
	fileFixed := "../fixtures/autofix-fixed.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assert.Equal(correctlyFixedResources, fixedResources)
}

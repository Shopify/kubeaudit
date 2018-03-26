package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-test/deep"
)

func TestCronjob(t *testing.T) {
	file := "../fixtures/cronjob.yml"
	fileFixed := "../fixtures/cronjob-fixed.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assert.Nil(deep.Equal(correctlyFixedResources, fixedResources))
}

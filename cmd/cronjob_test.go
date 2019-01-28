package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-test/deep"
)

func TestCronjobV1(t *testing.T) {
	file := "../fixtures/cronjob_v1beta1.yml"
	fileFixed := "../fixtures/cronjob-fixed_v1beta1.yml"
	assert := assert.New(t)
	resources, _, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources := fix(resources)
	correctlyFixedResources, _, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assert.Nil(deep.Equal(correctlyFixedResources, fixedResources))
}

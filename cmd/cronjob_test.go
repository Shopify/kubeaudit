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
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources, _ := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assert.Nil(deep.Equal(correctlyFixedResources[0].Object, fixedResources[0].Object))
}

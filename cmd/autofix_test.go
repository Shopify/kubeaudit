package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-test/deep"
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
	assert.Nil(deep.Equal(correctlyFixedResources, fixedResources))
}

func TestFixV1Beta1(t *testing.T) {
	file := "../fixtures/autofix_v1beta1.yml"
	fileFixed := "../fixtures/autofix-fixed_v1beta1.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assert.Nil(deep.Equal(correctlyFixedResources, fixedResources))
}

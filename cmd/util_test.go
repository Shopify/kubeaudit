package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownResourceV1(t *testing.T) {
	file := "../fixtures/unknown_type_v1.yml"
	assert := assert.New(t)
	_, err := getKubeResourcesManifest(file)

	assert.Nil(err)
}

func TestUnknownResourceV2(t *testing.T) {
	file := "../fixtures/unknown_type_v1.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Len(resources, 1)
	result, err, warn := newResultFromResource(resources[0])
	assert.Nil(err)
	assert.Nil(result)
	assert.NotNil(warn)
}

func TestUnknownResourceV3(t *testing.T) {
	file := "../fixtures/unknown_type_v1.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Len(resources, 1)
	result, err, warn := newResultFromResourceWithServiceAccountInfo(resources[0])
	assert.Nil(err)
	assert.Nil(result)
	assert.NotNil(warn)
}

func TestCertificateResourceV1(t *testing.T) {
	file := "../fixtures/certificate_unsupported_v1alpha1.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Len(resources, 1)
	assert.False(IsSupportedResourceType(resources[0]))
	assert.True(IsSupportedGroupVersionKind(resources[0]))
	result, err, warn := newResultFromResource(resources[0])
	assert.Nil(err)
	assert.Nil(result)
	assert.NotNil(warn)
}

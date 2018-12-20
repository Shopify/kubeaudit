package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixV1(t *testing.T) {
	file := "../fixtures/autofix_v1.yml"
	fileFixed := "../fixtures/autofix-fixed_v1.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assertEqualWorkloads(assert, correctlyFixedResources, fixedResources)
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
	assertEqualWorkloads(assert, correctlyFixedResources, fixedResources)
}

func TestPreserveComments(t *testing.T) {
	origFile := "../fixtures/autofix_v1.yml"
	expectedFile := "../fixtures/autofix-fixed_v1.yml"
	assert := assert.New(t)

	rootConfig.manifest = origFile

	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	assert.Nil(err)
	fixedResources := fix(resources)

	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix")
	assert.Nil(err)
	defer os.Remove(tmpFile.Name())

	err = writeManifestFile(fixedResources, tmpFile.Name())
	assert.Nil(err)
	fixedYaml, err := mergeYAML(rootConfig.manifest, tmpFile.Name())
	assert.Nil(err)
	expectedYaml, err := ioutil.ReadFile(expectedFile)
	assert.Nil(err)

	assert.Equal(string(fixedYaml), string(expectedYaml))
}

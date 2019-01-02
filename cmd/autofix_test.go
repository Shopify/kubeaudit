package cmd

import (
	"io"
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
	origFilename := "../fixtures/autofix_v1.yml"
	expectedFilename := "../fixtures/autofix-fixed_v1.yml"
	assert := assert.New(t)

	// Copy original yaml to a temp file because autofix modifies the input file
	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix_test")
	assert.Nil(err)
	defer os.Remove(tmpFile.Name())
	origFile, err := os.Open(origFilename)
	assert.Nil(err)
	_, err = io.Copy(tmpFile, origFile)
	assert.Nil(err)
	tmpFile.Close()
	origFile.Close()

	rootConfig.manifest = tmpFile.Name()
	autofix(nil, nil)

	expectedYaml, err := ioutil.ReadFile(expectedFilename)
	assert.Nil(err)
	fixedYaml, err := ioutil.ReadFile(tmpFile.Name())
	assert.Nil(err)

	assert.Equal(string(fixedYaml), string(expectedYaml))
}

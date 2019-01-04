package cmd

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func compareFiles(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}

	s1 := bufio.NewScanner(f1)
	s2 := bufio.NewScanner(f2)

	for s1.Scan() {
		s2.Scan()
		if !bytes.Equal(s1.Bytes(), s2.Bytes()) {
			return false
		}
	}
	return true
}

func assertEqualWorkloads(assert *assert.Assertions, resource1, resource2 []Resource) {
	file1 := "/tmp/dat1"
	file2 := "/tmp/dat2"
	err := WriteToFile(resource1[0], file1, false)
	assert.Nil(err)
	err = WriteToFile(resource2[0], file2, false)
	assert.Nil(err)
	assert.True(compareFiles(file1, file2))
	os.Remove(file1)
	os.Remove(file2)
}

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

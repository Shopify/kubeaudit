package cmd

import (
	"bufio"
	"fmt"
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
		text1 := s1.Text()
		text2 := s2.Text()
		if text1 != text2 {
			fmt.Printf("Files don't match here:\n%v\n%v\n\n", text1, text2)
			return false
		}
	}
	return true
}

func assertEqualWorkloads(assert *assert.Assertions, resource1, resource2 []Resource) {
	tmpfile1, err := WriteToTmpFile(resource1[0])
	assert.Nil(err)
	defer os.Remove(tmpfile1)
	tmpfile2, err := WriteToTmpFile(resource2[0])
	assert.Nil(err)
	defer os.Remove(tmpfile2)
	assert.True(compareFiles(tmpfile1, tmpfile2))
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

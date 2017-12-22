package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetContainers(t *testing.T) {
	assert := assert.New(t)
	obj := NewPod().DeepCopyObject()
	containers := getContainers(obj)
	containers[0].Name = "modified"
	setContainers(obj, containers)
	for _, container := range getContainers(obj) {
		assert.Equal(container.Name, "modified")
	}
}

func TestGetContainers(t *testing.T) {
	assert := assert.New(t)
	obj := NewPod().DeepCopyObject()
	for _, container := range getContainers(obj) {
		assert.Equal(container.Name, "container")
	}
}

func TestWriteToFile(t *testing.T) {
	file := "../fixtures/read_only_root_filesystem_false.yml"
	fileout := "out.yml"
	assert := assert.New(t)
	resource, err := getKubeResourcesManifest(file)
	assert.Equal(1, len(resource))
	assert.Nil(err)
	err = WriteToFile(resource[0], fileout)
	assert.Nil(err)
	resource2, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Equal(1, len(resource2))
	assert.Equal(resource, resource2)
	err = os.Remove(fileout)
	assert.Nil(err)
}

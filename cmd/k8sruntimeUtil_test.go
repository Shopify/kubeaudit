package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

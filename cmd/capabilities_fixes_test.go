package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixCapabilitiesNotDropped(t *testing.T) {
	assert := assert.New(t)
	file := filepath.Join(path, "capabilities_nil.yml")
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Equal(1, len(resources))
	resource := resources[0]
	results := getResults(resources, auditCapabilities)
	assert.Equal(1, len(results))
	result := results[0]
	for _, occurrence := range result.Occurrences {
		resource = fixCapabilityNotDropped(resource, occurrence)
	}
	add := []Capability{}
	drop := []Capability{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		for _, cap := range drop {
			assert.Contains(container.SecurityContext.Capabilities.Drop, cap)
		}
	}
}

func TestFixCapabilityAdded(t *testing.T) {
	assert := assert.New(t)
	file := filepath.Join(path, "capabilities_some_allowed.yml")
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Equal(1, len(resources))
	resource := resources[0]
	results := getResults(resources, auditCapabilities)
	assert.Equal(1, len(results))
	result := results[0]
	for _, occurrence := range result.Occurrences {
		resource = fixCapabilityAdded(resource, occurrence)
	}
	add := []Capability{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
	}
}

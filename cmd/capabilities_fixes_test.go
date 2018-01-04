package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixCapabilitiesNIL(t *testing.T) {
	assert := assert.New(t)
	file := filepath.Join(path, "capabilities_nil.yml")
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Equal(1, len(resources))
	resource := fixCapabilitiesNIL(resources[0])
	var add []Capability
	drop := []Capability{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assert.Equal(drop, container.SecurityContext.Capabilities.Drop)
	}
}

//func TestFixCapabilitiesNotDropeed(t *testing.T) {
//	assert := assert.New(t)
//	file := filepath.Join(path, "capabilities_nil.yml")
//	resources, err := getKubeResourcesManifest(file)
//	assert.Nil(err)
//	assert.Equal(1, len(resources))
//	resource := fixCapabilityNotDropped(resources[0], occurrence)
//	for _, container := range getContainers(resource) {
//		assert.Equal(add, container.SecurityContext.Capabilities.Add)
//		assert.Equal(drop, container.SecurityContext.Capabilities.Drop)
//	}
//}

//func TestFixCapabilityAdded(t *testing.T) {
//fixCapabilityAdded(resource k8sRuntime.Object, occurrence Occurrence) k8sRuntime.Object {
//	runTest(t, "capabilities_misconfigured_allow.yml", auditCapabilities, ErrorMisconfiguredKubeauditAllow)
//}

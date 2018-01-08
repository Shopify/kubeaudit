package cmd

import "testing"

func TestFixCapabilitiesNotDropped(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_nil.yml", fixCapabilityNotDropped)
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
	assert, resource := FixTestSetup(t, "capabilities_some_allowed.yml", fixCapabilityAdded)
	add := []Capability{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
	}
}

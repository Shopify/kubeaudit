package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAllDropped(assert *assert.Assertions, dropped []Capability, allowed ...[]Capability) {
	to_be_dropped := []Capability{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}
	for _, cap := range to_be_dropped {
		skip := false
		if allowed != nil {
			assert.Equal(1, len(allowed))
			for _, allowed_cap := range allowed[0] {
				if allowed_cap == cap {
					skip = true
				}
			}
		}
		if skip {
			continue
		}
		assert.Contains(dropped, cap)
	}
}

func TestFixCapabilitiesNotDropped(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_nil.yml", auditCapabilities)
	add := []Capability{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitySomeAllowed(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_some_allowed.yml", auditCapabilities)
	add := []Capability{"SYS_TIME"}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop, []Capability{"CHOWN"})
	}
}

func TestFixCapabilitiesNIL(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_nil.yml", auditCapabilities)
	add := []Capability{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesAdded(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_added.yml", auditCapabilities)
	add := []Capability{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesSomeDropped(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_some_dropped.yml", auditCapabilities)
	add := []Capability{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesMisconfiguredAllow(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_misconfigured_allow.yml", auditCapabilities)
	add := []Capability{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

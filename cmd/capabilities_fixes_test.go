package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAllDropped(assert *assert.Assertions, dropped []CapabilityV1, allowed ...[]CapabilityV1) {
	toBeDropped := []CapabilityV1{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}
	for _, cap := range toBeDropped {
		skip := false
		if allowed != nil {
			assert.Equal(1, len(allowed))
			for _, allowedCap := range allowed[0] {
				if allowedCap == cap {
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
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitySomeAllowed(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_some_allowed.yml", auditCapabilities)
	add := []CapabilityV1{"SYS_TIME"}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop, []CapabilityV1{"CHOWN"})
	}
}

func TestFixCapabilitiesNil(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_nil.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesAdded(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_added.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesSomeDropped(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_some_dropped.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesMisconfiguredAllow(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_misconfigured_allow.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		if container.SecurityContext.Capabilities.Add == nil {
			fmt.Println("it is nil!")
		}
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

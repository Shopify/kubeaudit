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

func TestFixCapabilitiesNotDroppedV1Beta2(t *testing.T) {
	rootConfig.auditConfig = ""
	assert, resource := FixTestSetup(t, "capabilities_nil_v1beta2.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitySomeAllowedV1Beta2(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_some_allowed_v1beta2.yml", auditCapabilities)
	add := []CapabilityV1{"SYS_TIME"}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop, []CapabilityV1{"CHOWN"})
	}
}

func TestFixCapabilitiesNilV1Beta2(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_nil_v1beta2.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesAddedV1Beta2(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_added_v1beta2.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesSomeDroppedV1Beta2(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_some_dropped_v1beta2.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitiesMisconfiguredAllowV1Beta2(t *testing.T) {
	assert, resource := FixTestSetup(t, "capabilities_misconfigured_allow_v1beta2.yml", auditCapabilities)
	add := []CapabilityV1{}
	for _, container := range getContainers(resource) {
		if container.SecurityContext.Capabilities.Add == nil {
			fmt.Println("it is nil!")
		}
		assert.Equal(add, container.SecurityContext.Capabilities.Add)
		assertAllDropped(assert, container.SecurityContext.Capabilities.Drop)
	}
}

func TestFixCapabilitySomeAllowedMultiContainersAllContainerLabelsV1Beta2(t *testing.T) {
	assert, resources := FixTestSetupMultipleResources(t, "capabilities_some_allowed_multi_containers_all_container_labels_v1beta2.yml", auditCapabilities)
	add := []CapabilityV1{"SYS_TIME"}
	for _, resource := range resources {
		for _, container := range getContainers(resource) {
			fmt.Println(container.Name)
			assert.Equal(add, container.SecurityContext.Capabilities.Add)
			assertAllDropped(assert, container.SecurityContext.Capabilities.Drop, []CapabilityV1{"CHOWN"})
		}
	}
}

func TestFixCapabilitySomeAllowedMultiContainersSomeContainerLabelsV1Beta2(t *testing.T) {
	assert, resources := FixTestSetupMultipleResources(t, "capabilities_some_allowed_multi_containers_some_container_labels_v1beta2.yml", auditCapabilities)
	for _, resource := range resources {
		for _, container := range getContainers(resource) {
			switch container.Name {
			case "fakeContainerSC":
				add := []CapabilityV1{"SYS_TIME"}
				assert.Equal(add, container.SecurityContext.Capabilities.Add)
				assertAllDropped(assert, container.SecurityContext.Capabilities.Drop, []CapabilityV1{"CHOWN"})
			case "fakeContainerSC2":
				add := []CapabilityV1{}
				assert.Equal(add, container.SecurityContext.Capabilities.Add)
				assertAllDropped(assert, container.SecurityContext.Capabilities.Drop, []CapabilityV1{"CHOWN", "SYS_TIME"})
			}
		}
	}
}

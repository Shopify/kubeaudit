package capabilities

import (
	"sort"

	"github.com/Shopify/kubeaudit/k8stypes"
)

// mergeCapabilities creates an array of all unique capabilities in the custom drop list and container add list
func mergeCapabilities(container *k8stypes.ContainerV1) []string {
	if !SecurityContextOrCapabilities(container) {
		return DefaultDropList
	}

	m := make(map[string]bool)

	for _, cap := range container.SecurityContext.Capabilities.Drop {
		m[string(cap)] = true
	}

	for _, cap := range container.SecurityContext.Capabilities.Add {
		m[string(cap)] = true
	}

	merged := make([]string, 0, len(m))

	// "ALL" and "all" mean the same thing so we don't want to count both
	all := false

	for k := range m {
		if isCapabilityAll(k) {
			all = true
			continue
		}

		merged = append(merged, k)
	}

	if all {
		merged = append(merged, "ALL")
	}

	sort.Strings(merged)

	return merged
}

func isCapabilityAdded(container *k8stypes.ContainerV1, capability string) bool {
	if isCapabilitiesNil(container) {
		return false
	}

	if isCapabilityAll(capability) && isAddAll(container) {
		return true
	}

	added := container.SecurityContext.Capabilities.Add

	return isCapabilityInCapArray(capability, added)
}

// isCapabilityNotDropped returns true if the capability should be dropped (it is in the dropList) but it is not
// dropped (it is not in the container's capability drop list)
func isCapabilityNotDropped(container *k8stypes.ContainerV1, capability string, dropList []string) bool {
	if isCapabilitiesNil(container) {
		return true
	}

	if isDropAll(container) {
		return false
	}

	return isCapabilityInArray(capability, dropList) && !isCapabilityDropped(container, capability)
}

// isCapabilityDropped returns true if the given capability is in the container's capability drop list.
// If the droplist is set to ALL, and the capability being tested is not ALL, it will return false even though
// that capability will be dropped as part of ALL.
func isCapabilityDropped(container *k8stypes.ContainerV1, capability string) bool {
	if isCapabilitiesNil(container) {
		return false
	}

	// Because there are multiple variations of "ALL" we need to explicitly check that both mean "ALL", not just that
	// they are equal strings
	if isCapabilityAll(capability) && isDropAll(container) {
		return true
	}

	dropped := container.SecurityContext.Capabilities.Drop

	return isCapabilityInCapArray(capability, dropped)
}

func isAddAll(container *k8stypes.ContainerV1) bool {
	if isCapabilitiesNil(container) {
		return false
	}

	added := container.SecurityContext.Capabilities.Add
	if len(added) == 0 {
		return false
	}

	return isCapabilityAll(string(added[0]))
}

func isDropAll(container *k8stypes.ContainerV1) bool {
	if isCapabilitiesNil(container) {
		return false
	}

	dropped := container.SecurityContext.Capabilities.Drop
	if len(dropped) == 0 {
		return false
	}

	return isCapabilityAll(string(dropped[0]))
}

func isCapabilityAll(capability string) bool {
	return capability == "ALL" || capability == "all"
}

func isCapabilitiesNil(container *k8stypes.ContainerV1) bool {
	return container.SecurityContext == nil || container.SecurityContext.Capabilities == nil
}

func isCapabilityInCapArray(capability string, capabilities []k8stypes.CapabilityV1) bool {
	if len(capabilities) == 0 {
		return false
	}

	for _, cap := range capabilities {
		if string(cap) == capability {
			return true
		}
	}

	return false
}

func isCapabilityInArray(capability string, capabilities []string) bool {
	if len(capabilities) == 0 {
		return false
	}

	for _, cap := range capabilities {
		if cap == capability {
			return true
		}
	}

	return false
}

package capabilities

import (
	"sort"
	"strings"

	"github.com/Shopify/kubeaudit/pkg/k8s"
)

// uniqueCapabilities creates an array of all unique capabilities in the custom drop list and container add list
func uniqueCapabilities(container *k8s.ContainerV1) []string {
	if !SecurityContextOrCapabilities(container) {
		return DefaultDropList
	}

	m := make(map[string]bool)

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

func isCapabilityAll(capability string) bool {
	return capability == "ALL" || capability == "all"
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

func IsDropAll(container *k8s.ContainerV1) bool {
	for _, cap := range container.SecurityContext.Capabilities.Drop {
		if strings.ToUpper(string(cap)) == "ALL" {
			return true
		}
	}

	return false
}

func IsCapabilityInAddList(container *k8s.ContainerV1, capability string) bool {
	for _, cap := range container.SecurityContext.Capabilities.Add {
		if string(cap) == capability {
			return true
		}
	}

	return false
}

func SecurityContextOrCapabilities(container *k8s.ContainerV1) bool {
	if container.SecurityContext == nil || container.SecurityContext.Capabilities == nil {
		return false
	}

	return true
}

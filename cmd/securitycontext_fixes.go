package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixSecurityContextNil(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if container.SecurityContext == nil {
			container.SecurityContext = &SecurityContextV1{Capabilities: &CapabilitiesV1{Drop: []CapabilityV1{}, Add: []CapabilityV1{}}}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

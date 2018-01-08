package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixSecurityContextNIL(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		if container.SecurityContext == nil {
			container.SecurityContext = &SecurityContext{Capabilities: &Capabilities{Drop: []Capability{}, Add: []Capability{}}}
		}
		if container.SecurityContext.Capabilities == nil {
			container.SecurityContext.Capabilities = &Capabilities{Drop: []Capability{}, Add: []Capability{}}
		}
		if container.SecurityContext.Capabilities.Drop == nil {
			container.SecurityContext.Capabilities.Drop = []Capability{}
		}
		if container.SecurityContext.Capabilities.Add == nil {
			container.SecurityContext.Capabilities.Add = []Capability{}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

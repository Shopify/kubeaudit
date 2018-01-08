package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixSecurityContextNIL(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		if container.SecurityContext == nil {
			container.SecurityContext = &SecurityContext{Capabilities: &Capabilities{Drop: []Capability{}, Add: []Capability{}}}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixAllowPrivilegeEscalation(resource k8sRuntime.Object, occurrence Occurrence) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.AllowPrivilegeEscalation = newFalse()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

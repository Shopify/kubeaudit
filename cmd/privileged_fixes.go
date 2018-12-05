package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixPrivileged(resource k8sRuntime.Object, occurrence Occurrence) k8sRuntime.Object {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.Privileged = newFalse()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

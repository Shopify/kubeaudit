package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixPrivileged(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext.Privileged = newFalse()
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

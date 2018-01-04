package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixCapabilitiesNIL(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext.Capabilities = &Capabilities{
			Drop: []Capability{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER",
				"FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP",
				"SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"},
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixCapabilityNotDropped(resource k8sRuntime.Object, occurrence Occurrence) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext.Capabilities.Drop = append(container.SecurityContext.Capabilities.Drop, Capability(occurrence.metadata["CapName"]))
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixCapabilityAdded(resource k8sRuntime.Object, occurrence Occurrence) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		var add []Capability
		for _, cap := range container.SecurityContext.Capabilities.Add {
			if string(cap) != occurrence.metadata["CapName"] {
				add = append(add, cap)
			}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

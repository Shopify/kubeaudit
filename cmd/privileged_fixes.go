package cmd

func fixPrivileged(result *Result, resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if labelExists, _ := getContainerOverrideLabelReason(result, container, "allow-privileged"); occurrence.container == container.Name && !labelExists {
			container.SecurityContext.Privileged = newFalse()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

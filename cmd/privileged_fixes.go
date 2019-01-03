package cmd

func fixPrivileged(resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.Privileged = newFalse()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

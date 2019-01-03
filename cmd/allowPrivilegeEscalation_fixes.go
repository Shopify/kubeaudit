package cmd

func fixAllowPrivilegeEscalation(resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.AllowPrivilegeEscalation = newFalse()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

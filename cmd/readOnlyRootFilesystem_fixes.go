package cmd

func fixReadOnlyRootFilesystem(resource Resource, occurrence Occurrence) Resource {
	var containers []Container
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.ReadOnlyRootFilesystem = newTrue()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

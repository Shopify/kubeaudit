package cmd

func fixReadOnlyRootFilesystem(result *Result, resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if labelExists, _ := getContainerOverrideLabelReason(result, container, "allow-read-only-root-filesystem-false"); occurrence.container == container.Name && !labelExists {
			container.SecurityContext.ReadOnlyRootFilesystem = newTrue()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

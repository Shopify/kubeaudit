package cmd

func fixRunAsNonRoot(result *Result, resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if labelExists, _ := getContainerOverrideLabelReason(result, container, "allow-run-as-root"); occurrence.container == container.Name && !labelExists {
			container.SecurityContext.RunAsNonRoot = newTrue()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

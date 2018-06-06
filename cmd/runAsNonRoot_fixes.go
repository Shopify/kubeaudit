package cmd

func fixRunAsNonRoot(resource Resource, occurrence Occurrence) Resource {
	var containers []Container
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.RunAsNonRoot = newTrue()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

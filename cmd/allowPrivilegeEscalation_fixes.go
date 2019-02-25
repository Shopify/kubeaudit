package cmd

func fixAllowPrivilegeEscalation(result *Result, resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {

		if labelExists, _ := getContainerOverrideLabelReason(result, container, "allow-privilege-escalation"); occurrence.container == container.Name && !labelExists {
			container.SecurityContext.AllowPrivilegeEscalation = newFalse()
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

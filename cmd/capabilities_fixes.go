package cmd

func fixCapabilitiesNil(resource Resource) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if container.SecurityContext.Capabilities == nil {
			container.SecurityContext.Capabilities = &CapabilitiesV1{}
		}
		if container.SecurityContext.Capabilities.Drop == nil {
			container.SecurityContext.Capabilities.Drop = []CapabilityV1{}
		}
		if container.SecurityContext.Capabilities.Add == nil {
			container.SecurityContext.Capabilities.Add = []CapabilityV1{}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixCapabilityNotDropped(resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.Capabilities.Drop = append(container.SecurityContext.Capabilities.Drop, CapabilityV1(occurrence.metadata["CapName"]))
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixCapabilityAdded(resource Resource, occurrence Occurrence) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			add := []CapabilityV1{}
			for _, cap := range container.SecurityContext.Capabilities.Add {
				if string(cap) != occurrence.metadata["CapName"] {
					add = append(add, cap)
				}
			}
			container.SecurityContext.Capabilities.Add = add
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

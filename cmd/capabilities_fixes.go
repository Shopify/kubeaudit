package cmd

func fixCapabilitiesNIL(resource Resource) Resource {
	var containers []Container
	for _, container := range getContainers(resource) {
		if container.SecurityContext.Capabilities == nil {
			container.SecurityContext.Capabilities = &Capabilities{}
		}
		if container.SecurityContext.Capabilities.Drop == nil {
			container.SecurityContext.Capabilities.Drop = []Capability{}
		}
		if container.SecurityContext.Capabilities.Add == nil {
			container.SecurityContext.Capabilities.Add = []Capability{}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixCapabilityNotDropped(resource Resource, occurrence Occurrence) Resource {
	var containers []Container
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			container.SecurityContext.Capabilities.Drop = append(container.SecurityContext.Capabilities.Drop, Capability(occurrence.metadata["CapName"]))
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixCapabilityAdded(resource Resource, occurrence Occurrence) Resource {
	var containers []Container
	for _, container := range getContainers(resource) {
		if occurrence.container == container.Name {
			add := []Capability{}
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

package cmd

func fixSecurityContextNil(resource Resource) Resource {
	var containers []ContainerV1
	for _, container := range getContainers(resource) {
		if container.SecurityContext == nil {
			container.SecurityContext = &SecurityContextV1{Capabilities: &CapabilitiesV1{Drop: []CapabilityV1{}, Add: []CapabilityV1{}}}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

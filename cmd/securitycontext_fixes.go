package cmd

func fixSecurityContextNIL(resource Resource) Resource {
	var containers []Container
	for _, container := range getContainers(resource) {
		if container.SecurityContext == nil {
			container.SecurityContext = &SecurityContext{Capabilities: &Capabilities{Drop: []Capability{}, Add: []Capability{}}}
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

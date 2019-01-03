package cmd

func fixAppArmor(resource Resource) Resource {
	annotations := getPodAnnotations(resource)

	if annotations == nil {
		annotations = make(map[string]string)
	}

	for _, container := range getContainers(resource) {
		containerAnnotation := ContainerAnnotationKeyPrefix + container.Name
		if val, ok := annotations[containerAnnotation]; !ok || badAppArmorProfileName(val) {
			annotations[containerAnnotation] = ProfileRuntimeDefault
		}
	}

	return setPodAnnotations(resource, annotations)
}

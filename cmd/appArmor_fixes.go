package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixAppArmor(resource k8sRuntime.Object) k8sRuntime.Object {
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

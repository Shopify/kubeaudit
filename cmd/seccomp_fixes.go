package cmd

import (
	apiv1 "k8s.io/api/core/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func fixSeccomp(resource k8sRuntime.Object) k8sRuntime.Object {
	annotations := getPodAnnotations(resource)

	if annotations == nil {
		annotations = make(map[string]string)
	}

	podAnnotation := apiv1.SeccompPodAnnotationKey
	podProfile, podOk := annotations[podAnnotation]
	podValid := podOk && !badSeccompProfileName(podProfile)

	// If there is no pod annotation or it is set to a bad value, set it to the default profile.
	if !podValid {
		annotations[podAnnotation] = apiv1.SeccompProfileRuntimeDefault
	}

	// If container annotation is set to invalid profile and
	// 1. pod annotation is set to default profile, then remove container annotation
	// 2. pod annotation is not set to default profile, then set container annotation to default
	for _, container := range getContainers(resource) {
		containerAnnotation := apiv1.SeccompContainerAnnotationKeyPrefix + container.Name
		containerProfile, containerOk := annotations[containerAnnotation]

		if containerOk && badSeccompProfileName(containerProfile) {
			if podValid && podProfile == apiv1.SeccompProfileRuntimeDefault {
				delete(annotations, containerAnnotation)
			} else {
				annotations[containerAnnotation] = apiv1.SeccompProfileRuntimeDefault
			}
		}
	}

	return setPodAnnotations(resource, annotations)
}

package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewTrue returns a pointer to a boolean variable set to true
func NewTrue() *bool {
	b := true
	return &b
}

// NewFalse returns a pointer to a boolean variable set to false
func NewFalse() *bool {
	return new(bool)
}

func GetContainers(resource Resource) []*ContainerV1 {
	podSpec := GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	var containers []*ContainerV1
	for i := range podSpec.Containers {
		containers = append(containers, &podSpec.Containers[i])
	}

	if len(podSpec.InitContainers) > 0 {
		containers = append(containers, GetInitContainers(resource)...)
	}
	return containers
}

func GetInitContainers(resource Resource) []*ContainerV1 {
	podSpec := GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	containers := make([]*ContainerV1, len(podSpec.InitContainers))
	for i := range podSpec.InitContainers {
		containers[i] = &podSpec.InitContainers[i]
	}
	return containers
}

// GetAnnotations returns the annotations at the pod level. If the resource does not have pods, then it returns
// the least-nested annotations
func GetAnnotations(resource Resource) map[string]string {
	objectMeta := GetPodObjectMeta(resource)
	if objectMeta != nil {
		return objectMeta.GetAnnotations()
	}

	return nil
}

// GetLabels returns the labels at the pod level. If the resource does not have pods, then it returns the
// least-nested labels
func GetLabels(resource Resource) map[string]string {
	objectMeta := GetPodObjectMeta(resource)
	if objectMeta != nil {
		return objectMeta.GetLabels()
	}

	return nil
}

// GetObjectMeta returns the highest-level ObjectMeta
func GetObjectMeta(resource Resource) metav1.Object {
	obj, _ := resource.(metav1.ObjectMetaAccessor)
	if obj != nil {
		return obj.GetObjectMeta()
	}
	return nil
}

// GetPodObjectMeta returns the ObjectMeta at the pod level. If the resource does not have pods, then it returns
// the highest-level ObjectMeta
func GetPodObjectMeta(resource Resource) metav1.Object {
	podTemplateSpec := GetPodTemplateSpec(resource)
	if podTemplateSpec != nil {
		return &podTemplateSpec.ObjectMeta
	}

	return GetObjectMeta(resource)
}

// GetPodSpec gets the PodSpec for a resource. Avoid using this function if you need support for Namespace or
// ServiceAccount resources, and write a helper functions in this package instead
func GetPodSpec(resource Resource) *PodSpecV1 {
	podTemplateSpec := GetPodTemplateSpec(resource)
	if podTemplateSpec != nil {
		return &podTemplateSpec.Spec
	}

	switch kubeType := resource.(type) {
	case *PodV1:
		return &kubeType.Spec
	case *NamespaceV1, *ServiceAccountV1:
		return nil
	}

	return nil
}

// GetPodTemplateSpec gets the PodTemplateSpec for a resource. Avoid using this function if you need support for
// Pod, Namespace, or ServiceAccount resources, and write a helper functions in this package instead
func GetPodTemplateSpec(resource Resource) *PodTemplateSpecV1 {
	switch kubeType := resource.(type) {
	case *CronJobV1Beta1:
		return &kubeType.Spec.JobTemplate.Spec.Template
	case *DaemonSetV1:
		return &kubeType.Spec.Template
	case *DeploymentV1:
		return &kubeType.Spec.Template
	case *JobV1:
		return &kubeType.Spec.Template
	case *PodTemplateV1:
		return &kubeType.Template
	case *ReplicationControllerV1:
		return kubeType.Spec.Template
	case *StatefulSetV1:
		return &kubeType.Spec.Template
	case *PodV1, *NamespaceV1:
		return nil
	}

	return nil
}

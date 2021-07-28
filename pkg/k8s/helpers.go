package k8s

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

	containers := make([]*ContainerV1, len(podSpec.Containers))
	for i := range podSpec.Containers {
		containers[i] = &podSpec.Containers[i]
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
func GetObjectMeta(resource Resource) *ObjectMetaV1 {
	switch kubeType := resource.(type) {
	case *CronJobV1Beta1:
		return &kubeType.ObjectMeta
	case *DaemonSetV1:
		return &kubeType.ObjectMeta
	case *DaemonSetV1Beta1:
		return &kubeType.ObjectMeta
	case *DaemonSetV1Beta2:
		return &kubeType.ObjectMeta
	case *DeploymentExtensionsV1Beta1:
		return &kubeType.ObjectMeta
	case *DeploymentV1:
		return &kubeType.ObjectMeta
	case *DeploymentV1Beta1:
		return &kubeType.ObjectMeta
	case *DeploymentV1Beta2:
		return &kubeType.ObjectMeta
	case *PodTemplateV1:
		return &kubeType.ObjectMeta
	case *ReplicationControllerV1:
		return &kubeType.ObjectMeta
	case *StatefulSetV1:
		return &kubeType.ObjectMeta
	case *StatefulSetV1Beta1:
		return &kubeType.ObjectMeta
	case *PodV1:
		return &kubeType.ObjectMeta
	case *NamespaceV1:
		return &kubeType.ObjectMeta
	case *NetworkPolicyV1:
		return &kubeType.ObjectMeta
	case *ServiceAccountV1:
		return &kubeType.ObjectMeta
	case *ServiceV1:
		return &kubeType.ObjectMeta
	}

	return nil
}

// GetPodObjectMeta returns the ObjectMeta at the pod level. If the resource does not have pods, then it returns
// the highest-level ObjectMeta
func GetPodObjectMeta(resource Resource) *ObjectMetaV1 {
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
	case *DaemonSetV1Beta1:
		return &kubeType.Spec.Template
	case *DaemonSetV1Beta2:
		return &kubeType.Spec.Template
	case *DeploymentExtensionsV1Beta1:
		return &kubeType.Spec.Template
	case *DeploymentV1:
		return &kubeType.Spec.Template
	case *DeploymentV1Beta1:
		return &kubeType.Spec.Template
	case *DeploymentV1Beta2:
		return &kubeType.Spec.Template
	case *PodTemplateV1:
		return &kubeType.Template
	case *ReplicationControllerV1:
		return kubeType.Spec.Template
	case *StatefulSetV1:
		return &kubeType.Spec.Template
	case *StatefulSetV1Beta1:
		return &kubeType.Spec.Template
	case *PodV1, *NamespaceV1:
		return nil
	}

	return nil
}

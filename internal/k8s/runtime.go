package k8s

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/Shopify/kubeaudit/k8stypes"
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

func DecodeResource(b []byte) (k8stypes.Resource, error) {
	decoder := codecs.UniversalDeserializer()
	return k8sRuntime.Decode(decoder, b)
}

func EncodeResource(resource k8stypes.Resource) ([]byte, error) {
	info, _ := k8sRuntime.SerializerInfoForMediaType(codecs.SupportedMediaTypes(), "application/yaml")
	groupVersion := schema.GroupVersion{Group: resource.GetObjectKind().GroupVersionKind().Group, Version: resource.GetObjectKind().GroupVersionKind().Version}
	encoder := codecs.EncoderForVersion(info.Serializer, groupVersion)
	return k8sRuntime.Encode(encoder, resource)
}

func GetContainers(resource k8stypes.Resource) []*k8stypes.ContainerV1 {
	podSpec := GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	containers := make([]*k8stypes.ContainerV1, len(podSpec.Containers))
	for i := range podSpec.Containers {
		containers[i] = &podSpec.Containers[i]
	}
	return containers
}

// GetAnnotations returns the annotations at the pod level. If the resource does not have pods, then it returns
// the least-nested annotations
func GetAnnotations(resource k8stypes.Resource) map[string]string {
	objectMeta := GetPodObjectMeta(resource)
	if objectMeta != nil {
		return objectMeta.GetAnnotations()
	}

	return nil
}

// GetLabels returns the labels at the pod level. If the resource does not have pods, then it returns the
// least-nested labels
func GetLabels(resource k8stypes.Resource) map[string]string {
	objectMeta := GetPodObjectMeta(resource)
	if objectMeta != nil {
		return objectMeta.GetLabels()
	}

	return nil
}

// GetObjectMeta returns the highest-level ObjectMeta
func GetObjectMeta(resource k8stypes.Resource) *metav1.ObjectMeta {
	switch kubeType := resource.(type) {
	case *k8stypes.CronJobV1Beta1:
		return &kubeType.ObjectMeta
	case *k8stypes.DaemonSetV1:
		return &kubeType.ObjectMeta
	case *k8stypes.DaemonSetV1Beta1:
		return &kubeType.ObjectMeta
	case *k8stypes.DaemonSetV1Beta2:
		return &kubeType.ObjectMeta
	case *k8stypes.DeploymentExtensionsV1Beta1:
		return &kubeType.ObjectMeta
	case *k8stypes.DeploymentV1:
		return &kubeType.ObjectMeta
	case *k8stypes.DeploymentV1Beta1:
		return &kubeType.ObjectMeta
	case *k8stypes.DeploymentV1Beta2:
		return &kubeType.ObjectMeta
	case *k8stypes.PodTemplateV1:
		return &kubeType.ObjectMeta
	case *k8stypes.ReplicationControllerV1:
		return &kubeType.ObjectMeta
	case *k8stypes.StatefulSetV1:
		return &kubeType.ObjectMeta
	case *k8stypes.StatefulSetV1Beta1:
		return &kubeType.ObjectMeta
	case *k8stypes.PodV1:
		return &kubeType.ObjectMeta
	case *k8stypes.NamespaceV1:
		return &kubeType.ObjectMeta
	case *k8stypes.NetworkPolicyV1:
		return &kubeType.ObjectMeta
	}

	return nil
}

// GetPodObjectMeta returns the ObjectMeta at the pod level. If the resource does not have pods, then it returns
// the highest-level ObjectMeta
func GetPodObjectMeta(resource k8stypes.Resource) *metav1.ObjectMeta {
	podTemplateSpec := GetPodTemplateSpec(resource)
	if podTemplateSpec != nil {
		return &podTemplateSpec.ObjectMeta
	}

	return GetObjectMeta(resource)
}

// GetPodSpec gets the PodSpec for a resource. Avoid using this function if you need support for Namespace resources,
// and write a helper functions in this package instead
func GetPodSpec(resource k8stypes.Resource) *k8stypes.PodSpecV1 {
	podTemplateSpec := GetPodTemplateSpec(resource)
	if podTemplateSpec != nil {
		return &podTemplateSpec.Spec
	}

	switch kubeType := resource.(type) {
	case *k8stypes.PodV1:
		return &kubeType.Spec
	}

	// Namespace
	return nil
}

// GetPodTemplateSpec gets the PodTemplateSpec for a resource. Avoid using this function if you need support for
// Pod or Namespace resources, and write a helper functions in this package instead
func GetPodTemplateSpec(resource k8stypes.Resource) *v1.PodTemplateSpec {
	switch kubeType := resource.(type) {
	case *k8stypes.CronJobV1Beta1:
		return &kubeType.Spec.JobTemplate.Spec.Template
	case *k8stypes.DaemonSetV1:
		return &kubeType.Spec.Template
	case *k8stypes.DaemonSetV1Beta1:
		return &kubeType.Spec.Template
	case *k8stypes.DaemonSetV1Beta2:
		return &kubeType.Spec.Template
	case *k8stypes.DeploymentExtensionsV1Beta1:
		return &kubeType.Spec.Template
	case *k8stypes.DeploymentV1:
		return &kubeType.Spec.Template
	case *k8stypes.DeploymentV1Beta1:
		return &kubeType.Spec.Template
	case *k8stypes.DeploymentV1Beta2:
		return &kubeType.Spec.Template
	case *k8stypes.PodTemplateV1:
		return &kubeType.Template
	case *k8stypes.ReplicationControllerV1:
		return kubeType.Spec.Template
	case *k8stypes.StatefulSetV1:
		return &kubeType.Spec.Template
	case *k8stypes.StatefulSetV1Beta1:
		return &kubeType.Spec.Template
	}

	// Pod, Namespace
	return nil
}

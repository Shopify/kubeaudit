package cmd

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CapabilitiesV1 is a type alias for the v1 version of the k8s API.
type CapabilitiesV1 = apiv1.Capabilities

// CapabilityV1 is a type alias for the v1 version of the k8s API.
type CapabilityV1 = apiv1.Capability

// ContainerV1 is a type alias for the v1 version of the k8s API.
type ContainerV1 = apiv1.Container

// CronJobV1Beta1 is a type alias for the v1beta1 version of the k8s batch API.
type CronJobV1Beta1 = batchv1beta1.CronJob

// DaemonSetListV1 is a type alias for the v1 version of the k8s apps API.
type DaemonSetListV1 = appsv1.DaemonSetList

// DaemonSetListV1Beta1 is a type alias for the v1beta1 version of the k8s extensions API.
type DaemonSetListV1Beta1 = extensionsv1beta1.DaemonSetList

// DaemonSetV1 is a type alias for the v1 version of the k8s API.
type DaemonSetV1 = appsv1.DaemonSet

// DaemonSetV1Beta1 is a type alias for the v1beta1 version of the k8s extensions API.
type DaemonSetV1Beta1 = extensionsv1beta1.DaemonSet

// DeploymentExtensionsV1Beta1 is a type alias for the v1beta1 version of the k8s extensions API.
type DeploymentExtensionsV1Beta1 = extensionsv1beta1.Deployment

// DeploymentListV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentListV1 = appsv1.DeploymentList

// DeploymentListV1Beta1 is a type alias for the v1beta1 version of the k8s apps API.
type DeploymentListV1Beta1 = appsv1beta1.DeploymentList

// DeploymentV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentV1 = appsv1.Deployment

// DeploymentV1Beta1 is a type alias for the v1beta1 version of the k8s apps API.
type DeploymentV1Beta1 = appsv1beta1.Deployment

// DeploymentV1Beta2 is a type alias for the v1beta2 version of the k8s apps API.
type DeploymentV1Beta2 = appsv1beta2.Deployment

// ListOptionsV1 is a type alias for the v1 version of the k8s meta API.
type ListOptionsV1 = metav1.ListOptions

// NamespaceListV1 is a type alias for the v1 version of the k8s API.
type NamespaceListV1 = apiv1.NamespaceList

// NetworkPolicyListV1 is a type alias for the v1 version of the k8s networking API.
type NetworkPolicyListV1 = networkingv1.NetworkPolicyList

// NetworkPolicyV1 is a type alias for the v1 version of the k8s API.
type NetworkPolicyV1 = networkingv1.NetworkPolicy

// ObjectMetaV1 is a type alias for the v1 version of the k8s API.
type ObjectMetaV1 = metav1.ObjectMeta

// PodListV1 is a type alias for the v1 version of the k8s API.
type PodListV1 = apiv1.PodList

// PodSpecV1 is a type alias for the v1 version of the k8s API.
type PodSpecV1 = apiv1.PodSpec

// PodV1 is a type alias for the v1 version of the k8s API.
type PodV1 = apiv1.Pod

// ReplicationControllerListV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerListV1 = apiv1.ReplicationControllerList

// ReplicationControllerV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerV1 = apiv1.ReplicationController

// SecurityContextV1 is a type alias for the v1 version of the k8s API.
type SecurityContextV1 = apiv1.SecurityContext

// StatefulSetListV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetListV1 = appsv1.StatefulSetList

// StatefulSetListV1Beta1 is a type alias for the v1beta1 version of the k8s apps API.
type StatefulSetListV1Beta1 = appsv1beta1.StatefulSetList

// StatefulSetV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetV1 = appsv1.StatefulSet

// StatefulSetV1Beta1 is a type alias for the v1beta1 version of the k8s API.
type StatefulSetV1Beta1 = appsv1beta1.StatefulSet

// Metadata holds metadata for a potential security issue.
type Metadata = map[string]string

// IsSupportedResourceType returns true if obj is a supported Kubernetes resource type
func IsSupportedResourceType(obj runtime.Object) bool {
	switch obj.(type) {
	case *CronJobV1Beta1, *DaemonSetV1Beta1, *NetworkPolicyV1, *PodV1, *ReplicationControllerV1, *StatefulSetV1Beta1,
		*DaemonSetListV1Beta1, *DeploymentListV1Beta1, *NamespaceListV1, *NetworkPolicyListV1, *PodListV1, *ReplicationControllerListV1,
		*StatefulSetListV1Beta1, *DeploymentV1Beta1, *DeploymentV1Beta2, *DeploymentExtensionsV1Beta1,
		*DeploymentListV1, *DeploymentV1, *DaemonSetV1, *StatefulSetListV1, *DaemonSetListV1, *StatefulSetV1:
		return true
	default:
		return false
	}
}

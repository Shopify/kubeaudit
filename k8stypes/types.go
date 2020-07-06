package k8stypes

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

// CapabilitiesV1 is a type alias for the v1 version of the k8s API.
type CapabilitiesV1 = apiv1.Capabilities

// CapabilityV1 is a type alias for the v1 version of the k8s API.
type CapabilityV1 = apiv1.Capability

// ContainerV1 is a type alias for the v1 version of the k8s API.
type ContainerV1 = apiv1.Container

// CronJobListV1Beta1 is a type alias for the v1beta1 version of the k8s batch API.
type CronJobListV1Beta1 = batchv1beta1.CronJobList

// CronJobV1Beta1 is a type alias for the v1beta1 version of the k8s batch API.
type CronJobV1Beta1 = batchv1beta1.CronJob

// CronJobSpecV1Beta1 is a type alias for the v1beta1 version of the k8s batch API.
type CronJobSpecV1Beta1 = batchv1beta1.CronJobSpec

// DaemonSetListV1 is a type alias for the v1 version of the k8s apps API.
type DaemonSetListV1 = appsv1.DaemonSetList

// DaemonSetSpecV1 is a type alias for the v1 version of the k8s apps API.
type DaemonSetSpecV1 = appsv1.DaemonSetSpec

// DaemonSetV1 is a type alias for the v1 version of the k8s API.
type DaemonSetV1 = appsv1.DaemonSet

// DaemonSetV1Beta1 is a type alias for the v1beta1 version of the k8s extensions API.
type DaemonSetV1Beta1 = extensionsv1beta1.DaemonSet

// DaemonSetV1Beta2 is a type alias for the v1beta2 version of the k8s extensions API.
type DaemonSetV1Beta2 = appsv1beta2.DaemonSet

// DeploymentExtensionsV1Beta1 is a type alias for the v1beta1 version of the k8s extensions API.
type DeploymentExtensionsV1Beta1 = extensionsv1beta1.Deployment

// DeploymentListV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentListV1 = appsv1.DeploymentList

// DeploymentSpecV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentSpecV1 = appsv1.DeploymentSpec

// DeploymentV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentV1 = appsv1.Deployment

// DeploymentV1Beta1 is a type alias for the v1beta1 version of the k8s apps API.
type DeploymentV1Beta1 = appsv1beta1.Deployment

// DeploymentV1Beta2 is a type alias for the v1beta2 version of the k8s apps API.
type DeploymentV1Beta2 = appsv1beta2.Deployment

// JobTemplateSpecV1Beta1 is a type alias for the v1beta1 version of the k8s batch API.
type JobTemplateSpecV1Beta1 = batchv1beta1.JobTemplateSpec

// JobSpecV1 is a type alias for the v1 version of the k8s batch API.
type JobSpecV1 = batchv1.JobSpec

// ListOptionsV1 is a type alias for the v1 version of the k8s meta API.
type ListOptionsV1 = metav1.ListOptions

// NamespaceV1 is a type alias for the v1 version of the k8s API.
type NamespaceV1 = apiv1.Namespace

// NamespaceSpecV1 is a type alias for the v1 version of the k8s API.
type NamespaceSpecV1 = apiv1.NamespaceSpec

// NamespaceListV1 is a type alias for the v1 version of the k8s API.
type NamespaceListV1 = apiv1.NamespaceList

// NetworkPolicyListV1 is a type alias for the v1 version of the k8s networking API.
type NetworkPolicyListV1 = networkingv1.NetworkPolicyList

// NetworkPolicySpecV1 is a type alias for the v1 version of the k8s networking API.
type NetworkPolicySpecV1 = networkingv1.NetworkPolicySpec

// NetworkPolicyV1 is a type alias for the v1 version of the k8s networking API.
type NetworkPolicyV1 = networkingv1.NetworkPolicy

// ObjectMetaV1 is a type alias for the v1 version of the k8s meta API.
type ObjectMetaV1 = metav1.ObjectMeta

// PodListV1 is a type alias for the v1 version of the k8s API.
type PodListV1 = apiv1.PodList

// PodSpecV1 is a type alias for the v1 version of the k8s API.
type PodSpecV1 = apiv1.PodSpec

// PodTemplateSpecV1 is a type alias for the v1 version of the k8s API.
type PodTemplateSpecV1 = apiv1.PodTemplateSpec

// PodTemplateListV1 is a type alias for the v1 version of the k8s API.
type PodTemplateListV1 = apiv1.PodTemplateList

// PodTemplateV1 is a type alias for the v1 version of the k8s API.
type PodTemplateV1 = apiv1.PodTemplate

// PodV1 is a type alias for the v1 version of the k8s API.
type PodV1 = apiv1.Pod

// PolicyTypeV1 is a type alias for the v1 version of the k8s networking API.
type PolicyTypeV1 = networkingv1.PolicyType

// ReplicationControllerListV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerListV1 = apiv1.ReplicationControllerList

// ReplicationControllerSpecV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerSpecV1 = apiv1.ReplicationControllerSpec

// ReplicationControllerV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerV1 = apiv1.ReplicationController

// Resource is a type alias for a runtime.Object
type Resource k8sRuntime.Object

// SecurityContextV1 is a type alias for the v1 version of the k8s API.
type SecurityContextV1 = apiv1.SecurityContext

// StatefulSetListV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetListV1 = appsv1.StatefulSetList

// StatefulSetSpecV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetSpecV1 = appsv1.StatefulSetSpec

// StatefulSetV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetV1 = appsv1.StatefulSet

// StatefulSetV1Beta1 is a type alias for the v1beta1 version of the k8s API.
type StatefulSetV1Beta1 = appsv1beta1.StatefulSet

// UnsupportedType is a type alias for v1 version of the k8s apps API, this is meant for testing
type UnsupportedType = apiv1.Binding

// IsSupportedResourceType returns true if obj is a supported Kubernetes resource type
func IsSupportedResourceType(obj Resource) bool {
	switch obj.(type) {
	case *CronJobV1Beta1,
		*DaemonSetV1, *DaemonSetV1Beta1, *DaemonSetV1Beta2,
		*DeploymentExtensionsV1Beta1, *DeploymentV1, *DeploymentV1Beta1, *DeploymentV1Beta2,
		*NamespaceV1,
		*NetworkPolicyV1,
		*PodV1,
		*PodTemplateV1,
		*ReplicationControllerV1,
		*StatefulSetV1, *StatefulSetV1Beta1:
		return true
	default:
		return false
	}
}

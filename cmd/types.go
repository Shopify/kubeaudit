package cmd

import (
	v1beta1 "k8s.io/api/apps/v1beta1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CronJob represents the same type from the v1beta1 version of the k8s API.
type CronJob = batchv1beta1.CronJob

// DaemonSet represents the same type from the v1beta1 version of the k8s API.
type DaemonSet = extensionsv1beta1.DaemonSet

// NetworkPolicy represents the same type from the v1 version of the k8s API.
type NetworkPolicy = networking.NetworkPolicy

// Pod represents the same type from the v1 version of the k8s API.
type Pod = apiv1.Pod

// ReplicationController represents the same type from the v1 version of the k8s API.
type ReplicationController = apiv1.ReplicationController

// SecurityContext represents the same type from the v1 version of the k8s API.
type SecurityContext = apiv1.SecurityContext

// StatefulSet represents the same type from the v1beta1 version of the k8s API.
type StatefulSet = v1beta1.StatefulSet

// ObjectMeta represents the same type from the v1 version of the k8s API.
type ObjectMeta = metav1.ObjectMeta

// PodSpec represents the same type from the v1 version of the k8s API.
type PodSpec = apiv1.PodSpec

// DaemonSetList represents the same type from the v1beta1 version of the k8s API.
type DaemonSetList = extensionsv1beta1.DaemonSetList

// DeploymentList represents the same type from the v1beta1 version of the k8s API.
type DeploymentList = v1beta1.DeploymentList

// NamespaceList represents the same type from the v1 version of the k8s API.
type NamespaceList = apiv1.NamespaceList

// NetworkPolicyList represents the same type from the v1 version of the k8s API.
type NetworkPolicyList = networking.NetworkPolicyList

// PodList represents the same type from the v1 version of the k8s API.
type PodList = apiv1.PodList

// ReplicationControllerList represents the same type from the v1 version of the k8s API.
type ReplicationControllerList = apiv1.ReplicationControllerList

// StatefulSetList represents the same type from the v1beta1 version of the k8s API.
type StatefulSetList = v1beta1.StatefulSetList

// Capabilities represents the same type from the v1 version of the k8s API.
type Capabilities = apiv1.Capabilities

// Capability represents the same type from the v1 version of the k8s API.
type Capability = apiv1.Capability

// Container represents the same type from the v1 version of the k8s API.
type Container = apiv1.Container

// ListOptions represents the same type from the v1 version of the k8s API.
type ListOptions = metav1.ListOptions

// DeploymentV1Beta1 represents the same type from the v1beta1 version of the k8s API.
type DeploymentV1Beta1 = v1beta1.Deployment

// DeploymentV1Beta2 represents the same type from the v1beta2 version of the k8s API.
type DeploymentV1Beta2 = v1beta2.Deployment

// DeploymentExtensionsV1Beta1 represents the same type from the v1beta1 version of the k8s API.
type DeploymentExtensionsV1Beta1 = extensionsv1beta1.Deployment

// Metadata holds metadata for a potential security issue.
type Metadata = map[string]string

// IsSupportedResourceType returns true if obj is a supported Kubernetes resource type
func IsSupportedResourceType(obj runtime.Object) bool {
	switch obj.(type) {
	case *CronJob, *DaemonSet, *NetworkPolicy, *Pod, *ReplicationController, *StatefulSet,
		*DaemonSetList, *DeploymentList, *NamespaceList, *NetworkPolicyList, *PodList, *ReplicationControllerList,
		*StatefulSetList, *DeploymentV1Beta1, *DeploymentV1Beta2, *DeploymentExtensionsV1Beta1:
		return true
	default:
		return false
	}
}

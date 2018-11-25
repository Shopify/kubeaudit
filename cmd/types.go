package cmd

import (
	appsv1 "k8s.io/api/apps/v1"
	v1beta1 "k8s.io/api/apps/v1beta1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CronJob is a type alias for the v1beta1 version of the k8s API.
type CronJob = batchv1beta1.CronJob

// DaemonSet is a type alias for the v1beta1 version of the k8s API.
type DaemonSet = extensionsv1beta1.DaemonSet

// DaemonSetV1 is a type alias for the v1 version of the k8s API.
type DaemonSetV1 = appsv1.DaemonSet

// NetworkPolicy is a type alias for the v1 version of the k8s API.
type NetworkPolicy = networking.NetworkPolicy

// Pod is a type alias for the v1 version of the k8s API.
type Pod = apiv1.Pod

// ReplicationController is a type alias for the v1 version of the k8s API.
type ReplicationController = apiv1.ReplicationController

// SecurityContext is a type alias for the v1 version of the k8s API.
type SecurityContext = apiv1.SecurityContext

// StatefulSet is a type alias for the v1beta1 version of the k8s API.
type StatefulSet = v1beta1.StatefulSet

// StatefulSetV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetV1 = appsv1.StatefulSet

// ObjectMeta is a type alias for the v1 version of the k8s API.
type ObjectMeta = metav1.ObjectMeta

// PodSpec is a type alias for the v1 version of the k8s API.
type PodSpec = apiv1.PodSpec

// DaemonSetList is a type alias for the v1beta1 version of the k8s API.
type DaemonSetList = extensionsv1beta1.DaemonSetList

// DaemonSetListV1 is a type alias for the v1 version of the k8s apps API.
type DaemonSetListV1 = appsv1.DaemonSetList

// DeploymentList is a type alias for the v1beta1 version of the k8s API.
type DeploymentList = v1beta1.DeploymentList

// DeploymentListV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentListV1 = appsv1.DeploymentList

// NamespaceList is a type alias for the v1 version of the k8s API.
type NamespaceList = apiv1.NamespaceList

// NetworkPolicyList is a type alias for the v1 version of the k8s API.
type NetworkPolicyList = networking.NetworkPolicyList

// PodList is a type alias for the v1 version of the k8s API.
type PodList = apiv1.PodList

// ReplicationControllerList is a type alias for the v1 version of the k8s API.
type ReplicationControllerList = apiv1.ReplicationControllerList

// StatefulSetList is a type alias for the v1beta1 version of the k8s API.
type StatefulSetList = v1beta1.StatefulSetList

// StatefulSetListV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetListV1 = appsv1.StatefulSetList

// Capabilities is a type alias for the v1 version of the k8s API.
type Capabilities = apiv1.Capabilities

// Capability is a type alias for the v1 version of the k8s API.
type Capability = apiv1.Capability

// Container is a type alias for the v1 version of the k8s API.
type Container = apiv1.Container

// ListOptions is a type alias for the v1 version of the k8s API.
type ListOptions = metav1.ListOptions

// DeploymentV1Beta1 is a type alias for the v1beta1 version of the k8s API.
type DeploymentV1Beta1 = v1beta1.Deployment

// DeploymentV1Beta2 is a type alias for the v1beta2 version of the k8s API.
type DeploymentV1Beta2 = v1beta2.Deployment

// DeploymentV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentV1 = appsv1.Deployment

// DeploymentExtensionsV1Beta1 is a type alias for the v1beta1 version of the k8s API.
type DeploymentExtensionsV1Beta1 = extensionsv1beta1.Deployment

// Metadata holds metadata for a potential security issue.
type Metadata = map[string]string

// IsSupportedResourceType returns true if obj is a supported Kubernetes resource type
func IsSupportedResourceType(obj runtime.Object) bool {
	switch obj.(type) {
	case *CronJob, *DaemonSet, *NetworkPolicy, *Pod, *ReplicationController, *StatefulSet,
		*DaemonSetList, *DeploymentList, *NamespaceList, *NetworkPolicyList, *PodList, *ReplicationControllerList,
		*StatefulSetList, *DeploymentV1Beta1, *DeploymentV1Beta2, *DeploymentExtensionsV1Beta1,
		*DeploymentListV1, *DeploymentV1, *DaemonSetV1, *StatefulSetListV1, *DaemonSetListV1, *StatefulSetV1:
		return true
	default:
		return false
	}
}

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
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
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

// DeploymentV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentV1 = appsv1.Deployment

// DeploymentV1Beta1 is a type alias for the v1beta1 version of the k8s apps API.
type DeploymentV1Beta1 = appsv1beta1.Deployment

// DeploymentV1Beta2 is a type alias for the v1beta2 version of the k8s apps API.
type DeploymentV1Beta2 = appsv1beta2.Deployment

// ListOptionsV1 is a type alias for the v1 version of the k8s meta API.
type ListOptionsV1 = metav1.ListOptions

// NamespaceV1 is a type alias for the v1 version of the k8s API.
type NamespaceV1 = apiv1.Namespace

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

// PodTemplateV1 is a type alias for the v1 version of the k8s API.
type PodTemplateV1 = apiv1.PodTemplate

// PodV1 is a type alias for the v1 version of the k8s API.
type PodV1 = apiv1.Pod

// ReplicationControllerListV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerListV1 = apiv1.ReplicationControllerList

// ReplicationControllerV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerV1 = apiv1.ReplicationController

// Resource is a type alias for a runtime.Object.
type Resource k8sRuntime.Object

// SecurityContextV1 is a type alias for the v1 version of the k8s API.
type SecurityContextV1 = apiv1.SecurityContext

// StatefulSetListV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetListV1 = appsv1.StatefulSetList

// StatefulSetV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetV1 = appsv1.StatefulSet

// StatefulSetV1Beta1 is a type alias for the v1beta1 version of the k8s API.
type StatefulSetV1Beta1 = appsv1beta1.StatefulSet

// Metadata holds metadata for a potential security issue.
type Metadata = map[string]string

// UnsupportedType is a type alias for v1 version of the k8s apps API, this is meant for testing
type UnsupportedType = apiv1.Binding

// ResourceTypes is a map of all the Kubernetes workloads kubeaudit can decode
var ResourceTypes = map[string]bool{"ReplicaSet" : true, "Endpoints" : true, "Ingress" : true, "Service" : true,
"ConfigMap" : true, "Secret" : true , "PersistentVolumeClaim" : true, "StorageClass" : true,
"Volume" : true , "VolumeAttachment" : true , "Certificate" : true,
"ControllerRevision" : true, "CustomResourceDefinition" : true, "Event" : true,
"LimitRange" : true, "HorizontalPodAutoscaler" : true, "InitializerConfiguration" : true,
"MutatingWebhookConfiguration" : true, "ValidatingWebhookConfiguration" : true, "PodTemplate" : true,
"PodDisruptionBudget" : true, "PriorityClass" : true,
"PodPreset" : true, "PodSecurityPolicy" : true, "APIService" : true, "Binding" : true,
"CertificateSigningRequest" : true, "ClusterRole" : true,
"ClusterRoleBinding" : true, "ComponentStatus" : true, "LocalSubjectAccessReview" : true, "Node" : true,
"PersistentVolume" : true, "ResourceQuota" : true,
"Role" : true, "RoleBinding" : true,
"SelfSubjectAccessReview" : true, "SelfSubjectRulesReview" : true,
"ServiceAccount" : true, "SubjectAccessReview" : true,
"TokenReview" : true}

// IsSupportedResourceType returns true if obj is a supported Kubernetes resource type
func IsSupportedResourceType(obj Resource) bool {
	switch obj.(type) {
	case *CronJobV1Beta1,
		*DaemonSetListV1, *DaemonSetV1, *DaemonSetV1Beta1, *DaemonSetV1Beta2,
		*DeploymentExtensionsV1Beta1, *DeploymentV1, *DeploymentV1Beta1, *DeploymentV1Beta2, *DeploymentListV1,
		*NamespaceListV1, *NamespaceV1,
		*NetworkPolicyListV1, *NetworkPolicyV1,
		*PodListV1, *PodV1, *PodTemplateV1,
		*ReplicationControllerListV1, *ReplicationControllerV1,
		*StatefulSetListV1, *StatefulSetV1, *StatefulSetV1Beta1:
		return true
	default:
		return false
	}
}

// IsSupportedGroupVersionKind returns false if resource is of Supported Kind but not of supported Group Version Kind
func IsSupportedGroupVersionKind(obj Resource) bool {
	if ( IsSupportedResourceType(obj) || ResourceTypes[obj.GetObjectKind().GroupVersionKind().Kind] ) {
		return true
	} else {
		return false
	}
}

// IsNamespaceType returns true if obj is of NamespaceV1 type
func IsNamespaceType(obj Resource) bool {
	switch obj.(type) {
	case *NamespaceV1:
		return true
	default:
		return false
	}
}

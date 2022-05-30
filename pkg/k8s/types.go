package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
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

// CronJobSpecV1Beta1 is a type alias for the v1beta1 version of the k8s batch API.
type CronJobSpecV1Beta1 = batchv1beta1.CronJobSpec

// DaemonSetSpecV1 is a type alias for the v1 version of the k8s apps API.
type DaemonSetSpecV1 = appsv1.DaemonSetSpec

// DaemonSetV1 is a type alias for the v1 version of the k8s API.
type DaemonSetV1 = appsv1.DaemonSet

// DeploymentSpecV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentSpecV1 = appsv1.DeploymentSpec

// DeploymentV1 is a type alias for the v1 version of the k8s apps API.
type DeploymentV1 = appsv1.Deployment

// JobTemplateSpecV1Beta1 is a type alias for the v1beta1 version of the k8s batch API.
type JobTemplateSpecV1Beta1 = batchv1beta1.JobTemplateSpec

// JobSpecV1 is a type alias for the v1 version of the k8s batch API.
type JobSpecV1 = batchv1.JobSpec

// JobV1 is a type alias for the v1 version of the k8s batch API.
type JobV1 = batchv1.Job

// ListOptionsV1 is a type alias for the v1 version of the k8s meta API.
type ListOptionsV1 = metav1.ListOptions

// NamespaceV1 is a type alias for the v1 version of the k8s API.
type NamespaceV1 = apiv1.Namespace

// NamespaceSpecV1 is a type alias for the v1 version of the k8s API.
type NamespaceSpecV1 = apiv1.NamespaceSpec

// NetworkPolicySpecV1 is a type alias for the v1 version of the k8s networking API.
type NetworkPolicySpecV1 = networkingv1.NetworkPolicySpec

// NetworkPolicyV1 is a type alias for the v1 version of the k8s networking API.
type NetworkPolicyV1 = networkingv1.NetworkPolicy

// ObjectMetaV1 is a type alias for the v1 version of the k8s meta API.
type ObjectMetaV1 = metav1.ObjectMeta

// PodSpecV1 is a type alias for the v1 version of the k8s API.
type PodSpecV1 = apiv1.PodSpec

// PodTemplateSpecV1 is a type alias for the v1 version of the k8s API.
type PodTemplateSpecV1 = apiv1.PodTemplateSpec

// PodTemplateV1 is a type alias for the v1 version of the k8s API.
type PodTemplateV1 = apiv1.PodTemplate

// PodV1 is a type alias for the v1 version of the k8s API.
type PodV1 = apiv1.Pod

// PolicyTypeV1 is a type alias for the v1 version of the k8s networking API.
type PolicyTypeV1 = networkingv1.PolicyType

// ReplicationControllerSpecV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerSpecV1 = apiv1.ReplicationControllerSpec

// ReplicationControllerV1 is a type alias for the v1 version of the k8s API.
type ReplicationControllerV1 = apiv1.ReplicationController

// Resource is a type alias for a runtime.Object
type Resource k8sRuntime.Object

// SecurityContextV1 is a type alias for the v1 version of the k8s API.
type SecurityContextV1 = apiv1.SecurityContext

// ServiceAccountV1 is a type alias for the v1 version of the k8s API.
type ServiceAccountV1 = apiv1.ServiceAccount

// ServiceV1 is a type alias for the v1 version of the k8s API.
type ServiceV1 = apiv1.Service

// ServiceV1Spec is a type alias for the v1 version of the k8s API.
type ServiceV1Spec = apiv1.ServiceSpec

// StatefulSetSpecV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetSpecV1 = appsv1.StatefulSetSpec

// StatefulSetV1 is a type alias for the v1 version of the k8s apps API.
type StatefulSetV1 = appsv1.StatefulSet

// TypeMetaV1 is a type alias for the v1 version of the k8s meta API.
type TypeMetaV1 = metav1.TypeMeta

// UnsupportedType is a type alias for v1 version of the k8s apps API, this is meant for testing
type UnsupportedType = apiv1.Binding

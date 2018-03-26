package cmd

import (
	v1beta1 "k8s.io/api/apps/v1beta1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJob = batchv1beta1.CronJob
type DaemonSet = extensionsv1beta1.DaemonSet
type Deployment = v1beta1.Deployment
type NetworkPolicy = networking.NetworkPolicy
type Pod = apiv1.Pod
type ReplicationController = apiv1.ReplicationController
type SecurityContext = apiv1.SecurityContext
type StatefulSet = v1beta1.StatefulSet

type ObjectMeta = metav1.ObjectMeta
type PodSpec = apiv1.PodSpec

type DaemonSetList = extensionsv1beta1.DaemonSetList
type DeploymentList = v1beta1.DeploymentList
type NamespaceList = apiv1.NamespaceList
type NetworkPolicyList = networking.NetworkPolicyList
type PodList = apiv1.PodList
type ReplicationControllerList = apiv1.ReplicationControllerList
type StatefulSetList = v1beta1.StatefulSetList

type Capabilities = apiv1.Capabilities
type Capability = apiv1.Capability
type Container = apiv1.Container
type ListOptions = metav1.ListOptions

type Metadata = map[string]string

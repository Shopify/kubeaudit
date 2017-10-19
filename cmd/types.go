package cmd

import (
	v1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod = apiv1.Pod
type ReplicationController = apiv1.ReplicationController
type DaemonSet = extensionsv1beta1.DaemonSet
type Deployment = v1beta1.Deployment
type StatefulSet = v1beta1.StatefulSet
type NetworkPolicy = networking.NetworkPolicy

type PodList = apiv1.PodList
type ReplicationControllerList = apiv1.ReplicationControllerList
type DaemonSetList = extensionsv1beta1.DaemonSetList
type DeploymentList = v1beta1.DeploymentList
type StatefulSetList = v1beta1.StatefulSetList
type NamespaceList = apiv1.NamespaceList
type NetworkPolicyList = networking.NetworkPolicyList

type Capability = apiv1.Capability
type Container = apiv1.Container
type ListOptions = metav1.ListOptions

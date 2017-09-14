package cmd

import (
	log "github.com/sirupsen/logrus"
	v1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func debugPrint() {
	if rootConfig.verbose {
		buf := make([]byte, 1<<16)
		stacklen := runtime.Stack(buf, true)
		log.Debugf("%s", buf[:stacklen])
	}
}

func convertDeploymentToDeploymentList(deployment v1beta1.Deployment) (deploymentList *v1beta1.DeploymentList) {
	deploymentList = &v1beta1.DeploymentList{
		Items: []v1beta1.Deployment{deployment},
	}
	return
}

type kubeAuditDeployments struct {
	list *v1beta1.DeploymentList
}

type kubeAuditStatefulSets struct {
	list *v1beta1.StatefulSetList
}

type kubeAuditDaemonSets struct {
	list *extensionsv1beta1.DaemonSetList
}

type kubeAuditPods struct {
	list *apiv1.PodList
}

type kubeAuditReplicationControllers struct {
	list *apiv1.ReplicationControllerList
}

type Result struct {
	err         int
	namespace   string
	name        string
	capsAdded   []apiv1.Capability
	img         string
	capsDropped bool
	kubeType    string
	dsa         string
	sa          string
	token       *bool
}

type Items interface {
	Iter() []interface{}
}

func (deployments kubeAuditDeployments) Iter() (it []interface{}) {
	for _, deployment := range deployments.list.Items {
		it = append(it, deployment)
	}
	return
}

func (statefulSets kubeAuditStatefulSets) Iter() (it []interface{}) {
	for _, statefulSet := range statefulSets.list.Items {
		it = append(it, statefulSet)
	}
	return
}

func (daemonSets kubeAuditDaemonSets) Iter() (it []interface{}) {
	for _, daemonSet := range daemonSets.list.Items {
		it = append(it, daemonSet)
	}
	return
}

func (pods kubeAuditPods) Iter() (it []interface{}) {
	count := 0
	for _, pod := range pods.list.Items {
		if rootConfig.allPods {
			if pod.OwnerReferences == nil {
				it = append(it, pod)
				count = count + 1
			}

		} else {
			if pod.OwnerReferences == nil && pod.Status.Phase == "Running" {
				it = append(it, pod)
				count = count + 1
			}
		}
	}
	it = it[:count]
	return
}

func (replicationControllers kubeAuditReplicationControllers) Iter() (it []interface{}) {
	for _, replicationController := range replicationControllers.list.Items {
		it = append(it, replicationController)
	}
	return
}

func containerIter(t interface{}) (containers []apiv1.Container, result *Result) {
	switch kubeType := t.(type) {
	case v1beta1.Deployment:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "deployment",
		}
		return

	case v1beta1.StatefulSet:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "statefulSet",
		}
		return

	case extensionsv1beta1.DaemonSet:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "daemonSet",
		}
		return

	case apiv1.Pod:
		containers = kubeType.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "pod",
		}
		return

	case apiv1.ReplicationController:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "replicationController",
		}
		return

	default:
		return
	}
}

func ServiceAccountIter(t interface{}) (result *Result) {
	switch kubeType := t.(type) {
	case v1beta1.Deployment:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.Template.Spec.ServiceAccountName,
			kubeType:  "deployment",
		}
		return

	case v1beta1.StatefulSet:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.Template.Spec.ServiceAccountName,
			kubeType:  "statefulSet",
		}
		return

	case extensionsv1beta1.DaemonSet:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.Template.Spec.ServiceAccountName,
			kubeType:  "daemonSet",
		}
		return

	case apiv1.Pod:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.ServiceAccountName,
			kubeType:  "pod",
		}
		return

	case apiv1.ReplicationController:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.Template.Spec.ServiceAccountName,
			kubeType:  "replicationController",
		}
		return

	default:
		return

	}
}

package cmd

import (
	"runtime"
	"sync"

	fakeaudit "github.com/Shopify/kubeaudit/fakeaudit"
	log "github.com/sirupsen/logrus"
	v1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
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

func convertDaemonSetToDaemonSetList(daemonSet extensionsv1beta1.DaemonSet) (daemonSetList *extensionsv1beta1.DaemonSetList) {
	daemonSetList = &extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet},
	}
	return

}

func convertPodToPodList(pod apiv1.Pod) (podList *apiv1.PodList) {
	podList = &apiv1.PodList{
		Items: []apiv1.Pod{pod},
	}
	return

}

func convertStatefulSetToStatefulSetList(statefulSet v1beta1.StatefulSet) (statefulSetList *v1beta1.StatefulSetList) {
	statefulSetList = &v1beta1.StatefulSetList{
		Items: []v1beta1.StatefulSet{statefulSet},
	}
	return

}

func convertReplicationControllerToReplicationList(replicationController apiv1.ReplicationController) (replicationControllerList *apiv1.ReplicationControllerList) {
	replicationControllerList = &apiv1.ReplicationControllerList{
		Items: []apiv1.ReplicationController{replicationController},
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
	imgName     string
	capsDropped bool
	kubeType    string
	dsa         string
	sa          string
	token       *bool
	imgTag      string
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

func getKubeResources(config string) (items []Items, err error) {
	resources, read_err := fakeaudit.ReadConfigFile(config)
	if err != nil {
		err = read_err
		return
	}
	for _, resource := range resources {
		switch resource := resource.(type) {
		case *v1beta1.Deployment:
			items = append(items, kubeAuditDeployments{list: convertDeploymentToDeploymentList(*resource)})
		case *v1beta1.StatefulSet:
			items = append(items, kubeAuditStatefulSets{list: convertStatefulSetToStatefulSetList(*resource)})
		case *extensionsv1beta1.DaemonSet:
			items = append(items, kubeAuditDaemonSets{list: convertDaemonSetToDaemonSetList(*resource)})
		case *apiv1.Pod:
			items = append(items, kubeAuditPods{list: convertPodToPodList(*resource)})
		case *apiv1.ReplicationController:
			items = append(items, kubeAuditReplicationControllers{list: convertReplicationControllerToReplicationList(*resource)})
		}
	}
	return
}

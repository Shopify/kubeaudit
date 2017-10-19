package cmd

import (
	"runtime"
	"sync"

	fakeaudit "github.com/Shopify/kubeaudit/fakeaudit"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func debugPrint() {
	if rootConfig.verbose {
		buf := make([]byte, 1<<16)
		stacklen := runtime.Stack(buf, true)
		log.Debugf("%s", buf[:stacklen])
	}
}

func convertDeploymentToDeploymentList(deployment Deployment) (deploymentList *DeploymentList) {
	deploymentList = &DeploymentList{
		Items: []Deployment{deployment},
	}
	return

}

func convertDaemonSetToDaemonSetList(daemonSet DaemonSet) (daemonSetList *DaemonSetList) {
	daemonSetList = &DaemonSetList{
		Items: []DaemonSet{daemonSet},
	}
	return

}

func convertPodToPodList(pod Pod) (podList *PodList) {
	podList = &PodList{
		Items: []Pod{pod},
	}
	return

}

func convertStatefulSetToStatefulSetList(statefulSet StatefulSet) (statefulSetList *StatefulSetList) {
	statefulSetList = &StatefulSetList{
		Items: []StatefulSet{statefulSet},
	}
	return

}

func convertReplicationControllerToReplicationList(replicationController ReplicationController) (replicationControllerList *ReplicationControllerList) {
	replicationControllerList = &ReplicationControllerList{
		Items: []ReplicationController{replicationController},
	}
	return

}

type kubeAuditDeployments struct {
	list *DeploymentList
}

type kubeAuditStatefulSets struct {
	list *StatefulSetList
}

type kubeAuditDaemonSets struct {
	list *DaemonSetList
}

type kubeAuditPods struct {
	list *PodList
}

type kubeAuditReplicationControllers struct {
	list *ReplicationControllerList
}

type Result struct {
	err         int
	namespace   string
	name        string
	capsAdded   []Capability
	imgName     string
	capsDropped []Capability
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

func containerIter(t interface{}) (containers []Container, result *Result) {
	switch kubeType := t.(type) {
	case Deployment:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "deployment",
		}
		return

	case StatefulSet:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "statefulSet",
		}
		return

	case DaemonSet:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "daemonSet",
		}
		return

	case Pod:
		containers = kubeType.Spec.Containers
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			kubeType:  "pod",
		}
		return

	case ReplicationController:
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
	case Deployment:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.Template.Spec.ServiceAccountName,
			kubeType:  "deployment",
		}
		return

	case StatefulSet:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.Template.Spec.ServiceAccountName,
			kubeType:  "statefulSet",
		}
		return

	case DaemonSet:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.Template.Spec.ServiceAccountName,
			kubeType:  "daemonSet",
		}
		return

	case Pod:
		result = &Result{
			name:      kubeType.Name,
			namespace: kubeType.Namespace,
			dsa:       kubeType.Spec.DeprecatedServiceAccount,
			token:     kubeType.Spec.AutomountServiceAccountToken,
			sa:        kubeType.Spec.ServiceAccountName,
			kubeType:  "pod",
		}
		return

	case ReplicationController:
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
		case *Deployment:
			items = append(items, kubeAuditDeployments{list: convertDeploymentToDeploymentList(*resource)})
		case *StatefulSet:
			items = append(items, kubeAuditStatefulSets{list: convertStatefulSetToStatefulSetList(*resource)})
		case *DaemonSet:
			items = append(items, kubeAuditDaemonSets{list: convertDaemonSetToDaemonSetList(*resource)})
		case *Pod:
			items = append(items, kubeAuditPods{list: convertPodToPodList(*resource)})
		case *ReplicationController:
			items = append(items, kubeAuditReplicationControllers{list: convertReplicationControllerToReplicationList(*resource)})
		}
	}
	return
}

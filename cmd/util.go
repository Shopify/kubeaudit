package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func debugPrint() {
	if rootConfig.verbose == "DEBUG" {
		buf := make([]byte, 1<<16)
		stacklen := runtime.Stack(buf, true)
		log.Debugf("%s", buf[:stacklen])
	}
}

func isInRootConfigNamespace(meta metav1.ObjectMeta) (valid bool) {
	return isInNamespace(meta, rootConfig.namespace)
}

func isInNamespace(meta metav1.ObjectMeta, namespace string) (valid bool) {
	return namespace == apiv1.NamespaceAll || namespace == meta.Namespace
}

func convertDeploymentToDeploymentList(deployment Deployment) (deploymentList *DeploymentList) {
	if isInRootConfigNamespace(deployment.ObjectMeta) {
		deploymentList = &DeploymentList{Items: []Deployment{deployment}}
	} else {
		deploymentList = &DeploymentList{Items: []Deployment{}}
	}
	return
}

func convertDaemonSetToDaemonSetList(daemonSet DaemonSet) (daemonSetList *DaemonSetList) {
	if isInRootConfigNamespace(daemonSet.ObjectMeta) {
		daemonSetList = &DaemonSetList{Items: []DaemonSet{daemonSet}}
	} else {
		daemonSetList = &DaemonSetList{Items: []DaemonSet{}}
	}
	return
}

func convertPodToPodList(pod Pod) (podList *PodList) {
	if isInRootConfigNamespace(pod.ObjectMeta) {
		podList = &PodList{Items: []Pod{pod}}
	} else {
		podList = &PodList{Items: []Pod{}}
	}
	return
}

func convertStatefulSetToStatefulSetList(statefulSet StatefulSet) (statefulSetList *StatefulSetList) {
	if isInRootConfigNamespace(statefulSet.ObjectMeta) {
		statefulSetList = &StatefulSetList{Items: []StatefulSet{statefulSet}}
	} else {
		statefulSetList = &StatefulSetList{Items: []StatefulSet{}}
	}
	return
}

func convertReplicationControllerToReplicationList(replicationController ReplicationController) (replicationControllerList *ReplicationControllerList) {
	if isInRootConfigNamespace(replicationController.ObjectMeta) {
		replicationControllerList = &ReplicationControllerList{Items: []ReplicationController{replicationController}}
	} else {
		replicationControllerList = &ReplicationControllerList{Items: []ReplicationController{}}
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
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			KubeType:  "deployment",
		}
		return

	case StatefulSet:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			KubeType:  "statefulSet",
		}
		return

	case DaemonSet:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			KubeType:  "daemonSet",
		}
		return

	case Pod:
		containers = kubeType.Spec.Containers
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			KubeType:  "pod",
		}
		return

	case ReplicationController:
		containers = kubeType.Spec.Template.Spec.Containers
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			KubeType:  "replicationController",
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
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			DSA:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			Token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			SA:        kubeType.Spec.Template.Spec.ServiceAccountName,
			KubeType:  "deployment",
		}
		return

	case StatefulSet:
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			DSA:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			Token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			SA:        kubeType.Spec.Template.Spec.ServiceAccountName,
			KubeType:  "statefulSet",
		}
		return

	case DaemonSet:
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			DSA:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			Token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			SA:        kubeType.Spec.Template.Spec.ServiceAccountName,
			KubeType:  "daemonSet",
		}
		return

	case Pod:
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			DSA:       kubeType.Spec.DeprecatedServiceAccount,
			Token:     kubeType.Spec.AutomountServiceAccountToken,
			SA:        kubeType.Spec.ServiceAccountName,
			KubeType:  "pod",
		}
		return

	case ReplicationController:
		result = &Result{
			Name:      kubeType.Name,
			Namespace: kubeType.Namespace,
			DSA:       kubeType.Spec.Template.Spec.DeprecatedServiceAccount,
			Token:     kubeType.Spec.Template.Spec.AutomountServiceAccountToken,
			SA:        kubeType.Spec.Template.Spec.ServiceAccountName,
			KubeType:  "replicationController",
		}
		return

	default:
		return

	}
}

func getKubeResources(clientset *kubernetes.Clientset) (items []Items) {
	// fetch deployments, statefulsets, daemonsets
	// and pods which do not belong to another abstraction
	deployments := getDeployments(clientset)
	statefulSets := getStatefulSets(clientset)
	daemonSets := getDaemonSets(clientset)
	pods := getPods(clientset)
	replicationControllers := getReplicationControllers(clientset)

	items = append(items, kubeAuditDeployments{deployments})
	items = append(items, kubeAuditStatefulSets{statefulSets})
	items = append(items, kubeAuditDaemonSets{daemonSets})
	items = append(items, kubeAuditPods{pods})
	items = append(items, kubeAuditReplicationControllers{replicationControllers})

	return
}

func getKubeResourcesManifest(config string) (items []Items, err error) {
	resources, read_err := readManifestFile(config)
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

func readManifestFile(filename string) (decoded []k8sRuntime.Object, err error) {
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Error("File not found")
		return
	}
	buf_slice := bytes.Split(buf, []byte("---"))

	decoder := scheme.Codecs.UniversalDeserializer()

	for _, b := range buf_slice {
		obj, _, err := decoder.Decode(b, nil, nil)
		if err == nil && obj != nil {
			decoded = append(decoded, obj)
		}
	}
	return
}

func getResources() (resources []Items, err error) {
	if rootConfig.manifest != "" {
		resources, err = getKubeResourcesManifest(rootConfig.manifest)
	} else {
		if kube, err := kubeClient(rootConfig.kubeConfig); err == nil {
			resources = getKubeResources(kube)
		}
	}
	return
}

func setFormatter() {
	if rootConfig.json {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func checkParams(auditFunc interface{}) (err error) {
	switch auditFunc.(type) {
	case (func(image imgFlags, item Items) (results []Result)):
		if len(imgConfig.img) == 0 {
			return errors.New("Empty image name. Are you missing the image flag?")
		}
		imgConfig.splitImageString()
		if len(imgConfig.tag) == 0 {
			return errors.New("Empty image tag. Are you missing the image tag?")
		}
	}
	return nil
}

func getResults(resources []Items, auditFunc interface{}) []Result {
	var wg sync.WaitGroup
	wg.Add(len(resources))
	resultsChannel := make(chan []Result, 1)
	go func() { resultsChannel <- []Result{} }()

	for _, resource := range resources {
		results := <-resultsChannel
		go func(item Items) {
			switch f := auditFunc.(type) {
			case func(item Items) (results []Result):
				resultsChannel <- append(results, f(item)...)
			case func(image imgFlags, item Items) (results []Result):
				resultsChannel <- append(results, f(imgConfig, item)...)
			case func(limits limitFlags, item Items) (results []Result):
				resultsChannel <- append(results, f(limitConfig, item)...)
			default:
				log.Fatal("Invalid audit function provided")
			}
			wg.Done()
		}(resource)
	}

	wg.Wait()
	close(resultsChannel)

	var results []Result
	for _, result := range <-resultsChannel {
		results = append(results, result)
	}
	return results
}

func runAudit(auditFunc interface{}) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := checkParams(auditFunc); err != nil {
			log.Error(err)
		}
		setFormatter()
		resources, err := getResources()
		if err != nil {
			log.Error(err)
			return
		}
		results := getResults(resources, auditFunc)
		for _, result := range results {
			result.Print()
		}
	}
}

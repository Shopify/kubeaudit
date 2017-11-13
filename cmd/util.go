package cmd

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	Err            int
	Occurrences    []Occurrence
	Namespace      string
	Name           string
	CapsAdded      []Capability
	ImageName      string
	CapsDropped    []Capability
	CapsNotDropped []Capability
	KubeType       string
	DSA            string
	SA             string
	Token          *bool
	ImageTag       string
	CPULimitActual string
	CPULimitMax    string
	MEMLimitActual string
	MEMLimitMax    string
}

func (res Result) Print() {
	for _, occ := range res.Occurrences {
		if occ.kind <= KubeauditLogLevels[rootConfig.verbose] {
			logger := log.WithFields(createFields(res, occ.id))
			switch occ.kind {
			case Debug:
				logger.Debug(occ.message)
			case Info:
				logger.Info(occ.message)
			case Warn:
				logger.Warn(occ.message)
			case Error:
				logger.Error(occ.message)
			}
		}
	}
}

func createFields(res Result, err int) (fields log.Fields) {
	fields = log.Fields{}
	v := reflect.ValueOf(res)
	for _, member := range shouldLog(err) {
		value := v.FieldByName(member)
		if value.IsValid() && value.Interface() != nil && value.Interface() != "" {
			fields[member] = value.Interface()
		}
	}
	return
}

func shouldLog(err int) (members []string) {
	members = []string{"Name", "Namespace", "KubeType"}
	switch err {
	case ErrorCapabilitiesAdded:
		members = append(members, "CapsAdded")
	case ErrorCapabilitiesSomeDropped:
		members = append(members, "CapsNotDropped")
	case ErrorServiceAccountTokenDeprecated:
		members = append(members, "DSA")
		members = append(members, "SA")
	case InfoImageCorrect:
	case ErrorImageTagMissing:
	case ErrorImageTagIncorrect:
		members = append(members, "ImageTag")
		members = append(members, "ImageName")
	case ErrorResourcesLimitsCpuExceeded:
		members = append(members, "CPULimitActual")
		members = append(members, "CPULimitMax")
	case ErrorResourcesLimitsMemoryExceeded:
		members = append(members, "MEMLimitActual")
		members = append(members, "MEMLimitMax")
	}
	return
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
	resources, read_err := readConfigFiles(config)
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

func readConfigFiles(filename string) (decoded []k8sRuntime.Object, err error) {
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

func runAudit(auditFunc interface{}) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		switch auditFunc.(type) {
		case (func(image imgFlags, item Items) (results []Result)):
			if len(imgConfig.img) == 0 {
				log.Error("Empty image name. Are you missing the image flag?")
				return
			}
			imgConfig.splitImageString()
			if len(imgConfig.tag) == 0 {
				log.Error("Empty image tag. Are you missing the image tag?")
				return
			}
		}

		if rootConfig.json {
			log.SetFormatter(&log.JSONFormatter{})
		}
		var resources []Items

		if rootConfig.manifest != "" {
			var err error
			resources, err = getKubeResourcesManifest(rootConfig.manifest)
			if err != nil {
				log.Error(err)
			}
		} else {
			kube, err := kubeClient(rootConfig.kubeConfig)
			if err != nil {
				log.Error(err)
			}
			resources = getKubeResources(kube)
		}

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

		for _, result := range <-resultsChannel {
			result.Print()
		}
	}
}

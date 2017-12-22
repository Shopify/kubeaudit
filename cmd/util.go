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

func getContainers(resource k8sRuntime.Object) (container []Container) {
	switch kubeType := resource.(type) {
	case *DaemonSet:
		container = kubeType.Spec.Template.Spec.Containers
	case *Deployment:
		container = kubeType.Spec.Template.Spec.Containers
	case *Pod:
		container = kubeType.Spec.Containers
	case *ReplicationController:
		container = kubeType.Spec.Template.Spec.Containers
	case *StatefulSet:
		container = kubeType.Spec.Template.Spec.Containers
	}
	return container
}

func newResultFromResource(resource k8sRuntime.Object) (result Result) {
	switch kubeType := resource.(type) {
	case *DaemonSet:
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
		result.KubeType = "daemonSet"
	case *Deployment:
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
		result.KubeType = "deployment"
	case *Pod:
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
		result.KubeType = "pod"
	case *ReplicationController:
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
		result.KubeType = "replicationController"
	case *StatefulSet:
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
		result.KubeType = "statefulSet"
	}
	return
}

func newResultFromResourceWithServiceAccountInfo(resource k8sRuntime.Object) Result {
	result := newResultFromResource(resource)
	switch kubeType := resource.(type) {
	case *DaemonSet:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
	case *Deployment:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
	case *Pod:
		result.DSA = kubeType.Spec.DeprecatedServiceAccount
		result.Token = kubeType.Spec.AutomountServiceAccountToken
		result.SA = kubeType.Spec.ServiceAccountName
	case *ReplicationController:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
	case *StatefulSet:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
	}
	return result
}

func getKubeResources(clientset *kubernetes.Clientset) (resources []k8sRuntime.Object) {
	for _, resource := range getDaemonSets(clientset).Items {
		if isInRootConfigNamespace(resource.ObjectMeta) {
			resources = append(resources, resource.DeepCopyObject())
		}
	}
	for _, resource := range getDeployments(clientset).Items {
		if isInRootConfigNamespace(resource.ObjectMeta) {
			resources = append(resources, resource.DeepCopyObject())
		}
	}
	for _, resource := range getPods(clientset).Items {
		if isInRootConfigNamespace(resource.ObjectMeta) {
			resources = append(resources, resource.DeepCopyObject())
		}
	}
	for _, resource := range getReplicationControllers(clientset).Items {
		if isInRootConfigNamespace(resource.ObjectMeta) {
			resources = append(resources, resource.DeepCopyObject())
		}
	}
	for _, resource := range getStatefulSets(clientset).Items {
		if isInRootConfigNamespace(resource.ObjectMeta) {
			resources = append(resources, resource.DeepCopyObject())
		}
	}
	return
}

func getKubeResourcesManifest(filename string) (decoded []k8sRuntime.Object, err error) {
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

func getResources() (resources []k8sRuntime.Object, err error) {
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
	case (func(image imgFlags, resource k8sRuntime.Object) (results []Result)):
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

func getResults(resources []k8sRuntime.Object, auditFunc interface{}) []Result {
	var wg sync.WaitGroup
	wg.Add(len(resources))
	resultsChannel := make(chan []Result, 1)
	go func() { resultsChannel <- []Result{} }()

	for _, resource := range resources {
		results := <-resultsChannel
		go func(resource k8sRuntime.Object) {
			switch f := auditFunc.(type) {
			case func(resource k8sRuntime.Object) (results []Result):
				resultsChannel <- append(results, f(resource)...)
			case func(image imgFlags, resource k8sRuntime.Object) (results []Result):
				resultsChannel <- append(results, f(imgConfig, resource)...)
			case func(limits limitFlags, resource k8sRuntime.Object) (results []Result):
				resultsChannel <- append(results, f(limitConfig, resource)...)
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

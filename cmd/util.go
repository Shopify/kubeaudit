package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	return new(bool)
}

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

func newResultFromResource(resource k8sRuntime.Object) (result *Result) {
	result = &Result{}

	switch kubeType := resource.(type) {
	case *DaemonSet:
		result.KubeType = "daemonSet"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *DeploymentV1Beta1:
		result.KubeType = "deployment"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *DeploymentV1Beta2:
		result.KubeType = "deployment"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *DeploymentExtensionsV1Beta1:
		result.KubeType = "deployment"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *Pod:
		result.KubeType = "pod"
		result.Labels = kubeType.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *ReplicationController:
		result.KubeType = "replicationController"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *StatefulSet:
		result.KubeType = "statefulSet"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	default:
		return nil
	}
	return
}

func newResultFromResourceWithServiceAccountInfo(resource k8sRuntime.Object) *Result {
	result := newResultFromResource(resource)
	switch kubeType := resource.(type) {
	case *DaemonSet:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *DeploymentV1Beta1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *DeploymentV1Beta2:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *DeploymentExtensionsV1Beta1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *Pod:
		result.DSA = kubeType.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.ServiceAccountName
		result.Token = kubeType.Spec.AutomountServiceAccountToken
	case *ReplicationController:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *StatefulSet:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	default:
		return nil
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

func writeManifestFile(decoded []k8sRuntime.Object, filename string) error {
	var toAppend bool
	for _, decode := range decoded {
		if err := WriteToFile(decode, filename, toAppend); err != nil {
			log.Error(err)
			return err
		}
		toAppend = true
	}
	return nil
}

func containerNamesUniq(resource k8sRuntime.Object) bool {
	names := make(map[string]bool)
	for _, container := range getContainers(resource) {
		if names[container.Name] {
			return false
		}
		names[container.Name] = true
	}
	return true
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
			if !containerNamesUniq(obj) {
				log.Error("Container names are not uniq")
				return nil, errors.New("Container names are not uniq")
			}
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
				name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
				log.Fatal("Invalid audit function provided: ", name)
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
			log.Error("Parameter check failed")
			log.Error(err)
		}
		setFormatter()
		resources, err := getResources()
		if err != nil {
			log.Error("getResources failed")
			log.Error(err)
			return
		}
		results := getResults(resources, auditFunc)
		for _, result := range results {
			result.Print()
		}
	}
}

func mergeAuditFunctions(auditFunctions []interface{}) func(resource k8sRuntime.Object) (results []Result) {
	return func(resource k8sRuntime.Object) (results []Result) {
		for _, function := range auditFunctions {
			for _, result := range getResults([]k8sRuntime.Object{resource}, function) {
				results = append(results, result)
			}
		}
		return results
	}
}

func prettifyReason(reason string) string {
	if strings.ToLower(reason) == "true" {
		return "Unspecified"
	}
	return reason
}

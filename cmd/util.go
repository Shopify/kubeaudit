package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func newResultFromResource(resource Resource) (*Result, error) {
	result := &Result{}
	switch kubeType := resource.(type) {
	case *CronJobV1Beta1:
		result.KubeType = "cronjob"
		result.Labels = kubeType.Spec.JobTemplate.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *DaemonSetV1:
		result.KubeType = "daemonSet"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *DaemonSetV1Beta1:
		result.KubeType = "daemonSet"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *DeploymentExtensionsV1Beta1:
		result.KubeType = "deployment"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *DeploymentV1:
		result.KubeType = "deployment"
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
	case *PodV1:
		result.KubeType = "pod"
		result.Labels = kubeType.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *ReplicationControllerV1:
		result.KubeType = "replicationController"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *StatefulSetV1:
		result.KubeType = "statefulSet"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *StatefulSetV1Beta1:
		result.KubeType = "statefulSet"
		result.Labels = kubeType.Spec.Template.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	case *NamespaceV1:
		result.KubeType = "namespace"
		result.Labels = kubeType.Labels
		result.Name = kubeType.Name
		result.Namespace = kubeType.Namespace
	default:
		return nil, fmt.Errorf("resource type %s not supported", resource.GetObjectKind().GroupVersionKind())
	}
	return result, nil
}

func newResultFromResourceWithServiceAccountInfo(resource Resource) (*Result, error) {
	result, err := newResultFromResource(resource)
	if err != nil {
		return nil, err
	}

	switch kubeType := resource.(type) {
	case *CronJobV1Beta1:
		result.DSA = kubeType.Spec.JobTemplate.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.JobTemplate.Spec.Template.Spec.AutomountServiceAccountToken
	case *DaemonSetV1Beta1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *DaemonSetV1:
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
	case *DeploymentV1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *DeploymentExtensionsV1Beta1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *PodV1:
		result.DSA = kubeType.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.ServiceAccountName
		result.Token = kubeType.Spec.AutomountServiceAccountToken
	case *ReplicationControllerV1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *StatefulSetV1Beta1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *StatefulSetV1:
		result.DSA = kubeType.Spec.Template.Spec.DeprecatedServiceAccount
		result.SA = kubeType.Spec.Template.Spec.ServiceAccountName
		result.Token = kubeType.Spec.Template.Spec.AutomountServiceAccountToken
	case *NamespaceV1:
		// We need to set this here so the audit function will ignore the namespace
		result.Token = newFalse()
	}

	return result, nil
}

func getKubeResources(clientset *kubernetes.Clientset) (resources []Resource) {
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
	for _, resource := range getNamespaces(clientset).Items {
		if isInRootConfigNamespace(resource.ObjectMeta) {
			resources = append(resources, resource.DeepCopyObject())
		}
	}

	return
}

func writeManifestFile(decoded []Resource, filename string) error {
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

func containerNamesUniq(resource Resource) bool {
	names := make(map[string]bool)
	for _, container := range getContainers(resource) {
		if names[container.Name] {
			return false
		}
		names[container.Name] = true
	}
	return true
}

func getKubeResourcesManifest(filename string) (decoded []Resource, err error) {
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Error("File not found")
		return
	}
	bufSlice := bytes.Split(buf, []byte("---"))

	decoder := scheme.Codecs.UniversalDeserializer()

	for _, b := range bufSlice {
		obj, _, err := decoder.Decode(b, nil, nil)
		if err == nil && obj != nil {
			if !IsSupportedResourceType(obj) {
				log.Warnf("Skipping unsupported resource type %s", obj.GetObjectKind().GroupVersionKind())
				continue
			}

			if !containerNamesUniq(obj) {
				log.Error("Container names are not uniq")
				return nil, errors.New("Container names are not uniq")
			}
			decoded = append(decoded, obj)
		}
	}
	return
}

func getResources() (resources []Resource, err error) {
	if rootConfig.manifest != "" {
		resources, err = getKubeResourcesManifest(rootConfig.manifest)
	} else {
		if kube, err := kubeClient(); err == nil {
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
	case (func(image imgFlags, resource Resource) (results []Result)):
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

func getResults(resources []Resource, auditFunc interface{}) []Result {
	var wg sync.WaitGroup
	wg.Add(len(resources))
	resultsChannel := make(chan []Result, 1)
	go func() { resultsChannel <- []Result{} }()

	for _, resource := range resources {
		results := <-resultsChannel
		go func(resource Resource) {
			switch f := auditFunc.(type) {
			case func(resource Resource) (results []Result):
				resultsChannel <- append(results, f(resource)...)
			case func(image imgFlags, resource Resource) (results []Result):
				resultsChannel <- append(results, f(imgConfig, resource)...)
			case func(limits limitFlags, resource Resource) (results []Result):
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

func mergeAuditFunctions(auditFunctions []interface{}) func(resource Resource) (results []Result) {
	return func(resource Resource) (results []Result) {
		for _, function := range auditFunctions {
			for _, result := range getResults([]Resource{resource}, function) {
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

func isCSCWellDefined(podSpec PodSpecV1, container ContainerV1) bool {
	if container.SecurityContext != nil && container.SecurityContext.RunAsNonRoot != nil {
		return true
	}
	if podSpec.SecurityContext == nil || podSpec.SecurityContext.RunAsNonRoot == nil {
		return true
	}
	return false
}

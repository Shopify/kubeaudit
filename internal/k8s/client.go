package k8s

import (
	"errors"
	"fmt"
	"os"

	"github.com/Shopify/kubeaudit/k8stypes"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ErrNoReadableKubeConfig represents any error that prevents the client from opening a kubeconfig file.
var ErrNoReadableKubeConfig = errors.New("unable to open kubeconfig file")

var DefaultClient = k8sClient{}

// Client abstracts the API to allow testing.
type Client interface {
	InClusterConfig() (*rest.Config, error)
}

// k8sClient wraps kubernetes client-go so it can be mocked.
type k8sClient struct{}

// InClusterConfig wraps the client-go method with the same name.
func (kc k8sClient) InClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}

func NewKubeClientLocal(configPath string) (*kubernetes.Clientset, error) {
	if _, err := os.Stat(configPath); err != nil {
		return nil, ErrNoReadableKubeConfig
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}
	kube, err := kubernetes.NewForConfig(config)
	return kube, err
}

func NewKubeClientCluster(client Client) (*kubernetes.Clientset, error) {
	config, err := client.InClusterConfig()
	if err != nil {
		return nil, err
	}
	log.Info("Running inside cluster, using the cluster config")
	kube, err := kubernetes.NewForConfig(config)
	return kube, err
}

func IsRunningInCluster(client Client) bool {
	_, err := client.InClusterConfig()
	return err == nil
}

func GetAllResources(clientset kubernetes.Interface, namespace string) []k8stypes.Resource {
	var resources []k8stypes.Resource

	for _, resource := range GetDaemonSets(clientset, namespace).Items {
		resources = append(resources, resource.DeepCopyObject())
	}
	for _, resource := range GetDeployments(clientset, namespace).Items {
		resources = append(resources, resource.DeepCopyObject())
	}
	for _, resource := range GetPods(clientset, namespace).Items {
		resources = append(resources, resource.DeepCopyObject())
	}
	for _, resource := range GetReplicationControllers(clientset, namespace).Items {
		resources = append(resources, resource.DeepCopyObject())
	}
	for _, resource := range GetStatefulSets(clientset, namespace).Items {
		resources = append(resources, resource.DeepCopyObject())
	}
	for _, resource := range GetNetworkPolicies(clientset, namespace).Items {
		resources = append(resources, resource.DeepCopyObject())
	}
	for _, resource := range GetNamespaces(clientset, namespace).Items {
		resources = append(resources, resource.DeepCopyObject())
	}

	return resources
}

func GetDaemonSets(clientset kubernetes.Interface, namespace string) *k8stypes.DaemonSetListV1 {
	daemonSetClient := clientset.AppsV1().DaemonSets(namespace)
	daemonSets, err := daemonSetClient.List(k8stypes.ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return daemonSets
}

func GetDeployments(clientset kubernetes.Interface, namespace string) *k8stypes.DeploymentListV1 {
	deploymentClient := clientset.AppsV1().Deployments(namespace)
	deployments, err := deploymentClient.List(k8stypes.ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return deployments
}

func GetPods(clientset kubernetes.Interface, namespace string) *k8stypes.PodListV1 {
	podClient := clientset.CoreV1().Pods(namespace)
	pods, err := podClient.List(k8stypes.ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return pods
}

func GetReplicationControllers(clientset kubernetes.Interface, namespace string) *k8stypes.ReplicationControllerListV1 {
	replicationControllerClient := clientset.CoreV1().ReplicationControllers(namespace)
	replicationControllers, err := replicationControllerClient.List(k8stypes.ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return replicationControllers
}

func GetStatefulSets(clientset kubernetes.Interface, namespace string) *k8stypes.StatefulSetListV1 {
	statefulSetClient := clientset.AppsV1().StatefulSets(namespace)
	statefulSets, err := statefulSetClient.List(k8stypes.ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return statefulSets
}

func GetNetworkPolicies(clientset kubernetes.Interface, namespace string) *k8stypes.NetworkPolicyListV1 {
	netPolClient := clientset.NetworkingV1().NetworkPolicies(namespace)
	netPols, err := netPolClient.List(k8stypes.ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return netPols
}

func GetNamespaces(clientset kubernetes.Interface, namespace string) *k8stypes.NamespaceListV1 {
	namespaceClient := clientset.CoreV1().Namespaces()
	listOptions := k8stypes.ListOptionsV1{}

	if namespace != "" {
		// Select only the specified namespace
		listOptions.FieldSelector = fmt.Sprintf("metadata.name=%s", namespace)
	}

	namespaces, err := namespaceClient.List(listOptions)

	if err != nil {
		log.Error(err)
	}

	return namespaces
}

func GetKubernetesVersion(clientset kubernetes.Interface) (*version.Info, error) {
	discoveryClient := clientset.Discovery()
	return discoveryClient.ServerVersion()
}

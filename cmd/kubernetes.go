package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"  // auth for GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // auth for OIDC
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client abstracts the API to allow testing.
type Client interface {
	InClusterConfig() (*rest.Config, error)
}

// K8sClient wraps kubernetes client-go so it can be mocked.
type K8sClient struct{}

// InClusterConfig wraps the client-go method with the same name.
func (kc K8sClient) InClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}

// ErrNoReadableKubeConfig represents any error that prevents the client from opening a kubeconfig file.
var ErrNoReadableKubeConfig = errors.New("unable to open kubeconfig file")

func kubeClient() (*kubernetes.Clientset, error) {
	return kubeClientType(K8sClient{})
}

func kubeClientType(kc Client) (*kubernetes.Clientset, error) {
	config, err := kubeClientConfig(kc)
	if err != nil {
		return nil, err
	}
	kube, err := kubernetes.NewForConfig(config)
	return kube, err
}

func kubeClientConfig(kc Client) (*rest.Config, error) {
	if rootConfig.kubeConfig != "" {
		return kubeClientConfigLocal()
	}

	if config, err := kc.InClusterConfig(); err == nil {
		log.Info("Running inside cluster, using the cluster config")
		return config, nil
	}

	log.Info("Not running inside cluster, using local config")
	home, ok := os.LookupEnv("HOME")
	if !ok || home == "" {
		log.Error("Unable to load kubeconfig. No config file specified and $HOME not found.")
		return nil, ErrNoReadableKubeConfig
	}

	rootConfig.kubeConfig = filepath.Join(home, ".kube", "config")
	return kubeClientConfigLocal()
}

func kubeClientConfigLocal() (*rest.Config, error) {
	if _, err := os.Stat(rootConfig.kubeConfig); err != nil {
		log.Errorf("Unable to load kubeconfig. Could not open file %s.", rootConfig.kubeConfig)
		return nil, ErrNoReadableKubeConfig
	}
	return clientcmd.BuildConfigFromFlags("", rootConfig.kubeConfig)
}

func getDeployments(clientset *kubernetes.Clientset) *DeploymentListV1 {
	deploymentClient := clientset.AppsV1().Deployments(rootConfig.namespace)
	deployments, err := deploymentClient.List(ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return deployments
}

func getStatefulSets(clientset *kubernetes.Clientset) *StatefulSetListV1 {
	statefulSetClient := clientset.AppsV1().StatefulSets(rootConfig.namespace)
	statefulSets, err := statefulSetClient.List(ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return statefulSets
}

func getDaemonSets(clientset *kubernetes.Clientset) *DaemonSetListV1 {
	daemonSetClient := clientset.AppsV1().DaemonSets(rootConfig.namespace)
	daemonSets, err := daemonSetClient.List(ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return daemonSets
}

func getPods(clientset *kubernetes.Clientset) *PodListV1 {
	podClient := clientset.CoreV1().Pods(rootConfig.namespace)
	pods, err := podClient.List(ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return pods
}

func getReplicationControllers(clientset *kubernetes.Clientset) *ReplicationControllerListV1 {
	replicationControllerClient := clientset.CoreV1().ReplicationControllers(rootConfig.namespace)
	replicationControllers, err := replicationControllerClient.List(ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return replicationControllers
}

func getNetworkPolicies(clientset *kubernetes.Clientset) *NetworkPolicyListV1 {
	netPolClient := clientset.NetworkingV1().NetworkPolicies(rootConfig.namespace)
	netPols, err := netPolClient.List(ListOptionsV1{})
	if err != nil {
		log.Error(err)
	}
	return netPols
}

func getNamespaces(clientset *kubernetes.Clientset) *NamespaceList {
	namespaceClient := clientset.CoreV1().Namespaces()
	listOptions := ListOptions{}

	if rootConfig.namespace != "" {
		// Select only the specified namespace
		listOptions = ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name=%s", rootConfig.namespace),
		}
	}

	namespaces, err := namespaceClient.List(listOptions)
	if err != nil {
		log.Error(err)
	}
	return namespaces
}

func getKubernetesVersion(clientset kubernetes.Interface) (*version.Info, error) {
	discoveryClient := clientset.Discovery()
	return discoveryClient.ServerVersion()
}

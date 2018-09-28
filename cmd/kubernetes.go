package cmd

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	networking "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/version"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"  // auth for GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // auth for OIDC
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ErrNoReadableKubeConfig represents any error that prevents the client from opening a kubeconfig file.
var ErrNoReadableKubeConfig = errors.New("unable to open kubeconfig file")

func kubeClientConfig() (*rest.Config, error) {
	if rootConfig.kubeConfig != "" {
		return kubeClientConfigLocal()
	}

	if config, err := rest.InClusterConfig(); err == nil {
		log.Info("Running inside cluster, using the cluster config")
		return config, nil
	}

	log.Info("Not running inside cluster, using local config")
	home, ok := os.LookupEnv("HOME")
	if !ok {
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

func kubeClient() (*kubernetes.Clientset, error) {
	config, err := kubeClientConfig()
	if err != nil {
		return nil, err
	}
	kube, err := kubernetes.NewForConfig(config)
	return kube, err
}

func getDeployments(clientset *kubernetes.Clientset) *DeploymentList {
	deploymentClient := clientset.AppsV1beta1().Deployments(rootConfig.namespace)
	deployments, err := deploymentClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return deployments
}

func getStatefulSets(clientset *kubernetes.Clientset) *StatefulSetList {
	statefulSetClient := clientset.AppsV1beta1().StatefulSets(rootConfig.namespace)
	statefulSets, err := statefulSetClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return statefulSets
}

func getDaemonSets(clientset *kubernetes.Clientset) *DaemonSetList {
	daemonSetClient := clientset.ExtensionsV1beta1().DaemonSets(rootConfig.namespace)
	daemonSets, err := daemonSetClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return daemonSets
}

func getPods(clientset *kubernetes.Clientset) *PodList {
	podClient := clientset.CoreV1().Pods(rootConfig.namespace)
	pods, err := podClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return pods
}

func getReplicationControllers(clientset *kubernetes.Clientset) *ReplicationControllerList {
	replicationControllerClient := clientset.CoreV1().ReplicationControllers(rootConfig.namespace)
	replicationControllers, err := replicationControllerClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return replicationControllers
}

func getNetworkPolicies(clientset *kubernetes.Clientset) *networking.NetworkPolicyList {
	netPolClient := clientset.NetworkingV1().NetworkPolicies(rootConfig.namespace)
	netPols, err := netPolClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return netPols
}

func printKubernetesVersion(clientset *kubernetes.Clientset) {
	discoveryClient := clientset.Discovery()
	serverInfo, err := discoveryClient.ServerVersion()
	if err != nil {
		log.Error(err)
	}
	log.WithFields(log.Fields{
		"Major":    serverInfo.Major,
		"Minor":    serverInfo.Minor,
		"Platform": serverInfo.Platform,
	}).Info("Kubernetes server version")

	clientInfo := version.Get()
	log.WithFields(log.Fields{
		"Major":    clientInfo.Major,
		"Minor":    clientInfo.Minor,
		"Platform": clientInfo.Platform,
	}).Info("Kubernetes client version")
}

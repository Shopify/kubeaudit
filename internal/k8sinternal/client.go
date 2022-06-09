package k8sinternal

import (
	"context"
	"errors"
	"os"

	"github.com/Shopify/kubeaudit/pkg/k8s"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	// add authentication support to the kubernetes code
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
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

// NewKubeClientLocal creates a new kube client for local mode
func NewKubeClientLocal(configPath string) (KubeClient, error) {
	var kubeconfig *rest.Config
	var err error

	if configPath == "" {
		kubeconfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}},
		).ClientConfig()
	} else {
		if _, err = os.Stat(configPath); err != nil {
			return nil, ErrNoReadableKubeConfig
		}
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", configPath)
	}

	if err != nil {
		return nil, err
	}

	return newKubeClientFromConfig(kubeconfig)
}

// NewKubeClientCluster creates a new kube client for cluster mode
func NewKubeClientCluster(client Client) (KubeClient, error) {
	config, err := client.InClusterConfig()
	if err != nil {
		return nil, err
	}
	log.Info("Running inside cluster, using the cluster config")
	return newKubeClientFromConfig(config)
}

// newKubeClientFromConfig creates a new dynamic client with discovery or returns an error.
func newKubeClientFromConfig(config *rest.Config) (KubeClient, error) {
	discovery, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	dynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return NewKubeClient(dynamic, discovery), nil
}

// IsRunningInCluster returns true if kubeaudit is running inside a cluster
func IsRunningInCluster(client Client) bool {
	_, err := client.InClusterConfig()
	return err == nil
}

type ClientOptions struct {
	// Namespace filters resources by namespace. Defaults to all namespaces.
	Namespace string
	// IncludeGenerated is a boolean option to include generated resources.
	IncludeGenerated bool
}

type KubeClient interface {
	// GetAllResources gets all supported resources from the cluster
	GetAllResources(options ClientOptions) ([]k8s.Resource, error)
	// GetKubernetesVersion returns the kubernetes client version
	GetKubernetesVersion() (*version.Info, error)
	// ServerPreferredResources returns the supported resources with the version preferred by the server.
	ServerPreferredResources() ([]*metav1.APIResourceList, error)
}

type kubeClient struct {
	dynamicClient   dynamic.Interface
	discoveryClient discovery.DiscoveryInterface
}

func NewKubeClient(dynamic dynamic.Interface, discovery discovery.DiscoveryInterface) KubeClient {
	return &kubeClient{dynamicClient: dynamic, discoveryClient: discovery}
}

// GetAllResources gets all supported resources from the cluster
func (kc kubeClient) GetAllResources(options ClientOptions) ([]k8s.Resource, error) {
	var resources []k8s.Resource

	lists, err := kc.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	if lists != nil {
		for _, list := range lists {
			if len(list.APIResources) == 0 {
				continue
			}
			gv, err := schema.ParseGroupVersion(list.GroupVersion)
			if err != nil {
				continue
			}
			for _, apiresource := range list.APIResources {
				if len(apiresource.Verbs) == 0 {
					continue
				}
				gvr := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: apiresource.Name}

				// Namespace has to be included as a resource to audit if it is specified.
				if apiresource.Name == "namespaces" && options.Namespace != "" {
					unstructured, err := kc.dynamicClient.Resource(gvr).Get(context.Background(), options.Namespace, metav1.GetOptions{})
					if err == nil {
						r, err := unstructuredToObject(unstructured)
						if err == nil {
							resources = append(resources, r)
						}
					}
				} else {
					unstructuredList, err := kc.dynamicClient.Resource(gvr).Namespace(options.Namespace).List(context.Background(), metav1.ListOptions{})
					if err == nil {
						for _, unstructured := range unstructuredList.Items {
							r, err := unstructuredToObject(&unstructured)
							if err == nil {
								resources = append(resources, r)
							}
						}
					}
				}
			}
		}
	}

	if !options.IncludeGenerated {
		resources = excludeGenerated(resources)
	}
	return resources, nil
}

// unstructuredToObject unstructured to Go typed object conversions
func unstructuredToObject(unstructured *unstructured.Unstructured) (k8s.Resource, error) {
	obj, err := scheme.New(unstructured.GroupVersionKind())
	if err == nil {
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.UnstructuredContent(), obj)
	}
	return obj, err
}

// excludeGenerated filters out generated resources (eg. pods generated by deployments)
func excludeGenerated(resources []k8s.Resource) []k8s.Resource {
	var filteredResources []k8s.Resource
	for _, resource := range resources {
		if resource != nil {
			obj, _ := resource.(metav1.ObjectMetaAccessor)
			if obj != nil {
				meta := obj.GetObjectMeta()
				if meta != nil {
					if len(meta.GetOwnerReferences()) == 0 {
						filteredResources = append(filteredResources, resource)
					}
				}
			}
		}
	}
	return filteredResources
}

// GetKubernetesVersion returns the kubernetes client version
func (kc kubeClient) GetKubernetesVersion() (*version.Info, error) {
	return kc.discoveryClient.ServerVersion()
}

// ServerPreferredResources returns the supported resources with the version preferred by the server.
func (kc kubeClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	list, err := discovery.ServerPreferredResources(kc.discoveryClient)
	// If a group is not served by the cluster the resources of this group will not be audited.
	var e *discovery.ErrGroupDiscoveryFailed
	if errors.As(err, &e) {
		return list, nil
	}
	return list, err
}

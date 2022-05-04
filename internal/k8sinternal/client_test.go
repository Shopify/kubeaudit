package k8sinternal_test

import (
	"errors"
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure" // auth for AKS clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"   // auth for GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"  // auth for OIDC
	"k8s.io/client-go/rest"
)

type MockK8sClient struct {
	mock.Mock
}

func (kc *MockK8sClient) InClusterConfig() (*rest.Config, error) {
	args := kc.Called()
	return args.Get(0).(*rest.Config), args.Error(1)
}

func TestKubeClientConfigLocal(t *testing.T) {
	assert := assert.New(t)

	_, err := k8sinternal.NewKubeClientLocal("/notarealfile")
	assert.Equal(k8sinternal.ErrNoReadableKubeConfig, err)

	_, err = k8sinternal.NewKubeClientLocal("client.go")
	assert.NotEqual(k8sinternal.ErrNoReadableKubeConfig, err)
	assert.NotNil(err)
}

func TestKubeClientConfigCluster(t *testing.T) {
	assert := assert.New(t)

	client := &MockK8sClient{}
	var config *rest.Config = nil
	client.On("InClusterConfig").Return(config, errors.New("mock error"))
	clientset, err := k8sinternal.NewKubeClientCluster(client)
	assert.Nil(clientset)
	assert.NotNil(err)

	client = &MockK8sClient{}
	client.On("InClusterConfig").Return(&rest.Config{}, nil)
	clientset, err = k8sinternal.NewKubeClientCluster(client)
	assert.NotNil(clientset)
	assert.NoError(err)
}

func TestIsRunningInCluster(t *testing.T) {
	assert := assert.New(t)

	client := &MockK8sClient{}
	var config *rest.Config = nil
	client.On("InClusterConfig").Return(config, errors.New("mock error"))
	assert.False(k8sinternal.IsRunningInCluster(client))

	client = &MockK8sClient{}
	client.On("InClusterConfig").Return(&rest.Config{}, nil)
	assert.True(k8sinternal.IsRunningInCluster(client))
}

func TestGetAllResources(t *testing.T) {
	resourceTemplates := []k8s.Resource{
		k8s.NewDeployment(),
		k8s.NewPod(),
		k8s.NewNamespace(),
		k8s.NewDaemonSet(),
		k8s.NewNetworkPolicy(),
		k8s.NewReplicationController(),
		k8s.NewStatefulSet(),
		k8s.NewPodTemplate(),
		k8s.NewCronJob(),
		k8s.NewServiceAccount(),
		k8s.NewService(),
		k8s.NewJob(),
	}
	namespaces := []string{"foo", "bar"}

	resources := make([]runtime.Object, 0, len(resourceTemplates)*len(namespaces))
	for _, template := range resourceTemplates {
		for _, namespace := range namespaces {
			resource := template.DeepCopyObject()
			setNamespace(resource, namespace)
			resources = append(resources, resource)
		}
	}

	clientset := fakeclientset.NewSimpleClientset(resources...)
	assert.Len(t, k8sinternal.GetAllResources(clientset, k8sinternal.ClientOptions{}), len(resourceTemplates)*len(namespaces))

	// Because field selectors are handled server-side, the fake clientset does not support them
	// which means the Namespace resources don't get filtered (this is not a problem when using
	// a real clientset)
	// See https://github.com/kubernetes/client-go/issues/326
	assert.Len(
		t,
		k8sinternal.GetAllResources(clientset, k8sinternal.ClientOptions{Namespace: namespaces[0]}),
		len(resourceTemplates)+(len(namespaces)-1),
	)
}

func setNamespace(resource k8s.Resource, namespace string) {
	if _, ok := resource.(*k8s.NamespaceV1); ok {
		k8s.GetObjectMeta(resource).Name = namespace
	} else {
		k8s.GetObjectMeta(resource).Namespace = namespace
	}
}

func TestGetKubernetesVersion(t *testing.T) {
	client := fakeclientset.NewSimpleClientset()
	fakeDiscovery, ok := client.Discovery().(*fakediscovery.FakeDiscovery)
	if !ok {
		t.Fatalf("couldn't mock server version")
	}

	fakeDiscovery.FakedServerVersion = &version.Info{
		Major:     "0",
		Minor:     "0",
		GitCommit: "0000",
		Platform:  "ACME 8-bit",
	}

	r, err := k8sinternal.GetKubernetesVersion(client)
	assert.Nil(t, err)
	assert.EqualValues(t, *fakeDiscovery.FakedServerVersion, *r)
}

func TestIncludeGenerated(t *testing.T) {
	// The "IncludeGenerated" option only applies to local and cluster mode
	if !test.UseKind() {
		return
	}

	namespace := "include-generated"
	defer test.DeleteNamespace(t, namespace)
	test.CreateNamespace(t, namespace)
	test.ApplyManifest(t, "./fixtures/include-generated.yml", namespace)

	clientset, err := k8sinternal.NewKubeClientLocal("")
	require.NoError(t, err)

	// Test IncludeGenerated = false
	resources := k8sinternal.GetAllResources(
		clientset,
		k8sinternal.ClientOptions{Namespace: namespace, IncludeGenerated: false},
	)
	assert.False(t, hasPod(resources), "Expected no pods for IncludeGenerated=false")

	// Test IncludeGenerated unspecified defaults to false
	resources = k8sinternal.GetAllResources(
		clientset,
		k8sinternal.ClientOptions{Namespace: namespace},
	)
	assert.False(t, hasPod(resources), "Expected no pods if IncludeGenerated is unspecified (ie. default to false)")

	// Test IncludeGenerated = true
	resources = k8sinternal.GetAllResources(
		clientset,
		k8sinternal.ClientOptions{Namespace: namespace, IncludeGenerated: true},
	)
	assert.True(t, hasPod(resources), "Expected pods for IncludeGenerated=true")
}

func hasPod(resources []k8s.Resource) bool {
	for _, resource := range resources {
		if k8s.IsPodV1(resource) {
			return true
		}
	}
	return false
}
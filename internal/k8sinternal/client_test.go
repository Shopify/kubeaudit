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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	fakedynamic "k8s.io/client-go/dynamic/fake"
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
	kubeclient, err := k8sinternal.NewKubeClientCluster(client)
	assert.Nil(kubeclient)
	assert.NotNil(err)

	client = &MockK8sClient{}
	client.On("InClusterConfig").Return(&rest.Config{}, nil)
	kubeclient, err = k8sinternal.NewKubeClientCluster(client)
	assert.NotNil(kubeclient)
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

	client := newFakeKubeClient(resources...)
	k8sresources, err := client.GetAllResources(k8sinternal.ClientOptions{})
	require.NoError(t, err)
	assert.Len(t, k8sresources, len(resourceTemplates)*len(namespaces))

	k8sresources, err = client.GetAllResources(k8sinternal.ClientOptions{Namespace: namespaces[0]})
	require.NoError(t, err)
	assert.Len(t, k8sresources, len(resourceTemplates))
}

func setNamespace(resource k8s.Resource, namespace string) {
	if _, ok := resource.(*k8s.NamespaceV1); ok {
		k8s.GetObjectMeta(resource).SetName(namespace)
	} else {
		k8s.GetObjectMeta(resource).SetNamespace(namespace)
	}
}

func TestGetKubernetesVersion(t *testing.T) {
	serverVersion := &version.Info{
		Major:     "0",
		Minor:     "0",
		GitCommit: "0000",
		Platform:  "ACME 8-bit",
	}

	client := newFakeKubeClientWithServerVersion(serverVersion)
	r, err := client.GetKubernetesVersion()
	assert.Nil(t, err)
	assert.EqualValues(t, *serverVersion, *r)
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

	client, err := k8sinternal.NewKubeClientLocal("")
	require.NoError(t, err)

	// Test IncludeGenerated = false
	resources, err := client.GetAllResources(
		k8sinternal.ClientOptions{Namespace: namespace, IncludeGenerated: false},
	)
	require.NoError(t, err)
	assert.False(t, hasPod(resources), "Expected no pods for IncludeGenerated=false")

	// Test IncludeGenerated unspecified defaults to false
	resources, err = client.GetAllResources(
		k8sinternal.ClientOptions{Namespace: namespace},
	)
	require.NoError(t, err)
	assert.False(t, hasPod(resources), "Expected no pods if IncludeGenerated is unspecified (ie. default to false)")

	// Test IncludeGenerated = true
	resources, err = client.GetAllResources(
		k8sinternal.ClientOptions{Namespace: namespace, IncludeGenerated: true},
	)
	require.NoError(t, err)
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

func newFakeKubeClient(resources ...runtime.Object) k8sinternal.KubeClient {
	return newFakeKubeClientWithServerVersion(nil, resources...)
}

func newFakeKubeClientWithServerVersion(serverversion *version.Info, resources ...runtime.Object) k8sinternal.KubeClient {
	clientset := fakeclientset.NewSimpleClientset()
	fakeDiscovery, _ := clientset.Discovery().(*fakediscovery.FakeDiscovery)
	if serverversion != nil {
		fakeDiscovery.FakedServerVersion = serverversion
	}
	unstructuredresources := []runtime.Object{}
	gvrToListKind := map[schema.GroupVersionResource]string{}
	gvAPIResources := map[string][]metav1.APIResource{}
	for _, r := range resources {
		gvk := r.GetObjectKind().GroupVersionKind()
		listGVK := gvk
		listGVK.Kind += "List"

		u := unstructured.Unstructured{}
		u.SetGroupVersionKind(r.GetObjectKind().GroupVersionKind())
		u.SetName(k8s.GetObjectMeta(r).GetName())
		u.SetNamespace(k8s.GetObjectMeta(r).GetNamespace())
		unstructuredresources = (append(unstructuredresources, &u))

		kind := r.GetObjectKind().GroupVersionKind().Kind
		plural, _ := meta.UnsafeGuessKindToResource(r.GetObjectKind().GroupVersionKind())
		apiresource := metav1.APIResource{Name: plural.Resource, Namespaced: false, Group: gvk.Group, Version: gvk.Version, Kind: kind, Verbs: metav1.Verbs{"list"}}
		gvr := schema.GroupVersionResource{Group: apiresource.Group, Version: apiresource.Version, Resource: apiresource.Name}
		if _, ok := gvrToListKind[gvr]; !ok {
			gvrToListKind[gvr] = kind + "List"
			gv := gvk.GroupVersion().String()
			gvAPIResources[gv] = append(gvAPIResources[gv], apiresource)
		}
	}
	for gv, apiresources := range gvAPIResources {
		fakeDiscovery.Resources = append(fakeDiscovery.Resources, &metav1.APIResourceList{
			GroupVersion: gv,
			APIResources: apiresources})
	}
	fakedynamic := fakedynamic.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), gvrToListKind, unstructuredresources...)
	return k8sinternal.NewKubeClient(fakedynamic, fakeDiscovery)
}

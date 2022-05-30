package k8sinternal_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fixtureDir = "../test/fixtures"

func TestNewTrue(t *testing.T) {
	assert.True(t, *k8s.NewTrue())
}

func TestNewFalse(t *testing.T) {
	assert.False(t, *k8s.NewFalse())
}

func TestEncodeDecode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	deployment := k8s.NewDeployment()
	deployment.ObjectMeta = k8s.ObjectMetaV1{Namespace: "foo"}
	deployment.Spec.Template.Spec.Containers = []k8s.ContainerV1{{Name: "bar"}}

	expectedManifest, err := ioutil.ReadFile("fixtures/test-encode-decode.yml")
	require.NoError(err)

	encoded, err := k8sinternal.EncodeResource(deployment)
	require.NoError(err)
	assert.Equal(string(expectedManifest), string(encoded))

	decoded, err := k8sinternal.DecodeResource(expectedManifest)
	require.NoError(err)
	assert.Equal(deployment, decoded)
}

func TestGetContainers(t *testing.T) {
	for _, resource := range getAllResources(t) {
		containers := k8s.GetContainers(resource)
		switch resource.(type) {
		case *k8s.NamespaceV1:
			assert.Nil(t, containers)
		default:
			assert.NotEmpty(t, containers)
		}
	}
}

func TestGetAnnotations(t *testing.T) {
	annotations := map[string]string{"foo": "bar"}
	deployment := k8s.NewDeployment()
	deployment.Spec.Template.ObjectMeta.SetAnnotations(annotations)
	assert.Equal(t, annotations, k8s.GetAnnotations(deployment))
}

func TestGetLabels(t *testing.T) {
	labels := map[string]string{"foo": "bar"}
	deployment := k8s.NewDeployment()
	deployment.Spec.Template.ObjectMeta.SetLabels(labels)
	assert.Equal(t, labels, k8s.GetLabels(deployment))
}

func TestGetObjectMeta(t *testing.T) {
	assert := assert.New(t)
	objectMeta := k8s.ObjectMetaV1{Name: "foo"}
	podObjectMeta := k8s.ObjectMetaV1{Name: "bar"}

	deployment := k8s.NewDeployment()
	deployment.ObjectMeta = objectMeta
	deployment.Spec.Template.ObjectMeta = podObjectMeta
	assert.Equal(&objectMeta, k8s.GetObjectMeta(deployment))
	assert.Equal(&podObjectMeta, k8s.GetPodObjectMeta(deployment))

	pod := k8s.NewPod()
	pod.ObjectMeta = objectMeta
	assert.Equal(&objectMeta, k8s.GetObjectMeta(pod))
	assert.Equal(&objectMeta, k8s.GetPodObjectMeta(pod))

	namespace := k8s.NewNamespace()
	namespace.ObjectMeta = objectMeta
	assert.Equal(&objectMeta, k8s.GetObjectMeta(namespace))
	assert.Equal(&objectMeta, k8s.GetPodObjectMeta(namespace))
}

func TestGetPodTemplateSpec(t *testing.T) {
	for _, resource := range getAllResources(t) {
		podTemplateSpec := k8s.GetPodTemplateSpec(resource)
		switch resource.(type) {
		case *k8s.PodV1, *k8s.NamespaceV1:
			assert.Nil(t, podTemplateSpec)
		default:
			assert.NotNil(t, podTemplateSpec)
		}
	}
}

func TestUnsupportedResource(t *testing.T) {
	unsupported := &k8s.UnsupportedType{}
	assert.Nil(t, k8s.GetAnnotations(unsupported))
	assert.Nil(t, k8s.GetLabels(unsupported))
	assert.Nil(t, k8s.GetContainers(unsupported))
}

func getAllResources(t *testing.T) (resources []k8s.Resource) {
	fixtureDir := "../test/fixtures/all_resources"
	for _, fixture := range test.GetAllFileNames(t, fixtureDir) {
		resources = append(resources, getResourcesFromManifest(t, path.Join(fixtureDir, fixture))...)
	}
	return
}

func getResourcesFromManifest(t *testing.T, manifest string) (resources []k8s.Resource) {
	assert := assert.New(t)

	data, err := ioutil.ReadFile(manifest)
	require.NoError(t, err)

	bufSlice := bytes.Split(data, []byte("---"))

	for _, b := range bufSlice {
		resource, err := k8sinternal.DecodeResource(b)
		if err != nil {
			continue
		}
		assert.NotNil(resource)
		resources = append(resources, resource)
	}
	return
}

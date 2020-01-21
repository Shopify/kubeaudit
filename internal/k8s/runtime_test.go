package k8s

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/Shopify/kubeaudit/k8stypes"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "../test/fixtures"

func TestNewTrue(t *testing.T) {
	assert.True(t, *NewTrue())
}

func TestNewFalse(t *testing.T) {
	assert.False(t, *NewFalse())
}

func TestEncodeDecode(t *testing.T) {
	assert := assert.New(t)

	deployment := k8stypes.NewDeployment()
	deployment.ObjectMeta = k8stypes.ObjectMetaV1{Namespace: "foo"}
	deployment.Spec.Template.Spec.Containers = []k8stypes.ContainerV1{{Name: "bar"}}

	expectedManifest, err := ioutil.ReadFile("fixtures/test_encode_decode.yml")
	assert.NoError(err)

	encoded, err := EncodeResource(deployment)
	assert.Nil(err)
	assert.Equal(string(expectedManifest), string(encoded))

	decoded, err := DecodeResource(expectedManifest)
	assert.Nil(err)
	assert.Equal(deployment, decoded)
}

func TestGetContainers(t *testing.T) {
	for _, resource := range getAllResources(t) {
		if !k8stypes.IsSupportedResourceType(resource) {
			continue
		}
		containers := GetContainers(resource)
		switch resource.(type) {
		case *k8stypes.NamespaceV1:
			assert.Nil(t, containers)
		default:
			assert.NotEmpty(t, containers)
		}
	}
}

func TestGetAnnotations(t *testing.T) {
	annotations := map[string]string{"foo": "bar"}
	deployment := k8stypes.NewDeployment()
	deployment.Spec.Template.ObjectMeta.SetAnnotations(annotations)
	assert.Equal(t, annotations, GetAnnotations(deployment))
}

func TestGetLabels(t *testing.T) {
	labels := map[string]string{"foo": "bar"}
	deployment := k8stypes.NewDeployment()
	deployment.Spec.Template.ObjectMeta.SetLabels(labels)
	assert.Equal(t, labels, GetLabels(deployment))
}

func TestGetObjectMeta(t *testing.T) {
	assert := assert.New(t)
	objectMeta := k8stypes.ObjectMetaV1{Name: "foo"}
	podObjectMeta := k8stypes.ObjectMetaV1{Name: "bar"}

	deployment := k8stypes.NewDeployment()
	deployment.ObjectMeta = objectMeta
	deployment.Spec.Template.ObjectMeta = podObjectMeta
	assert.Equal(objectMeta, *GetObjectMeta(deployment))
	assert.Equal(podObjectMeta, *GetPodObjectMeta(deployment))

	pod := k8stypes.NewPod()
	pod.ObjectMeta = objectMeta
	assert.Equal(objectMeta, *GetObjectMeta(pod))
	assert.Equal(objectMeta, *GetPodObjectMeta(pod))

	namespace := k8stypes.NewNamespace()
	namespace.ObjectMeta = objectMeta
	assert.Equal(objectMeta, *GetObjectMeta(namespace))
	assert.Equal(objectMeta, *GetPodObjectMeta(namespace))
}

func TestGetPodTemplateSpec(t *testing.T) {
	for _, resource := range getAllResources(t) {
		if !k8stypes.IsSupportedResourceType(resource) {
			continue
		}
		podTemplateSpec := GetPodTemplateSpec(resource)
		switch resource.(type) {
		case *k8stypes.PodV1, *k8stypes.NamespaceV1:
			assert.Nil(t, podTemplateSpec)
		default:
			assert.NotNil(t, podTemplateSpec)
		}
	}
}

func TestUnsupportedResource(t *testing.T) {
	unsupported := &k8stypes.UnsupportedType{}
	assert.Nil(t, GetAnnotations(unsupported))
	assert.Nil(t, GetLabels(unsupported))
	assert.Nil(t, GetContainers(unsupported))
}

func getAllResources(t *testing.T) (resources []k8stypes.Resource) {
	fixtures := []string{"all-resources_v1.yml", "all-resources_v1beta1.yml", "all-resources_v1beta2.yml"}
	for _, fixture := range fixtures {
		resources = append(resources, getResourcesFromManifest(t, path.Join(fixtureDir, fixture))...)
	}
	return
}

func getResourcesFromManifest(t *testing.T, manifest string) (resources []k8stypes.Resource) {
	assert := assert.New(t)

	data, err := ioutil.ReadFile(manifest)
	assert.Nil(err)

	bufSlice := bytes.Split(data, []byte("---"))

	for _, b := range bufSlice {
		resource, err := DecodeResource(b)
		if err != nil {
			continue
		}
		assert.NotNil(resource)
		resources = append(resources, resource)
	}
	return
}

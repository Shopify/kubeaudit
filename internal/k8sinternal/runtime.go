package k8sinternal

import (
	"github.com/Shopify/kubeaudit/pkg/k8s"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func DecodeResource(b []byte) (k8s.Resource, error) {
	decoder := codecs.UniversalDeserializer()
	return k8sRuntime.Decode(decoder, b)
}

func EncodeResource(resource k8s.Resource) ([]byte, error) {
	info, _ := k8sRuntime.SerializerInfoForMediaType(codecs.SupportedMediaTypes(), "application/yaml")
	groupVersion := schema.GroupVersion{Group: resource.GetObjectKind().GroupVersionKind().Group, Version: resource.GetObjectKind().GroupVersionKind().Version}
	encoder := codecs.EncoderForVersion(info.Serializer, groupVersion)
	return k8sRuntime.Encode(encoder, resource)
}

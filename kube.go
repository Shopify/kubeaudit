package kubeaudit

import "github.com/Shopify/kubeaudit/pkg/k8s"

// ErrorUnsupportedResource occurs when Kubeaudit doesn't know how to audit the resource
const ErrorUnsupportedResource = "Unsupported resource"

// RedundantAuditorOverride is the audit result name given when an override label is used to disable an auditor,
// but that auditor found no security issues so the label is redundant
const RedundantAuditorOverride = "RedundantAuditorOverride"

// KubeResource is a wrapper around a Kubernetes object
type KubeResource interface {
	// Object is a pointer to a Kubernetes resource. The resource may be modified by multiple auditors
	Object() k8s.Resource
	// Bytes is the original byte representation of the resource
	Bytes() []byte
}

// Implements KubeResource
type kubeResource struct {
	object k8s.Resource
	bytes  []byte
}

func (k *kubeResource) Object() k8s.Resource {
	return k.object
}

func (k *kubeResource) Bytes() []byte {
	return k.bytes
}

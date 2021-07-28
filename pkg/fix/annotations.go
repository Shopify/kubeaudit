package fix

import (
	"fmt"

	"github.com/Shopify/kubeaudit/pkg/k8s"
)

// FixBySettingPodAnnotation implements PendingFix
type BySettingPodAnnotation struct {
	Key   string
	Value string
}

// Apply sets the pod annotation to the specified value
func (pending *BySettingPodAnnotation) Apply(resource k8s.Resource) []k8s.Resource {
	objectMeta := k8s.GetPodObjectMeta(resource)

	if objectMeta.GetAnnotations() == nil {
		objectMeta.SetAnnotations(map[string]string{})
	}

	objectMeta.GetAnnotations()[pending.Key] = pending.Value

	return nil
}

// Plan is a description of what apply will do
func (pending *BySettingPodAnnotation) Plan() string {
	return fmt.Sprintf("Set pod-level annotation '%v' to '%v'", pending.Key, pending.Value)
}

// FixByAddingPodAnnotation implements PendingFix
type ByAddingPodAnnotation struct {
	Key   string
	Value string
}

// Apply adds the pod annotation
func (pending *ByAddingPodAnnotation) Apply(resource k8s.Resource) []k8s.Resource {
	objectMeta := k8s.GetPodObjectMeta(resource)

	if objectMeta.GetAnnotations() == nil {
		objectMeta.SetAnnotations(map[string]string{})
	}

	objectMeta.GetAnnotations()[pending.Key] = pending.Value

	return nil
}

// Plan is a description of what apply will do
func (pending *ByAddingPodAnnotation) Plan() string {
	return fmt.Sprintf("Add pod-level annotation '%v: %v'", pending.Key, pending.Value)
}

type ByRemovingPodAnnotation struct {
	Key string
}

// Apply removes the pod annotation
func (pending *ByRemovingPodAnnotation) Apply(resource k8s.Resource) []k8s.Resource {
	objectMeta := k8s.GetPodObjectMeta(resource)

	if objectMeta.GetAnnotations() == nil {
		return nil
	}

	delete(objectMeta.GetAnnotations(), pending.Key)

	return nil
}

// Plan is a description of what apply will do
func (pending *ByRemovingPodAnnotation) Plan() string {
	return fmt.Sprintf("Remove pod-level annotation '%v'", pending.Key)
}

package seccomp

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/fix"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	apiv1 "k8s.io/api/core/v1"
)

const Name = "seccomp"

const (
	// SeccompAnnotationMissing occurs when there are no seccomp annotations (pod nor container level)
	SeccompAnnotationMissing = "SeccompAnnotationMissing"
	// SeccompDeprecatedPod occurs when the pod-level seccomp annotation is set to a deprecated value
	SeccompDeprecatedPod = "SeccompDeprecatedPod"
	// SeccompDisabledPod occurs when the pod-level seccomp annotation is set to a value which disables seccomp
	SeccompDisabledPod = "SeccompDisabledPod"
	// SeccompDeprecatedContainer occurs when the container-level seccomp annotation is set to a deprecated value
	SeccompDeprecatedContainer = "SeccompDeprecatedContainer"
	// SeccompDisabledContainer occurs when the container-level seccomp annotation is set to a value which disables seccomp
	SeccompDisabledContainer = "SeccompDisabledContainer"
)

const (
	// ContainerAnnotationKeyPrefix represents the key of a seccomp profile applied to one container of a pod
	ContainerAnnotationKeyPrefix = apiv1.SeccompContainerAnnotationKeyPrefix
	// PodAnnotationKey represents the key of a seccomp profile applied to all containers of a pod
	PodAnnotationKey = apiv1.SeccompPodAnnotationKey
	// ProfileRuntimeDefault represents the default seccomp profile used by container runtime
	ProfileRuntimeDefault = apiv1.SeccompProfileRuntimeDefault
	// ProfileNamePrefix is the prefix for a custom seccomp profile
	ProfileNamePrefix = "localhost/"
	// DeprecatedProfileRuntimeDefault represents the default seccomp profile used by docker.
	// This is now deprecated and should be replaced by SeccompProfileRuntimeDefault
	DeprecatedProfileRuntimeDefault = apiv1.DeprecatedSeccompProfileDockerDefault
)

// Seccomp implements Auditable
type Seccomp struct{}

func New() *Seccomp {
	return &Seccomp{}
}

// Audit checks that Seccomp is enabled for all containers
func (a *Seccomp) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	auditResult := auditPod(resource)
	if auditResult != nil {
		auditResults = append(auditResults, auditResult)
	}

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container, resource)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditPod(resource k8s.Resource) *kubeaudit.AuditResult {
	annotations := k8s.GetAnnotations(resource)
	PodAnnotationKey := apiv1.SeccompPodAnnotationKey

	if isSeccompAnnotationMissing(annotations, PodAnnotationKey) {
		// If all the containers have container-level seccomp annotations then we don't need a pod-level annotation
		if isSeccompEnabledForContainers(annotations, resource) {
			return nil
		}

		return &kubeaudit.AuditResult{
			Name:     SeccompAnnotationMissing,
			Severity: kubeaudit.Error,
			Message:  fmt.Sprintf("Seccomp annotation is missing. The annotation %s: %s should be added.", PodAnnotationKey, ProfileRuntimeDefault),
			PendingFix: &fix.ByAddingPodAnnotation{
				Key:   PodAnnotationKey,
				Value: ProfileRuntimeDefault,
			},
			Metadata: kubeaudit.Metadata{
				"MissingAnnotation": PodAnnotationKey,
			},
		}
	}

	podSeccompProfile := annotations[PodAnnotationKey]

	if isSeccompProfileDeprecated(podSeccompProfile) {
		return &kubeaudit.AuditResult{
			Name:     SeccompDeprecatedPod,
			Severity: kubeaudit.Error,
			Message:  fmt.Sprintf("Seccomp pod annotation is set to deprecated value %s. It should be set to %s instead.", podSeccompProfile, ProfileRuntimeDefault),
			PendingFix: &fix.BySettingPodAnnotation{
				Key:   PodAnnotationKey,
				Value: ProfileRuntimeDefault,
			},
			Metadata: kubeaudit.Metadata{
				"AnnotationKey":   PodAnnotationKey,
				"AnnotationValue": podSeccompProfile,
			},
		}
	}

	if !isSeccompEnabled(podSeccompProfile) {
		return &kubeaudit.AuditResult{
			Name:     SeccompDisabledPod,
			Severity: kubeaudit.Error,
			Message:  fmt.Sprintf("Seccomp pod annotation is set to %s which disables Seccomp. It should be set to the default profile %s or should start with %s.", podSeccompProfile, ProfileRuntimeDefault, ProfileNamePrefix),
			PendingFix: &fix.BySettingPodAnnotation{
				Key:   PodAnnotationKey,
				Value: ProfileRuntimeDefault,
			},
			Metadata: kubeaudit.Metadata{
				"AnnotationKey":   PodAnnotationKey,
				"AnnotationValue": podSeccompProfile,
			},
		}
	}

	return nil
}

func auditContainer(container *k8s.ContainerV1, resource k8s.Resource) *kubeaudit.AuditResult {
	annotations := k8s.GetAnnotations(resource)
	containerAnnotationKey := getContainerAnnotationKey(container)
	PodAnnotationKey := apiv1.SeccompPodAnnotationKey

	// Assume that the container will be covered by the pod-level seccomp annotation. If there is no pod-level
	// seccomp annotation, assume that it will be added as part of the pod annotation audit / autofix
	if isSeccompAnnotationMissing(annotations, containerAnnotationKey) {
		return nil
	}

	// If the pod seccomp profile is a custom profile, and the container seccomp profile is set to a bad value,
	// then set the container annotation to the default profile. Otherwise, if the container annotation is set to a
	// bad value, then remove the container annotation in favour of the pod annotation (assumes the pod annotation is
	// the default profile because even if the pod annotation is set to a bad value, it will be autofixed to be the
	// default profile)
	var pendingFix kubeaudit.PendingFix
	podSeccompProfile := annotations[PodAnnotationKey]
	if isSeccompProfileCustom(podSeccompProfile) {
		pendingFix = &fix.BySettingPodAnnotation{Key: containerAnnotationKey, Value: ProfileRuntimeDefault}
	} else {
		pendingFix = &fix.ByRemovingPodAnnotation{Key: containerAnnotationKey}
	}

	containerSeccompProfile := annotations[containerAnnotationKey]

	if isSeccompProfileDeprecated(containerSeccompProfile) {
		return &kubeaudit.AuditResult{
			Name:       SeccompDeprecatedContainer,
			Severity:   kubeaudit.Error,
			Message:    fmt.Sprintf("Seccomp container annotation is set to deprecated value %s. It should be set to %s instead.", containerSeccompProfile, ProfileRuntimeDefault),
			PendingFix: pendingFix,
			Metadata: kubeaudit.Metadata{
				"AnnotationKey":   containerAnnotationKey,
				"AnnotationValue": containerSeccompProfile,
			},
		}
	}

	if !isSeccompEnabled(containerSeccompProfile) {
		return &kubeaudit.AuditResult{
			Name:       SeccompDisabledContainer,
			Severity:   kubeaudit.Error,
			Message:    fmt.Sprintf("Seccomp container annotation is set to %s which disables Seccomp. It should be set to the default profile %s or should start with %s.", containerSeccompProfile, ProfileRuntimeDefault, ProfileNamePrefix),
			PendingFix: pendingFix,
			Metadata: kubeaudit.Metadata{
				"AnnotationKey":   containerAnnotationKey,
				"AnnotationValue": containerSeccompProfile,
			},
		}
	}

	return nil
}

func isSeccompAnnotationMissing(annotations map[string]string, annotationKey string) bool {
	_, ok := annotations[annotationKey]
	return !ok
}

// returns false if there is at least one container that is not covered by a container-level seccomp annotation
func isSeccompEnabledForContainers(annotations map[string]string, resource k8s.Resource) bool {
	for _, container := range k8s.GetContainers(resource) {
		containerAnnotationKey := getContainerAnnotationKey(container)
		if isSeccompAnnotationMissing(annotations, containerAnnotationKey) {
			return false
		}

		containerSeccompProfile := annotations[containerAnnotationKey]
		if !isSeccompEnabled(containerSeccompProfile) {
			return false
		}
	}

	return true
}

func isSeccompProfileDeprecated(seccompProfile string) bool {
	return seccompProfile == DeprecatedProfileRuntimeDefault
}

func isSeccompProfileCustom(seccompProfile string) bool {
	return strings.HasPrefix(seccompProfile, ProfileNamePrefix)
}

func isSeccompEnabled(seccompProfile string) bool {
	return isSeccompProfileDefault(seccompProfile) || isSeccompProfileCustom(seccompProfile)
}

func isSeccompProfileDefault(seccompProfile string) bool {
	return seccompProfile == ProfileRuntimeDefault
}

func getContainerAnnotationKey(container *k8s.ContainerV1) string {
	return ContainerAnnotationKeyPrefix + container.Name
}

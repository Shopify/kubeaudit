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
	// SeccompDeprecatedAnnotations occurs when deprecated seccomp annotations are present
	SeccompDeprecatedAnnotations = "SeccompDeprecatedAnnotations"
	// SeccompProfileMissing occurs when there are no seccomp profiles (pod nor container level)
	SeccompProfileMissing = "SeccompProfileMissing"
	// SeccompDisabledPod occurs when the pod-level seccomp profile is set to a value which disables seccomp
	SeccompDisabledPod = "SeccompDisabledPod"
	// SeccompDisabledContainer occurs when the container-level seccomp profile is set to a value which disables seccomp
	SeccompDisabledContainer = "SeccompDisabledContainer"
)

const (
	// ProfileRuntimeDefault represents the default seccomp profile used by container runtime
	ProfileRuntimeDefault = apiv1.SeccompProfileTypeRuntimeDefault
	// ProfileLocalhost represents the localhost seccomp profile used by container runtime
	ProfileLocalhost = apiv1.SeccompProfileTypeLocalhost
	// ContainerAnnotationKeyPrefix represents the key of a seccomp profile applied to one container of a pod
	ContainerAnnotationKeyPrefix = apiv1.SeccompContainerAnnotationKeyPrefix
	// PodAnnotationKey represents the key of a seccomp profile applied to all containers of a pod
	PodAnnotationKey = apiv1.SeccompPodAnnotationKey
)

// Seccomp implements Auditable
type Seccomp struct{}

func New() *Seccomp {
	return &Seccomp{}
}

// Audit checks that Seccomp is enabled for all containers
func (a *Seccomp) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	annotationAuditResult := auditAnnotations(resource)
	auditResults = appendNotNil(auditResults, annotationAuditResult)

	auditResult := auditPod(resource)
	auditResults = appendNotNil(auditResults, auditResult)

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container, resource)
		auditResults = appendNotNil(auditResults, auditResult)
	}

	return auditResults, nil
}

func appendNotNil(auditResults []*kubeaudit.AuditResult, auditResult *kubeaudit.AuditResult) []*kubeaudit.AuditResult {
	if auditResult != nil {
		return append(auditResults, auditResult)
	}
	return auditResults
}

func auditAnnotations(resource k8s.Resource) *kubeaudit.AuditResult {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	// We check annotations only when seccomp profile is missing for both pod and all containers
	// This way we ensure that we're in Manifest mode. 
	// In Local and Cluster mode Kubernetes automatically populates seccomp profile in Security context when seccomp annotations are provided.
	if !isPodSeccompProfileMissing(podSpec.SecurityContext) || !isSeccompProfileMissingForAllContainers(resource) {
		return nil
	}

	seccompAnnotations := findSeccompAnnotations(resource)

	if len(seccompAnnotations) > 0 {

		return &kubeaudit.AuditResult{
			Auditor:    Name,
			Rule:       SeccompDeprecatedAnnotations,
			Severity:   kubeaudit.Warn,
			Message:    "Pod Seccomp annotations are deprecated. Seccomp profile should be added to the pod SecurityContext.",
			PendingFix: &fix.ByRemovingPodAnnotations{Keys: seccompAnnotations},
			Metadata:   kubeaudit.Metadata{"AnnotationKeys": strings.Join(seccompAnnotations, ", ")},
		}
	}

	return nil
}

func auditPod(resource k8s.Resource) *kubeaudit.AuditResult {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if isPodSeccompProfileMissing(podSpec.SecurityContext) {
		// If all the containers have container-level seccomp profiles then we don't need a pod-level profile
		if isSeccompEnabledForAllContainers(resource) {
			return nil
		}

		return &kubeaudit.AuditResult{
			Auditor:    Name,
			Rule:       SeccompProfileMissing,
			Severity:   kubeaudit.Error,
			Message:    "Pod Seccomp profile is missing. Seccomp profile should be added to the pod SecurityContext.",
			PendingFix: &BySettingSeccompProfile{seccompProfileType: ProfileRuntimeDefault},
			Metadata:   kubeaudit.Metadata{},
		}
	}

	podSeccompProfileType := podSpec.SecurityContext.SeccompProfile.Type

	if !isSeccompEnabled(podSeccompProfileType) {
		return &kubeaudit.AuditResult{
			Auditor:    Name,
			Rule:       SeccompDisabledPod,
			Severity:   kubeaudit.Error,
			Message:    fmt.Sprintf("Pod Seccomp profile is set to %s which disables Seccomp. It should be set to the `%s` or `%s`.", podSeccompProfileType, ProfileRuntimeDefault, ProfileLocalhost),
			PendingFix: &BySettingSeccompProfile{seccompProfileType: ProfileRuntimeDefault},
			Metadata:   kubeaudit.Metadata{"SeccompProfileType": string(podSeccompProfileType)},
		}
	}

	return nil
}

func auditContainer(container *k8s.ContainerV1, resource k8s.Resource) *kubeaudit.AuditResult {
	// Assume that the container will be covered by the pod-level seccomp profile. If there is no pod-level
	// seccomp profile, assume that it will be added as part of the pod seccomp profile audit / autofix
	if isContainerSeccompProfileMissing(container.SecurityContext) {
		return nil
	}

	containerSeccompProfile := container.SecurityContext.SeccompProfile.Type
	if !isSeccompEnabled(containerSeccompProfile) {

		// If the pod seccomp profile is set to Localhost, and the container seccomp profile is disabled,
		// then set the container seccomp profile to the default profile.
		// Otherwise, remove the container seccomp profile in favour of the pod profile.
		var pendingFix kubeaudit.PendingFix
		var msg string

		podSpec := k8s.GetPodSpec(resource)
		if isPodSeccompProfileMissing(podSpec.SecurityContext) || isSeccompProfileDefault(podSpec.SecurityContext.SeccompProfile.Type) {
			pendingFix = &ByRemovingSeccompProfileInContainer{container: container}
			msg = fmt.Sprintf("Container Seccomp profile is set to %s which disables Seccomp. It should be removed from the container SecurityContext, as the pod SeccompProfile is set.", containerSeccompProfile)

		} else {
			pendingFix = &BySettingSeccompProfileInContainer{container: container, seccompProfileType: ProfileRuntimeDefault}
			msg = fmt.Sprintf("Container Seccomp profile is set to %s which disables Seccomp. It should be set to the `%s` or `%s`.", containerSeccompProfile, ProfileRuntimeDefault, ProfileLocalhost)
		}

		return &kubeaudit.AuditResult{
			Auditor:    Name,
			Rule:       SeccompDisabledContainer,
			Severity:   kubeaudit.Error,
			Message:    msg,
			PendingFix: pendingFix,
			Metadata: kubeaudit.Metadata{
				"Container":          container.Name,
				"SeccompProfileType": string(containerSeccompProfile),
			},
		}
	}

	return nil
}

func isPodSeccompProfileMissing(securityContext *apiv1.PodSecurityContext) bool {
	return securityContext == nil || securityContext.SeccompProfile == nil
}

func isContainerSeccompProfileMissing(securityContext *apiv1.SecurityContext) bool {
	return securityContext == nil || securityContext.SeccompProfile == nil
}

// returns false if there is at least one container that is not covered by a container-level seccomp annotation
func isSeccompEnabledForAllContainers(resource k8s.Resource) bool {
	for _, container := range k8s.GetContainers(resource) {
		securityContext := container.SecurityContext
		if isContainerSeccompProfileMissing(securityContext) {
			return false
		}

		containerSeccompProfileType := securityContext.SeccompProfile.Type
		if !isSeccompEnabled(containerSeccompProfileType) {
			return false
		}
	}

	return true
}

func isSeccompProfileMissingForAllContainers(resource k8s.Resource) bool {
	for _, container := range k8s.GetContainers(resource) {
		securityContext := container.SecurityContext
		if !isContainerSeccompProfileMissing(securityContext) {
			return false
		}
	}

	return true
}

func isSeccompEnabled(seccompProfileType apiv1.SeccompProfileType) bool {
	return isSeccompProfileDefault(seccompProfileType) || isSeccompProfileLocalhost(seccompProfileType)
}

func isSeccompProfileDefault(seccompProfileType apiv1.SeccompProfileType) bool {
	return seccompProfileType == apiv1.SeccompProfileTypeRuntimeDefault
}

func isSeccompProfileLocalhost(seccompProfileType apiv1.SeccompProfileType) bool {
	return seccompProfileType == apiv1.SeccompProfileTypeLocalhost
}

func findSeccompAnnotations(resource k8s.Resource) []string {
	annotations := k8s.GetAnnotations(resource)
	seccompAnnotations := []string{}
	for annotation := range annotations {
		if annotation == PodAnnotationKey || strings.HasPrefix(annotation, ContainerAnnotationKeyPrefix) {
			seccompAnnotations = append(seccompAnnotations, annotation)
		}
	}

	return seccompAnnotations
}

package seccomp

import (
	"fmt"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	apiv1 "k8s.io/api/core/v1"
)

const Name = "seccomp"

const (
	// SeccompProfileMissing occurs when there are no seccomp annotations (pod nor container level)
	SeccompProfileMissing = "SeccompProfileMissing"
	// SeccompDisabledPod occurs when the pod-level seccomp annotation is set to a value which disables seccomp
	SeccompDisabledPod = "SeccompDisabledPod"
	// SeccompDisabledContainer occurs when the container-level seccomp annotation is set to a value which disables seccomp
	SeccompDisabledContainer = "SeccompDisabledContainer"
)

const (
	// ProfileRuntimeDefault represents the default seccomp profile used by container runtime
	ProfileRuntimeDefault = apiv1.SeccompProfileTypeRuntimeDefault
	// ProfileLocalhost represents the localhost seccomp profile used by container runtime
	ProfileLocalhost = apiv1.SeccompProfileTypeLocalhost
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
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if isPodSeccompProfileMissing(podSpec.SecurityContext) {
		// If all the containers have container-level seccomp profiles then we don't need a pod-level annotation
		if isSeccompEnabledForContainers(resource) {
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

		// If the pod seccomp profile is set, and the container seccomp profile is disabled,
		// then remove the container seccomp profile in favour of the pod profile.
		// Otherwise, set the container seccomp profile to the default profile.
		var pendingFix kubeaudit.PendingFix
		var msg string

		podSpec := k8s.GetPodSpec(resource)
		if isPodSeccompProfileMissing(podSpec.SecurityContext) {
			pendingFix = &BySettingSeccompProfileInContainer{container: container, seccompProfileType: ProfileRuntimeDefault}
			msg = fmt.Sprintf("Container Seccomp profile is set to %s which disables Seccomp. It should be set to the `%s` or `%s`.", containerSeccompProfile, ProfileRuntimeDefault, ProfileLocalhost)
		} else {
			pendingFix = &ByRemovingSeccompProfileInContainer{container: container}
			msg = fmt.Sprintf("Container Seccomp profile is set to %s which disables Seccomp. It should be removed from the container SecurityContext, as the pod SeccompProfile is set.", containerSeccompProfile)
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
func isSeccompEnabledForContainers(resource k8s.Resource) bool {
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

func isSeccompEnabled(seccompProfileType apiv1.SeccompProfileType) bool {
	return isSeccompProfileDefault(seccompProfileType) || isSeccompProfileLocalhost(seccompProfileType)
}

func isSeccompProfileDefault(seccompProfileType apiv1.SeccompProfileType) bool {
	return seccompProfileType == apiv1.SeccompProfileTypeRuntimeDefault
}

func isSeccompProfileLocalhost(seccompProfileType apiv1.SeccompProfileType) bool {
	return seccompProfileType == apiv1.SeccompProfileTypeLocalhost
}

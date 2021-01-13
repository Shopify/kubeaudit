package nonroot

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const Name = "nonroot"

const (
	// RunAsUserCSCRoot occurs when runAsUser is set to 0 in the container SecurityContext
	RunAsUserCSCRoot = "RunAsUserCSCRoot"
	// RunAsUserPSCRoot occurs when runAsUser is set to 0 in the pod SecurityContext
	RunAsUserPSCRoot = "RunAsUserPSCRoot"
	// RunAsNonRootCSCFalse occurs when runAsNonRoot is set to false in the container SecurityContext
	RunAsNonRootCSCFalse = "RunAsNonRootCSCFalse"
	// RunAsNonRootPSCNilCSCNil occurs when runAsNonRoot is not set in the container SecurityContext nor the pod
	// security context. runAsNonRoot defaults to false so this is bad
	RunAsNonRootPSCNilCSCNil = "RunAsNonRootPSCNilCSCNil"
	// RunAsNonRootPSCFalseCSCNil occurs when runAsNonRoot is not set in the container SecurityContext and is set to
	// false in the PodSecurityContext
	RunAsNonRootPSCFalseCSCNil = "RunAsNonRootPSCFalseCSCNil"
)

const OverrideLabel = "allow-run-as-root"

// RunAsNonRoot implements Auditable
type RunAsNonRoot struct{}

func New() *RunAsNonRoot {
	return &RunAsNonRoot{}
}

// Audit checks that runAsNonRoot is set to true in every container's security context
func (a *RunAsNonRoot) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container, resource)
		auditResult = override.ApplyOverride(auditResult, container.Name, resource, OverrideLabel)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditContainer(container *k8stypes.ContainerV1, resource k8stypes.Resource) *kubeaudit.AuditResult {
	if !isContainerRunAsUserNil(container) {
		if *container.SecurityContext.RunAsUser == 0 {
			return &kubeaudit.AuditResult{
				Name:     RunAsUserCSCRoot,
				Severity: kubeaudit.Error,
				Message:  "container running with root user (UID 0) as configured with runAsUser SecurityContext. It should be set to > 0.",
				Metadata: kubeaudit.Metadata{
					"Container": container.Name,
				},
			}
		}

		return nil
	}

	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if !isPodRunAsUserNil(podSpec) {
		if *podSpec.SecurityContext.RunAsUser == 0 {
			return &kubeaudit.AuditResult{
				Name:     RunAsUserPSCRoot,
				Severity: kubeaudit.Error,
				Message:  "runAsUser is not set in container SecurityContext and is set to 0 in the PodSecurityContext. It should be set to > 0 in at least one of the two.",
				Metadata: kubeaudit.Metadata{
					"Container": container.Name,
				},
			}
		}

		return nil
	}

	if isContainerRunAsNonRootCSCFalse(container) {
		return &kubeaudit.AuditResult{
			Name:     RunAsNonRootCSCFalse,
			Severity: kubeaudit.Error,
			Message:  "runAsNonRoot is set to false and runAsUser is not set in container SecurityContext. Either the first should be set to true or the later set to a non-root UID.",
			PendingFix: &fixRunAsNonRoot{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	if isContainerRunAsNonRootNil(container) {
		if isPodRunAsNonRootNil(podSpec) {
			return &kubeaudit.AuditResult{
				Name:     RunAsNonRootPSCNilCSCNil,
				Severity: kubeaudit.Error,
				Message:  "runAsNonRoot should be set to 'true' or runAsUser should be set to a non-root UID either in the container SecurityContext or PodSecurityContext.",
				PendingFix: &fixRunAsNonRoot{
					container: container,
				},
				Metadata: kubeaudit.Metadata{
					"Container": container.Name,
				},
			}
		}

		if isPodRunAsNonRootFalse(podSpec) {
			return &kubeaudit.AuditResult{
				Name:     RunAsNonRootPSCFalseCSCNil,
				Severity: kubeaudit.Error,
				Message:  "runAsNonRoot is not set in container SecurityContext but is set to false in the PodSecurityContext while runAsUser is not set. runAsNonRoot should be set to 'true' or runAsUser should be set to a non-root UID either in container SecurityContext or PodSecurityContext.",
				PendingFix: &fixRunAsNonRoot{
					container: container,
				},
				Metadata: kubeaudit.Metadata{
					"Container": container.Name,
				},
			}
		}
	}

	return nil
}

// returns true if runAsNonRoot is explicitly set to false in the pod's security context. Returns true if the
// security context is nil even though the default value for runAsNonRoot is false
func isPodRunAsNonRootFalse(podSpec *k8stypes.PodSpecV1) bool {
	if isPodRunAsNonRootNil(podSpec) {
		return false
	}

	return !*podSpec.SecurityContext.RunAsNonRoot
}

func isPodRunAsNonRootNil(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.SecurityContext == nil || podSpec.SecurityContext.RunAsNonRoot == nil
}

// returns true if runAsNonRoot is explicitly set to false in the container's security context. Returns true if the
// security context is nil even though the default value for runAsNonRoot is false
func isContainerRunAsNonRootCSCFalse(container *k8stypes.ContainerV1) bool {
	if isContainerRunAsNonRootNil(container) {
		return false
	}

	return !*container.SecurityContext.RunAsNonRoot
}

func isContainerRunAsNonRootNil(container *k8stypes.ContainerV1) bool {
	return container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil
}

func isContainerRunAsUserNil(container *k8stypes.ContainerV1) bool {
	return container.SecurityContext == nil || container.SecurityContext.RunAsUser == nil
}

func isPodRunAsUserNil(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.SecurityContext == nil || podSpec.SecurityContext.RunAsUser == nil
}

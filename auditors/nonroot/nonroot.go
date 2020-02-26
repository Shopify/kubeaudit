package nonroot

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const Name = "nonroot"

const (
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
	if isContainerRunAsNonRootCSCFalse(container) {
		return &kubeaudit.AuditResult{
			Name:     RunAsNonRootCSCFalse,
			Severity: kubeaudit.Error,
			Message:  "runAsNonRoot is set to false in container SecurityContext. It should be set to true.",
			PendingFix: &fixRunAsNonRoot{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}
	if isContainerRunAsNonRootNil(container) {
		if isPodRunAsNonRootNil(podSpec) {
			return &kubeaudit.AuditResult{
				Name:     RunAsNonRootPSCNilCSCNil,
				Severity: kubeaudit.Error,
				Message:  "runAsNonRoot is not set in container SecurityContext nor the PodSecurityContext. It should be set to 'true' in at least one of the two.",
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
				Message:  "runAsNonRoot is not set in container SecurityContext and is set to false in the PodSecurityContext. It should be set to 'true' in at least one of the two.",
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

// returns true if runAsNonRoot is explicilty set to false in the pod's security context. Returns true if the
// security context is nil even though the default value for runAsNonRoot is false
func isPodRunAsNonRootFalse(podSpec *k8stypes.PodSpecV1) bool {
	if isPodRunAsNonRootNil(podSpec) {
		return false
	}

	return !*podSpec.SecurityContext.RunAsNonRoot
}

func isPodRunAsNonRootNil(podSpec *k8stypes.PodSpecV1) bool {
	if podSpec.SecurityContext == nil || podSpec.SecurityContext.RunAsNonRoot == nil {
		return true
	}

	return false
}

// returns true if runAsNonRoot is explicilty set to false in the containers's security context. Returns true if the
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

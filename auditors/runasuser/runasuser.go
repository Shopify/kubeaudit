package runasuser

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const Name = "runasuser"

const (
	// RunAsUserCSCRoot occurs when runAsUser is set to 0 in the container SecurityContext
	RunAsUserCSCRoot = "RunAsUserCSCRoot"
	// RunAsUserPSCNilCSCNil occurs when runAsUser is not set in the container SecurityContext nor the pod
	// security context. runAsUser defaults to 0 so this is bad
	RunAsUserPSCNilCSCNil = "RunAsUserPSCNilCSCNil"
	// RunAsUserPSCRootCSCNil occurs when runAsUser is not set in the container SecurityContext and is set to
	// 0 in the PodSecurityContext
	RunAsUserPSCRootCSCNil = "RunAsUserPSCRootCSCNil"
)

const OverrideLabel = "allow-not-overridden-non-root-user"

// RunAsUser implements Auditable
type RunAsUser struct{}

func New() *RunAsUser {
	return &RunAsUser{}
}

// Audit checks that runAsUser is set to 0 in every container's security context
func (a *RunAsUser) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
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
	if isContainerRunAsUserSetToNonRoot(container) {
		return &kubeaudit.AuditResult{
			Name:     RunAsUserCSCRoot,
			Severity: kubeaudit.Error,
			Message:  "container user ID not overridden to non-root user using runAsUser SecurityContext. It should be set to > 0.",
			PendingFix: &fixRunAsUser{
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
	if isContainerRunAsUserNil(container) {
		if isPodRunAsUserNil(podSpec) {
			return &kubeaudit.AuditResult{
				Name:     RunAsUserPSCNilCSCNil,
				Severity: kubeaudit.Error,
				Message:  "runAsUser is not set in container SecurityContext nor the PodSecurityContext. It should be set to > 0 in at least one of the two.",
				PendingFix: &fixRunAsUser{
					container: container,
				},
				Metadata: kubeaudit.Metadata{
					"Container": container.Name,
				},
			}
		}

		if isPodRunAsUserSetToNonRoot(podSpec) {
			return &kubeaudit.AuditResult{
				Name:     RunAsUserPSCRootCSCNil,
				Severity: kubeaudit.Error,
				Message:  "runAsUser is not set in container SecurityContext and is set to 0 in the PodSecurityContext. It should be set to > 0 in at least one of the two.",
				PendingFix: &fixRunAsUser{
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

// returns true if runAsUser is explicilty set to 0 in the pod's security context. Returns true if the
// security context is nil even though the default value for runAsUser is 0
func isPodRunAsUserSetToNonRoot(podSpec *k8stypes.PodSpecV1) bool {
	if isPodRunAsUserNil(podSpec) {
		return false
	}

	return *podSpec.SecurityContext.RunAsUser == 0
}

func isPodRunAsUserNil(podSpec *k8stypes.PodSpecV1) bool {
	if podSpec.SecurityContext == nil || podSpec.SecurityContext.RunAsUser == nil {
		return true
	}

	return false
}

// returns true if runAsUser is explicilty set to 0 in the containers's security context. Returns true if the
// security context is nil even though the default value for runAsUser is 0
func isContainerRunAsUserSetToNonRoot(container *k8stypes.ContainerV1) bool {
	if isContainerRunAsUserNil(container) {
		return false
	}

	return *container.SecurityContext.RunAsUser == 0
}

func isContainerRunAsUserNil(container *k8stypes.ContainerV1) bool {
	return container.SecurityContext == nil || container.SecurityContext.RunAsUser == nil
}

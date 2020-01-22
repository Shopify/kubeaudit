package privileged

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const (
	// PrivilegedTrue occurs when privileged is set to true in the container SecurityContext
	PrivilegedTrue = "PrivilegedTrue"
	// PrivilegedNil occurs when privileged is not set in the container SecurityContext.
	// Prvivileged defaults to false so this is ok
	PrivilegedNil = "PrivilegedNil"
)

const OverrideLabel = "allow-privileged"

// Privileged implements Auditable
type Privileged struct{}

func New() *Privileged {
	return &Privileged{}
}

// Audit checks that privileged is set to false in every container's security context
func (a *Privileged) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
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
	if isPrivilegedNil(container) {
		// TODO:  should "warning" results include a fix?
		return &kubeaudit.AuditResult{
			Name:     PrivilegedNil,
			Severity: kubeaudit.Warn,
			Message:  "privileged is not set in container SecurityContext. Privileged defaults to 'false' but it should be explicitly set to 'false'.",
			PendingFix: &fixPrivileged{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	if isPrivilegedTrue(container) {
		return &kubeaudit.AuditResult{
			Name:     PrivilegedTrue,
			Severity: kubeaudit.Error,
			Message:  "privileged is set to 'true' in container SecurityContext. It should be set to 'false'.",
			PendingFix: &fixPrivileged{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	return nil
}

func isPrivilegedTrue(container *k8stypes.ContainerV1) bool {
	if isPrivilegedNil(container) {
		return false
	}

	return *container.SecurityContext.Privileged
}

func isPrivilegedNil(container *k8stypes.ContainerV1) bool {
	return container.SecurityContext == nil || container.SecurityContext.Privileged == nil
}

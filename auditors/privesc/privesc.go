package privesc

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const Name = "privesc"

const (
	// AllowPrivilegeEscalationNil occurs when the AllowPrivilegeEscalation field is missing or unset in the
	// container SecurityContext
	AllowPrivilegeEscalationNil = "AllowPrivilegeEscalationNil"
	// AllowPrivilegeEscalationTrue occurs when the AllowPrivilegeEscalation field is set to true in the container
	// security context
	AllowPrivilegeEscalationTrue = "AllowPrivilegeEscalationTrue"
)

const OverrideLabel = "allow-privilege-escalation"

// AllowPrivilegeEscalation implements Auditable
type AllowPrivilegeEscalation struct{}

func New() *AllowPrivilegeEscalation {
	return &AllowPrivilegeEscalation{}
}

// Audit checks that AllowPrivilegeEscalation is disabled (set to false) in the container SecurityContext
func (a *AllowPrivilegeEscalation) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container)
		auditResult = override.ApplyOverride(auditResult, container.Name, resource, OverrideLabel)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditContainer(container *k8s.ContainerV1) *kubeaudit.AuditResult {
	if isAllowPrivilegeEscalationNil(container) {
		return &kubeaudit.AuditResult{
			Name:     AllowPrivilegeEscalationNil,
			Severity: kubeaudit.Error,
			Message:  "allowPrivilegeEscalation not set which allows privilege escalation. It should be set to 'false'.",
			PendingFix: &fixBySettingAllowPrivilegeEscalationFalse{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	if isAllowPrivilegeEscalationTrue(container) {
		return &kubeaudit.AuditResult{
			Name:     AllowPrivilegeEscalationTrue,
			Severity: kubeaudit.Error,
			Message:  "allowPrivilegeEscalation set to 'true'. It should be set to 'false'.",
			PendingFix: &fixBySettingAllowPrivilegeEscalationFalse{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	return nil
}

func isAllowPrivilegeEscalationNil(container *k8s.ContainerV1) bool {
	return container.SecurityContext == nil || container.SecurityContext.AllowPrivilegeEscalation == nil
}

func isAllowPrivilegeEscalationTrue(container *k8s.ContainerV1) bool {
	return container.SecurityContext != nil && *container.SecurityContext.AllowPrivilegeEscalation
}

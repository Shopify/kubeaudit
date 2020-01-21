package asat

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const (
	// AutomountServiceAccountTokenDeprecated occurs when the deprecated serviceAccount field is non-empty
	AutomountServiceAccountTokenDeprecated = "AutomountServiceAccountTokenDeprecated"
	// AutomountServiceAccountTokenTrueAndDefaultSA occurs when automountServiceAccountToken is either not set
	// (which defaults to true) or explicitly set to true, and serviceAccountName is either not set or set to "default"
	AutomountServiceAccountTokenTrueAndDefaultSA = "AutomountServiceAccountTokenTrueAndDefaultSA"
)

const OverrideLabel = "allow-automount-service-account-token"

// AutomountServiceAccountToken implements Auditable
type AutomountServiceAccountToken struct{}

func New() *AutomountServiceAccountToken {
	return &AutomountServiceAccountToken{}
}

// Audit checks that the deprecated serviceAccount field is not used and that the default service account is not
// being automatically mounted
func (a *AutomountServiceAccountToken) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	auditResult := auditResource(resource)
	auditResult = override.ApplyOverride(auditResult, "", resource, OverrideLabel)
	if auditResult != nil {
		return []*kubeaudit.AuditResult{auditResult}, nil
	}

	return nil, nil
}

func auditResource(resource k8stypes.Resource) *kubeaudit.AuditResult {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if isDeprecatedServiceAccountName(podSpec) {
		return &kubeaudit.AuditResult{
			Name:     AutomountServiceAccountTokenDeprecated,
			Severity: kubeaudit.Warn,
			Message:  "serviceAccount is a deprecated alias for ServiceAccountName. serviceAccountName should be used instead.",
			PendingFix: &fixDeprecatedServiceAccountName{
				podSpec: podSpec,
			},
			Metadata: kubeaudit.Metadata{
				"DeprecatedServiceAccount": podSpec.DeprecatedServiceAccount,
			},
		}
	}

	if isDefaultServiceAccountWithAutomountToken(podSpec) {
		return &kubeaudit.AuditResult{
			Name:     AutomountServiceAccountTokenTrueAndDefaultSA,
			Severity: kubeaudit.Error,
			Message:  "Default serviceAccount with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.",
			PendingFix: &fixDefaultServiceAccountWithAutomountToken{
				podSpec: podSpec,
			},
		}
	}

	return nil
}

func isDeprecatedServiceAccountName(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.DeprecatedServiceAccount != ""
}

func isDefaultServiceAccountWithAutomountToken(podSpec *k8stypes.PodSpecV1) bool {
	return isAutomountTokenTrue(podSpec) && isDefaultServiceAccount(podSpec)
}

func isAutomountTokenTrue(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.AutomountServiceAccountToken == nil || *podSpec.AutomountServiceAccountToken
}

func isDefaultServiceAccount(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.ServiceAccountName == "" || podSpec.ServiceAccountName == "default"
}

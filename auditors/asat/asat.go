package asat

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const Name = "asat"

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
func (a *AutomountServiceAccountToken) Audit(resource k8stypes.Resource, resources []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	auditResult := auditResource(resource, resources)
	auditResult = override.ApplyOverride(auditResult, "", resource, OverrideLabel)
	if auditResult != nil {
		return []*kubeaudit.AuditResult{auditResult}, nil
	}

	return nil, nil
}

func auditResource(resource k8stypes.Resource, resources []k8stypes.Resource) *kubeaudit.AuditResult {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if isDeprecatedServiceAccountName(podSpec) && !hasServiceAccountName(podSpec) {
		return &kubeaudit.AuditResult{
			Name:     AutomountServiceAccountTokenDeprecated,
			Severity: kubeaudit.Warn,
			Message:  "serviceAccount is a deprecated alias for serviceAccountName. serviceAccountName should be used instead.",
			PendingFix: &fixDeprecatedServiceAccountName{
				podSpec: podSpec,
			},
			Metadata: kubeaudit.Metadata{
				"DeprecatedServiceAccount": podSpec.DeprecatedServiceAccount,
			},
		}
	}

	defaultServiceAccount := getDefaultServiceAccount(resources)
	if usesDefaultServiceAccount(podSpec) && isAutomountTokenTrue(podSpec, defaultServiceAccount) {
		return &kubeaudit.AuditResult{
			Name:     AutomountServiceAccountTokenTrueAndDefaultSA,
			Severity: kubeaudit.Error,
			Message:  "Default service account with token mounted. automountServiceAccountToken should be set to 'false' on either the ServiceAccount or on the PodSpec or a non-default service account should be used.",
			PendingFix: &fixDefaultServiceAccountWithAutomountToken{
				podSpec:               podSpec,
				defaultServiceAccount: defaultServiceAccount,
			},
		}
	}

	return nil
}

func isDeprecatedServiceAccountName(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.DeprecatedServiceAccount != ""
}

func hasServiceAccountName(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.ServiceAccountName != ""
}

func isAutomountTokenTrue(podSpec *k8stypes.PodSpecV1, defaultServiceAccount *k8stypes.ServiceAccountV1) bool {
	if podSpec.AutomountServiceAccountToken != nil {
		return *podSpec.AutomountServiceAccountToken
	}

	return defaultServiceAccount == nil ||
		defaultServiceAccount.AutomountServiceAccountToken == nil ||
		*defaultServiceAccount.AutomountServiceAccountToken
}

func usesDefaultServiceAccount(podSpec *k8stypes.PodSpecV1) bool {
	return podSpec.ServiceAccountName == "" || podSpec.ServiceAccountName == "default"
}

func getDefaultServiceAccount(resources []k8stypes.Resource) (serviceAccount *k8stypes.ServiceAccountV1) {
	for _, resource := range resources {
		serviceAccount, ok := resource.(*k8stypes.ServiceAccountV1)
		if ok && (k8s.GetObjectMeta(serviceAccount).GetName() == "default") {
			return serviceAccount
		}
	}
	return
}

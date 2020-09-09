package asat

import (
	"fmt"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
)

type fixDeprecatedServiceAccountName struct {
	podSpec *k8stypes.PodSpecV1
}

func (f *fixDeprecatedServiceAccountName) Plan() string {
	return fmt.Sprintf("Set serviceAccountName to '%s' and set serviceAccount to '' in PodSpec", f.podSpec.DeprecatedServiceAccount)
}

func (f *fixDeprecatedServiceAccountName) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	f.podSpec.ServiceAccountName = f.podSpec.DeprecatedServiceAccount
	f.podSpec.DeprecatedServiceAccount = ""
	return nil
}

type fixDefaultServiceAccountWithAutomountToken struct {
	podSpec               *k8stypes.PodSpecV1
	defaultServiceAccount *k8stypes.ServiceAccountV1
}

func (f *fixDefaultServiceAccountWithAutomountToken) Plan() string {
	if f.defaultServiceAccount != nil {
		plan := "Set automountServiceAccountToken to 'false' in ServiceAccount"
		if f.podSpec.AutomountServiceAccountToken != nil && *(f.podSpec.AutomountServiceAccountToken) {
			plan += " and set automountServiceAccountToken to 'nil' in PodSpec"
		}
		return plan
	}
	return "Set automountServiceAccountToken to 'false' in PodSpec"
}

func (f *fixDefaultServiceAccountWithAutomountToken) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	if f.defaultServiceAccount != nil {
		f.defaultServiceAccount.AutomountServiceAccountToken = k8s.NewFalse()
		if (f.podSpec.AutomountServiceAccountToken != nil) && *(f.podSpec.AutomountServiceAccountToken) {
			f.podSpec.AutomountServiceAccountToken = nil
		}
	} else {
		f.podSpec.AutomountServiceAccountToken = k8s.NewFalse()
	}
	return nil
}

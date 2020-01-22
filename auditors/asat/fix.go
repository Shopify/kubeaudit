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
	podSpec *k8stypes.PodSpecV1
}

func (f *fixDefaultServiceAccountWithAutomountToken) Plan() string {
	return fmt.Sprintf("Set automountServiceAccountToken to 'false' in PodSpec")
}

func (f *fixDefaultServiceAccountWithAutomountToken) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	f.podSpec.AutomountServiceAccountToken = k8s.NewFalse()
	return nil
}

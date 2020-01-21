package all

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
)

func Auditors() []kubeaudit.Auditable {
	return []kubeaudit.Auditable{
		apparmor.New(),
	}
}

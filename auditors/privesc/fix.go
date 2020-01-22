package privesc

import (
	"fmt"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
)

type fixBySettingAllowPrivilegeEscalationFalse struct {
	container *k8stypes.ContainerV1
}

func (f *fixBySettingAllowPrivilegeEscalationFalse) Plan() string {
	return fmt.Sprintf("Set AllowPrivilegeEscalation to 'false' in the container SecurityContext for container %s", f.container.Name)
}

func (f *fixBySettingAllowPrivilegeEscalationFalse) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8stypes.SecurityContextV1{}
	}
	f.container.SecurityContext.AllowPrivilegeEscalation = k8s.NewFalse()
	return nil
}

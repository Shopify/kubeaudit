package privesc

import (
	"fmt"

	"github.com/Shopify/kubeaudit/pkg/k8s"
)

type fixBySettingAllowPrivilegeEscalationFalse struct {
	container *k8s.ContainerV1
}

func (f *fixBySettingAllowPrivilegeEscalationFalse) Plan() string {
	return fmt.Sprintf("Set AllowPrivilegeEscalation to 'false' in the container SecurityContext for container %s", f.container.Name)
}

func (f *fixBySettingAllowPrivilegeEscalationFalse) Apply(resource k8s.Resource) []k8s.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8s.SecurityContextV1{}
	}
	f.container.SecurityContext.AllowPrivilegeEscalation = k8s.NewFalse()
	return nil
}

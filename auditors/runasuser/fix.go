package runasuser

import (
	"fmt"
	"github.com/Shopify/kubeaudit/k8stypes"
)

type fixRunAsUser struct {
	container *k8stypes.ContainerV1
}

func (f *fixRunAsUser) Plan() string {
	return fmt.Sprintf("Set runAsUser to 1 in container SecurityContext for container %s", f.container.Name)
}

func (f *fixRunAsUser) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8stypes.SecurityContextV1{}
	}
	var uid int64 = 1
	f.container.SecurityContext.RunAsUser = &uid
	return nil
}

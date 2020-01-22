package privileged

import (
	"fmt"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
)

type fixPrivileged struct {
	container *k8stypes.ContainerV1
}

func (f *fixPrivileged) Plan() string {
	return fmt.Sprintf("Set privileged to 'false' in container SecurityContext for container %s", f.container.Name)
}

func (f *fixPrivileged) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8stypes.SecurityContextV1{}
	}
	f.container.SecurityContext.Privileged = k8s.NewFalse()
	return nil
}

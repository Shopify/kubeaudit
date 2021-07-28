package privileged

import (
	"fmt"

	"github.com/Shopify/kubeaudit/pkg/k8s"
)

type fixPrivileged struct {
	container *k8s.ContainerV1
}

func (f *fixPrivileged) Plan() string {
	return fmt.Sprintf("Set privileged to 'false' in container SecurityContext for container %s", f.container.Name)
}

func (f *fixPrivileged) Apply(resource k8s.Resource) []k8s.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8s.SecurityContextV1{}
	}
	f.container.SecurityContext.Privileged = k8s.NewFalse()
	return nil
}

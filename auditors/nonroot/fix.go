package nonroot

import (
	"fmt"

	"github.com/Shopify/kubeaudit/pkg/k8s"
)

type fixRunAsNonRoot struct {
	container *k8s.ContainerV1
}

func (f *fixRunAsNonRoot) Plan() string {
	return fmt.Sprintf("Set runAsNonRoot to 'true' in container SecurityContext for container %s", f.container.Name)
}

func (f *fixRunAsNonRoot) Apply(resource k8s.Resource) []k8s.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8s.SecurityContextV1{}
	}

	if f.container.SecurityContext.RunAsUser != nil && *f.container.SecurityContext.RunAsUser == 0 {
		f.container.SecurityContext.RunAsUser = nil
	}

	f.container.SecurityContext.RunAsNonRoot = k8s.NewTrue()
	return nil
}

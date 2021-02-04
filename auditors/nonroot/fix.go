package nonroot

import (
	"fmt"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
)

type fixRunAsNonRoot struct {
	container *k8stypes.ContainerV1
}

func (f *fixRunAsNonRoot) Plan() string {
	return fmt.Sprintf("Set runAsNonRoot to 'true' in container SecurityContext for container %s", f.container.Name)
}

func (f *fixRunAsNonRoot) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8stypes.SecurityContextV1{}
	}

	if f.container.SecurityContext.RunAsUser != nil && *f.container.SecurityContext.RunAsUser == 0 {
		f.container.SecurityContext.RunAsUser = nil
	}

	f.container.SecurityContext.RunAsNonRoot = k8s.NewTrue()
	return nil
}

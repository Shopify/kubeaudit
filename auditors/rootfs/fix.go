package rootfs

import (
	"fmt"

	"github.com/Shopify/kubeaudit/pkg/k8s"
)

type fixReadOnlyRootFilesystem struct {
	container *k8s.ContainerV1
}

func (f *fixReadOnlyRootFilesystem) Plan() string {
	return fmt.Sprintf("Set readOnlyRootFilesystem to 'true' in container SecurityContext for container %s", f.container.Name)
}

func (f *fixReadOnlyRootFilesystem) Apply(resource k8s.Resource) []k8s.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8s.SecurityContextV1{}
	}
	f.container.SecurityContext.ReadOnlyRootFilesystem = k8s.NewTrue()
	return nil
}

package rootfs

import (
	"fmt"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
)

type fixReadOnlyRootFilesystem struct {
	container *k8stypes.ContainerV1
}

func (f *fixReadOnlyRootFilesystem) Plan() string {
	return fmt.Sprintf("Set readOnlyRootFilesystem to 'true' in container SecurityContext for container %s", f.container.Name)
}

func (f *fixReadOnlyRootFilesystem) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	if f.container.SecurityContext == nil {
		f.container.SecurityContext = &k8stypes.SecurityContextV1{}
	}
	f.container.SecurityContext.ReadOnlyRootFilesystem = k8s.NewTrue()
	return nil
}

package capabilities

import (
	"fmt"

	"github.com/Shopify/kubeaudit/k8stypes"
)

type fixCapabilityAdded struct {
	container  *k8stypes.ContainerV1
	capability string
}

func (f *fixCapabilityAdded) Plan() string {
	return fmt.Sprintf("Remove capability '%s' from the capability add list in the container SecurityContext for container %s", f.capability, f.container.Name)
}

func (f *fixCapabilityAdded) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	removeCapabilityFromAddList(f.container, f.capability)
	return nil
}

type fixCapabilityNotDropped struct {
	container  *k8stypes.ContainerV1
	capability string
}

func (f *fixCapabilityNotDropped) Plan() string {
	return fmt.Sprintf("Add capability '%s' to the capability drop list in the container SecurityContext for container %s", f.capability, f.container.Name)
}

func (f *fixCapabilityNotDropped) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	dropCapability(f.container, f.capability)
	return nil
}

func dropCapability(container *k8stypes.ContainerV1, capability string) {
	if container.SecurityContext == nil {
		container.SecurityContext = &k8stypes.SecurityContextV1{}
	}

	if container.SecurityContext.Capabilities == nil {
		container.SecurityContext.Capabilities = &k8stypes.CapabilitiesV1{}
	}

	if container.SecurityContext.Capabilities.Drop == nil {
		container.SecurityContext.Capabilities.Drop = []k8stypes.CapabilityV1{}
	}

	dropped := container.SecurityContext.Capabilities.Drop
	dropped = append(dropped, k8stypes.CapabilityV1(capability))
	container.SecurityContext.Capabilities.Drop = dropped
}

func removeCapabilityFromAddList(container *k8stypes.ContainerV1, capability string) {
	added := container.SecurityContext.Capabilities.Add
	for i, add := range added {
		if string(add) == capability {
			added = append(added[:i], added[i+1:]...)
			break
		}
	}

	container.SecurityContext.Capabilities.Add = added
}

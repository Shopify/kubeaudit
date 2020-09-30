package capabilities

import (
	"fmt"

	"github.com/Shopify/kubeaudit/k8stypes"
	v1 "k8s.io/api/core/v1"
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

type fixCapabilityNotDroppedAll struct {
	container  *k8stypes.ContainerV1
	capability string
}

func (f *fixCapabilityNotDroppedAll) Plan() string {
	return fmt.Sprintf("Remove '%s' capability from drop list in the container SecurityContext for container %s", f.capability, f.container.Name)
}

func (f *fixCapabilityNotDroppedAll) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	dropCapabilityFromDropList(f.container, f.capability)
	return nil
}

func dropCapabilityFromDropList(container *k8stypes.ContainerV1, capability string) {
	if container.SecurityContext == nil {
		container.SecurityContext = &k8stypes.SecurityContextV1{}
	}

	if container.SecurityContext.Capabilities == nil {
		container.SecurityContext.Capabilities = &k8stypes.CapabilitiesV1{}
	}

	if container.SecurityContext.Capabilities.Drop == nil {
		container.SecurityContext.Capabilities.Drop = []k8stypes.CapabilityV1{}
	}

	container.SecurityContext.Capabilities.Drop = []v1.Capability{"ALL"}
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

type fixMissingSecurityContextOrCapability struct {
	container *k8stypes.ContainerV1
}

func (f *fixMissingSecurityContextOrCapability) Plan() string {
	return fmt.Sprintf("Adds security context and capabilities to %s. The capabilities Drop list is set to ALL.", f.container.Name)
}

func (f *fixMissingSecurityContextOrCapability) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	setDropToAll(f.container)
	return nil
}

func setDropToAll(container *k8stypes.ContainerV1) {
	if container.SecurityContext == nil {
		container.SecurityContext = &k8stypes.SecurityContextV1{}
	}

	if container.SecurityContext.Capabilities == nil {
		container.SecurityContext.Capabilities = &k8stypes.CapabilitiesV1{}
	}

	if container.SecurityContext.Capabilities.Drop == nil {
		container.SecurityContext.Capabilities.Drop = []k8stypes.CapabilityV1{}
	}

	container.SecurityContext.Capabilities.Drop = []v1.Capability{"ALL"}
}

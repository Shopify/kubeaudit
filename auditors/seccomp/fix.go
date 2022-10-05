package seccomp

import (
	"fmt"

	"github.com/Shopify/kubeaudit/pkg/k8s"
	apiv1 "k8s.io/api/core/v1"
)

type BySettingSeccompProfile struct {
	seccompProfileType apiv1.SeccompProfileType
}

func (pending *BySettingSeccompProfile) Plan() string {
	return fmt.Sprintf("Set SeccompProfile type to '%s' in pod SecurityContext", pending.seccompProfileType)
}

func (pending *BySettingSeccompProfile) Apply(resource k8s.Resource) []k8s.Resource {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec.SecurityContext == nil {
		podSpec.SecurityContext = &apiv1.PodSecurityContext{}
	}
	podSpec.SecurityContext.SeccompProfile = &apiv1.SeccompProfile{Type: pending.seccompProfileType}

	return nil
}

type BySettingSeccompProfileInContainer struct {
	container          *k8s.ContainerV1
	seccompProfileType apiv1.SeccompProfileType
}

func (pending *BySettingSeccompProfileInContainer) Plan() string {
	return fmt.Sprintf("Set SeccompProfile type to '%s' in SecurityContext for container `%s`", pending.seccompProfileType, pending.container.Name)
}

func (pending *BySettingSeccompProfileInContainer) Apply(resource k8s.Resource) []k8s.Resource {
	if pending.container.SecurityContext == nil {
		pending.container.SecurityContext = &apiv1.SecurityContext{}
	}
	pending.container.SecurityContext.SeccompProfile = &apiv1.SeccompProfile{Type: pending.seccompProfileType}
	return nil
}

type ByRemovingSeccompProfileInContainer struct {
	container *k8s.ContainerV1
}

func (pending *ByRemovingSeccompProfileInContainer) Plan() string {
	return fmt.Sprintf("Remove SeccompProfile in SecurityContext for container `%s`", pending.container.Name)
}

func (pending *ByRemovingSeccompProfileInContainer) Apply(resource k8s.Resource) []k8s.Resource {
	if pending.container.SecurityContext == nil {
		return nil
	}
	pending.container.SecurityContext.SeccompProfile = nil
	return nil
}

package mounts

import (
	"fmt"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/Shopify/kubeaudit/pkg/override"
	v1 "k8s.io/api/core/v1"
)

const Name = "mounts"

const (
	// SensitivePathsMounted occurs when a container has sensitive host paths mounted
	SensitivePathsMounted = "SensitivePathsMounted"
)

// DefaultSensitivePaths is the default list of sensitive mount paths (from Falco rule: https://github.com/falcosecurity/falco/blob/master/rules/k8s_audit_rules.yaml#L155)
var DefaultSensitivePaths = []string{"/proc", "/var/run/docker.sock", "/", "/etc", "/root", "/var/run/crio/crio.sock", "/home/admin", "/var/lib/kubelet", "/var/lib/kubelet/pki", "/etc/kubernetes", "/etc/kubernetes/manifests"}

const overrideLabelPrefix = "allow-host-path-mount-"

const (
	MountNameMetadataKey     = "MountName"
	MountPathMetadataKey     = "MountPath"
	MountReadOnlyMetadataKey = "MountReadOnly"
	MountVolumeNameKey       = "MountVolume"
	MountVolumeHostPathKey   = "MountVolumeHostPath"
)

// SensitivePathMounts implements Auditable
type SensitivePathMounts struct {
	sensitivePaths map[string]bool
}

func New(config Config) *SensitivePathMounts {
	paths := make(map[string]bool)
	for _, path := range config.GetSensitivePaths() {
		paths[path] = true
	}
	return &SensitivePathMounts{
		sensitivePaths: paths,
	}
}

// Audit checks that the container does not have any sensitive host path
func (sensitive *SensitivePathMounts) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	spec := k8s.GetPodSpec(resource)
	if spec == nil {
		return auditResults, nil
	}

	sensitiveVolumes := auditPodVolumes(spec, sensitive.sensitivePaths)

	if len(sensitiveVolumes) == 0 {
		return auditResults, nil
	}

	for _, container := range k8s.GetContainers(resource) {
		for _, auditResult := range auditContainer(container, sensitiveVolumes) {
			auditResult = override.ApplyOverride(auditResult, container.Name, resource, getOverrideLabel(auditResult.Metadata[MountNameMetadataKey]))
			if auditResult != nil {
				auditResults = append(auditResults, auditResult)
			}
		}
	}

	return auditResults, nil
}

func auditPodVolumes(podSpec *k8s.PodSpecV1, sensitivePaths map[string]bool) map[string]v1.Volume {
	if podSpec.Volumes == nil {
		return nil
	}

	found := make(map[string]v1.Volume)
	for _, volume := range podSpec.Volumes {
		if volume.HostPath == nil {
			continue
		}

		if _, ok := sensitivePaths[volume.HostPath.Path]; ok {
			found[volume.Name] = volume
		}
	}

	return found
}

func auditContainer(container *k8s.ContainerV1, sensitiveVolumes map[string]v1.Volume) []*kubeaudit.AuditResult {
	if container.VolumeMounts == nil {
		return nil
	}

	var auditResults []*kubeaudit.AuditResult

	for _, mount := range container.VolumeMounts {
		if volume, ok := sensitiveVolumes[mount.Name]; ok {
			auditResults = append(auditResults, &kubeaudit.AuditResult{
				Name:     SensitivePathsMounted,
				Severity: kubeaudit.Error,
				Message:  fmt.Sprintf("Sensitive path mounted as volume: %s (hostPath: %s). It should be removed from the container's mounts list.", mount.Name, volume.HostPath.Path),
				Metadata: kubeaudit.Metadata{
					"Container":              container.Name,
					MountNameMetadataKey:     mount.Name,
					MountPathMetadataKey:     mount.MountPath,
					MountReadOnlyMetadataKey: fmt.Sprintf("%t", mount.ReadOnly),
					MountVolumeNameKey:       volume.Name,
					MountVolumeHostPathKey:   volume.HostPath.Path,
				},
			})
		}
	}

	return auditResults
}

func getOverrideLabel(mountName string) string {
	return overrideLabelPrefix + mountName
}

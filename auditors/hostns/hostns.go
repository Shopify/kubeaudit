package hostns

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const Name = "hostns"

const (
	// NamespaceHostNetworkTrue occurs when hostNetwork is set to true in the container podspec
	NamespaceHostNetworkTrue = "NamespaceHostNetworkTrue"
	// NamespaceHostIPCTrue occurs when hostIPC is set to true in the container podspec
	NamespaceHostIPCTrue = "NamespaceHostIPCTrue"
	// NamespaceHostPIDTrue occurs when hostPID is set to true in the container podspec
	NamespaceHostPIDTrue = "NamespaceHostPIDTrue"
)

// HostNamespaces implements Auditable
type HostNamespaces struct{}

func New() *HostNamespaces {
	return &HostNamespaces{}
}

const HostNetworkOverrideLabel = "allow-namespace-host-network"
const HostIPCOverrideLabel = "allow-namespace-host-IPC"
const HostPIDOverrideLabel = "allow-namespace-host-PID"

// Audit checks that hostNetwork, hostIPC and hostPID are set to false in container podSpecs
func (a *HostNamespaces) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil, nil
	}

	for _, check := range []struct {
		auditFunc     func(*k8s.PodSpecV1) *kubeaudit.AuditResult
		overrideLabel string
	}{
		{auditHostNetwork, HostNetworkOverrideLabel},
		{auditHostIPC, HostIPCOverrideLabel},
		{auditHostPID, HostPIDOverrideLabel},
	} {
		auditResult := check.auditFunc(podSpec)
		auditResult = override.ApplyOverride(auditResult, "", resource, check.overrideLabel)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditHostNetwork(podSpec *k8s.PodSpecV1) *kubeaudit.AuditResult {
	if podSpec.HostNetwork {
		metadata := kubeaudit.Metadata{}
		if podSpec.Hostname != "" {
			metadata["PodHost"] = podSpec.Hostname
		}
		return &kubeaudit.AuditResult{
			Name:     NamespaceHostNetworkTrue,
			Severity: kubeaudit.Error,
			Message:  "hostNetwork is set to 'true' in PodSpec. It should be set to 'false'.",
			PendingFix: &fixHostNetworkTrue{
				podSpec: podSpec,
			},
			Metadata: metadata,
		}
	}

	return nil
}

func auditHostIPC(podSpec *k8s.PodSpecV1) *kubeaudit.AuditResult {
	if podSpec.HostIPC {
		metadata := kubeaudit.Metadata{}
		if podSpec.Hostname != "" {
			metadata["PodHost"] = podSpec.Hostname
		}
		return &kubeaudit.AuditResult{
			Name:     NamespaceHostIPCTrue,
			Severity: kubeaudit.Error,
			Message:  "hostIPC is set to 'true' in PodSpec. It should be set to 'false'.",
			PendingFix: &fixHostIPCTrue{
				podSpec: podSpec,
			},
			Metadata: metadata,
		}
	}

	return nil
}

func auditHostPID(podSpec *k8s.PodSpecV1) *kubeaudit.AuditResult {
	if podSpec.HostPID {
		metadata := kubeaudit.Metadata{}
		if podSpec.Hostname != "" {
			metadata["PodHost"] = podSpec.Hostname
		}
		return &kubeaudit.AuditResult{
			Name:     NamespaceHostPIDTrue,
			Severity: kubeaudit.Error,
			Message:  "hostPID is set to 'true' in PodSpec. It should be set to 'false'.",
			PendingFix: &fixHostPIDTrue{
				podSpec: podSpec,
			},
			Metadata: metadata,
		}
	}

	return nil
}

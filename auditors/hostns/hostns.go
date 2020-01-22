package hostns

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

// TODO this check used to only work on Pod resources, but now it works on any resource that has pods (including
// pods themselves). Is that ok or is there a reason it only worked on Pod before?

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
func (a *HostNamespaces) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, check := range []struct {
		auditFunc     func(k8stypes.Resource) *kubeaudit.AuditResult
		overrideLabel string
	}{
		{auditHostNetwork, HostNetworkOverrideLabel},
		{auditHostIPC, HostIPCOverrideLabel},
		{auditHostPID, HostPIDOverrideLabel},
	} {
		auditResult := check.auditFunc(resource)
		auditResult = override.ApplyOverride(auditResult, "", resource, check.overrideLabel)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditHostNetwork(resource k8stypes.Resource) *kubeaudit.AuditResult {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if podSpec.HostNetwork {
		return &kubeaudit.AuditResult{
			Name:     NamespaceHostNetworkTrue,
			Severity: kubeaudit.Error,
			Message:  "hostNetwork is set to 'true' in PodSpec. It should be set to 'false'.",
			PendingFix: &fixHostNetworkTrue{
				podSpec: podSpec,
			},
			Metadata: kubeaudit.Metadata{
				"PodHost": podSpec.Hostname,
			},
		}
	}

	return nil
}

func auditHostIPC(resource k8stypes.Resource) *kubeaudit.AuditResult {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if podSpec.HostIPC {
		return &kubeaudit.AuditResult{
			Name:     NamespaceHostIPCTrue,
			Severity: kubeaudit.Error,
			Message:  "hostIPC is set to 'true' in PodSpec. It should be set to 'false'.",
			PendingFix: &fixHostIPCTrue{
				podSpec: podSpec,
			},
			Metadata: kubeaudit.Metadata{
				"PodHost": podSpec.Hostname,
			},
		}
	}

	return nil
}

func auditHostPID(resource k8stypes.Resource) *kubeaudit.AuditResult {
	podSpec := k8s.GetPodSpec(resource)
	if podSpec == nil {
		return nil
	}

	if podSpec.HostPID {
		return &kubeaudit.AuditResult{
			Name:     NamespaceHostPIDTrue,
			Severity: kubeaudit.Error,
			Message:  "hostPID is set to 'true' in PodSpec. It should be set to 'false'.",
			PendingFix: &fixHostPIDTrue{
				podSpec: podSpec,
			},
			Metadata: kubeaudit.Metadata{
				"PodHost": podSpec.Hostname,
			},
		}
	}

	return nil
}

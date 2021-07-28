package override

import (
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
)

const (
	// ContainerOverrideLabelPrefix is used to disable an auditor for a specific container
	ContainerOverrideLabelPrefix = "container.audit.kubernetes.io/"
	// PodOverrideLabelPrefix is used to disable an auditor for a specific pod
	PodOverrideLabelPrefix = "audit.kubernetes.io/pod."
	// NamespaceOverrideLabelPrefix is used to disable an auditor for a specific namespace resource
	NamespaceOverrideLabelPrefix = "audit.kubernetes.io/namespace."
)

// GetOverriddenResultName takes an audit result name and modifies it to indicate that the security issue was
// ignored by an override label
func GetOverriddenResultName(resultName string) string {
	return resultName + "Allowed"
}

// NewRedundantOverrideResult creates a new AuditResult at warning level telling the user to remove the override
// label because there are no security issues found, so the label is redundant
func NewRedundantOverrideResult(containerName string, overrideReason, overrideLabel string) *kubeaudit.AuditResult {
	return &kubeaudit.AuditResult{
		Name:     kubeaudit.RedundantAuditorOverride,
		Severity: kubeaudit.Warn,
		Message:  "Auditor is disabled via label but there were no security issues found by the auditor. The label should be removed.",
		Metadata: kubeaudit.Metadata{
			"Container":     containerName,
			"OverrideLabel": overrideLabel,
		},
	}
}

// ApplyOverride checks if hasOverride is true. If it is, it changes the severity of the audit result from error to
// warn, adds the override reason to the metadata and removes the pending fix
func ApplyOverride(auditResult *kubeaudit.AuditResult, containerName string, resource k8s.Resource, overrideLabel string) *kubeaudit.AuditResult {
	hasOverride, overrideReason := GetContainerOverrideReason(containerName, resource, overrideLabel)

	if !hasOverride {
		return auditResult
	}

	if auditResult == nil {
		return NewRedundantOverrideResult(containerName, overrideReason, overrideLabel)
	}

	auditResult.Name = GetOverriddenResultName(auditResult.Name)
	auditResult.PendingFix = nil
	auditResult.Severity = kubeaudit.Info
	auditResult.Message = "Audit result overridden: " + auditResult.Message

	if overrideReason != "" && strings.ToLower(overrideReason) != "true" {
		if auditResult.Metadata == nil {
			auditResult.Metadata = make(kubeaudit.Metadata)
		}
		auditResult.Metadata["OverrideReason"] = overrideReason
	}

	return auditResult
}

// GetContainerOverrideReason returns true if the resource has a pod-level label disabling a given auditor and the
// value of the label which is meant to represent the reason for overriding the auditor
//
// Container override labels disable the auditor for that specific container and have the following format:
// 		container.audit.kubernetes.io/[container name].[auditor override label]
//
// If there is no container override label, it calls GetResourceOverrideReason()
func GetContainerOverrideReason(containerName string, resource k8s.Resource, overrideLabel string) (hasOverride bool, reason string) {
	labels := k8s.GetLabels(resource)

	if containerName != "" {
		if reason, hasOverride = labels[GetContainerOverrideLabel(containerName, overrideLabel)]; hasOverride {
			return
		}
	}

	return GetResourceOverrideReason(resource, overrideLabel)
}

// GetResourceOverrideReason returns true if the resource has a label disabling a given auditor and the value of the
// label which is meant to represent the reason for overriding the auditor
//
// Pod override labels disable the auditor for the pod and all containers within the pod and have the following format:
// 		audit.kubernetes.io/pod.[auditor override label]
// Namespace override labels disable the auditor for the namespace resource and have the following format:
// 		audit.kubernetes.io/namespace.[auditor override label]
func GetResourceOverrideReason(resource k8s.Resource, auditorOverrideLabel string) (hasOverride bool, reason string) {
	labelFuncs := []func(overrideLabel string) string{
		GetPodOverrideLabel,
		GetNamespaceOverrideLabel,
	}

	labels := k8s.GetLabels(resource)
	for _, getLabel := range labelFuncs {
		if reason, hasOverride = labels[getLabel(auditorOverrideLabel)]; hasOverride {
			return
		}
	}

	return false, ""
}

func GetPodOverrideLabel(overrideLabel string) string {
	return PodOverrideLabelPrefix + overrideLabel
}

func GetNamespaceOverrideLabel(overrideLabel string) string {
	return NamespaceOverrideLabelPrefix + overrideLabel
}

func GetContainerOverrideLabel(containerName, overrideLabel string) string {
	return ContainerOverrideLabelPrefix + containerName + "." + overrideLabel
}

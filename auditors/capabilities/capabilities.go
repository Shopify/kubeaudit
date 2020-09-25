package capabilities

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const Name = "capabilities"

const (
	// CapabilityAdded occurs when a capability is in the capability add list of a container's security context
	CapabilityAdded = "CapabilityAdded"
	// CapabilityShouldDropAll occurs when there's a drop list instead of having drop "ALL"
	CapabilityShouldDropAll = "CapabilityShouldDropAll"
)

const overrideLabelPrefix = "allow-capability-"

var DefaultDropList = []string{"ALL"}

// Capabilities implements Auditable
type Capabilities struct {
	dropList []string
}

func New(config Config) *Capabilities {
	return &Capabilities{
		dropList: config.GetDropList(),
	}
}

// Audit checks that bad capabilities are dropped with ALL and no capabilities are added
func (a *Capabilities) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		for _, capability := range mergeCapabilities(container) {
			for _, auditResult := range auditContainer(container, capability) {
				auditResult = override.ApplyOverride(auditResult, container.Name, resource, getOverrideLabel(capability))
				if auditResult != nil {
					auditResults = append(auditResults, auditResult)
				}
			}
		}
	}

	return auditResults, nil
}

func getOverrideLabel(capability string) string {
	return overrideLabelPrefix + strings.Replace(strings.ToLower(capability), "_", "-", -1)
}

func auditContainer(container *k8stypes.ContainerV1, capability string) []*kubeaudit.AuditResult {
	var auditResults []*kubeaudit.AuditResult

	if SecurityContextOrCapabilities(container) {
		var message string

		if IsCapabilityInAddList(container, capability) {
			message = fmt.Sprintf("Capability added. It should be removed from the capability add list. If you need this capability, add an override label such as '%s: SomeReason'.", override.GetContainerOverrideLabel(container.Name, getOverrideLabel(capability)))
			auditResult := &kubeaudit.AuditResult{
				Name:     CapabilityAdded,
				Severity: kubeaudit.Error,
				Message:  message,
				PendingFix: &fixCapabilityAdded{
					container:  container,
					capability: capability,
				},
				Metadata: kubeaudit.Metadata{
					"Container": container.Name,
				},
			}
			auditResults = append(auditResults, auditResult)
		}

		if !IsDropAll(container) && !IsCapabilityInAddList(container, capability) {
			message = "Capabily not set to ALL. Ideally, you should drop ALL capabilities and add the specific ones you need to the add list."
			auditResult := &kubeaudit.AuditResult{
				Name:     CapabilityShouldDropAll,
				Severity: kubeaudit.Error,
				Message:  message,
				PendingFix: &fixCapabilityNotDroppedAll{
					container:  container,
					capability: capability,
				},
				Metadata: kubeaudit.Metadata{
					"Container": container.Name,
				},
			}
			auditResults = append(auditResults, auditResult)
		}
	} else {
		message := "Security Context not set. Ideally, the Security Context should be specified. All capacities should be dropped by setting drop to ALL."
		auditResult := &kubeaudit.AuditResult{
			Name:     CapabilityShouldDropAll,
			Severity: kubeaudit.Error,
			Message:  message,
			PendingFix: &fixMissingSecurityContextOrCapability{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
		auditResults = append(auditResults, auditResult)
	}
	// We need the audit result to be nil for ApplyOverride to check for RedundantAuditorOverride errors
	if len(auditResults) == 0 {
		return []*kubeaudit.AuditResult{nil}
	}

	return auditResults
}

func IsDropAll(container *k8stypes.ContainerV1) bool {
	for _, cap := range container.SecurityContext.Capabilities.Drop {
		if strings.ToUpper(string(cap)) == "ALL" {
			return true
		}
	}

	return false
}

func IsCapabilityInAddList(container *k8stypes.ContainerV1, capability string) bool {
	for _, cap := range container.SecurityContext.Capabilities.Add {
		if string(cap) == capability {
			return true
		}
	}

	return false
}

func SecurityContextOrCapabilities(container *k8stypes.ContainerV1) bool {
	if container.SecurityContext == nil || container.SecurityContext.Capabilities == nil {
		return false
	}

	return true
}

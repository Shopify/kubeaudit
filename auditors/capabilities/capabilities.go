package capabilities

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const Name = "capabilities"

const (
	// CapabilityAdded occurs when a capability is in the capability add list of a container's security context
	CapabilityAdded = "CapabilityAdded"
	// CapabilityShouldDropAll occurs when there's a drop list instead of having drop "ALL"
	CapabilityShouldDropAll = "CapabilityShouldDropAll"
	// CapabilityOrSecurityContextMissing  occurs when either the Security Context or Capabilities are not specified
	CapabilityOrSecurityContextMissing = "CapabilityOrSecurityContextMissing"
)

const overrideLabelPrefix = "allow-capability-"

var DefaultDropList = []string{"ALL"}

var DefaultAllowAddList = []string{""}

// Capabilities implements Auditable
type Capabilities struct {
	allowAddList []string
}

func New(config Config) *Capabilities {
	return &Capabilities{
		allowAddList: config.GetAllowAddList(),
	}
}

// Audit checks that bad capabilities are dropped with ALL and no capabilities are added
func (a *Capabilities) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainerForDropAll(container)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}

		for _, capability := range uniqueCapabilities(container) {
			for _, auditResult := range auditContainer(container, capability, a.allowAddList) {
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

func auditContainer(container *k8s.ContainerV1, capability string, allowAddList []string) []*kubeaudit.AuditResult {
	var auditResults []*kubeaudit.AuditResult

	if isCapabilityInArray(capability, allowAddList) {
		return auditResults
	}

	if SecurityContextOrCapabilities(container) {
		if IsCapabilityInAddList(container, capability) {
			message := fmt.Sprintf("Capability \"%s\" added. It should be removed from the capability add list. If you need this capability, add an override label such as '%s: SomeReason'.", capability, override.GetContainerOverrideLabel(container.Name, getOverrideLabel(capability)))
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
					"Metadata":  capability,
				},
			}
			auditResults = append(auditResults, auditResult)
		}
	}
	// We need the audit result to be nil for ApplyOverride to check for RedundantAuditorOverride errors

	if len(auditResults) == 0 {
		return []*kubeaudit.AuditResult{nil}
	}

	return auditResults
}

func auditContainerForDropAll(container *k8s.ContainerV1) *kubeaudit.AuditResult {
	if !SecurityContextOrCapabilities(container) {
		message := "Security Context not set. The Security Context should be specified and all Capabilities should be dropped by setting the Drop list to ALL."
		return &kubeaudit.AuditResult{
			Name:     CapabilityOrSecurityContextMissing,
			Severity: kubeaudit.Error,
			Message:  message,
			PendingFix: &fixMissingSecurityContextOrCapability{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	if !IsDropAll(container) {
		message := "Capability Drop list should be set to ALL. Add the specific ones you need to the Add list and set an override label."
		return &kubeaudit.AuditResult{
			Name:     CapabilityShouldDropAll,
			Severity: kubeaudit.Error,
			Message:  message,
			PendingFix: &fixCapabilityNotDroppedAll{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}
	return nil
}

package rootfs

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const Name = "rootfs"

const (
	// ReadOnlyRootFilesystemFalse occurs when readOnlyRootFilesystem is set to false in the container SecurityContext
	ReadOnlyRootFilesystemFalse = "ReadOnlyRootFilesystemFalse"
	// ReadOnlyRootFilesystemNil occurs when readOnlyRootFilesystem is not set in the container SecurityContext.
	// readOnlyRootFilesystem defaults to false so this is bad
	ReadOnlyRootFilesystemNil = "ReadOnlyRootFilesystemNil"
)

const OverrideLabel = "allow-read-only-root-filesystem-false"

// ReadOnlyRootFilesystem implements Auditable
type ReadOnlyRootFilesystem struct{}

func New() *ReadOnlyRootFilesystem {
	return &ReadOnlyRootFilesystem{}
}

// Audit checks that readOnlyRootFilesystem is set to true in every container's security context
func (a *ReadOnlyRootFilesystem) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container, resource)
		auditResult = override.ApplyOverride(auditResult, container.Name, resource, OverrideLabel)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditContainer(container *k8stypes.ContainerV1, resource k8stypes.Resource) *kubeaudit.AuditResult {
	if isReadOnlyRootFilesystemNil(container) {
		return &kubeaudit.AuditResult{
			Name:     ReadOnlyRootFilesystemNil,
			Severity: kubeaudit.Error,
			Message:  "readOnlyRootFilesystem is not set in container SecurityContext. It should be set to 'true'.",
			PendingFix: &fixReadOnlyRootFilesystem{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	if isReadOnlyRootFilesystemFalse(container) {
		return &kubeaudit.AuditResult{
			Name:     ReadOnlyRootFilesystemFalse,
			Severity: kubeaudit.Error,
			Message:  "readOnlyRootFilesystem is set to 'false' in container SecurityContext. It should be set to 'true'.",
			PendingFix: &fixReadOnlyRootFilesystem{
				container: container,
			},
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	return nil
}

func isReadOnlyRootFilesystemFalse(container *k8stypes.ContainerV1) bool {
	if isReadOnlyRootFilesystemNil(container) {
		return true
	}

	return !*container.SecurityContext.ReadOnlyRootFilesystem
}

func isReadOnlyRootFilesystemNil(container *k8stypes.ContainerV1) bool {
	return container.SecurityContext == nil || container.SecurityContext.ReadOnlyRootFilesystem == nil
}

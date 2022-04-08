package deprecatedapis

import (
	"fmt"
	"strconv"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const Name = "deprecatedapis"

const (
	// DeprecatedAPIUsed occurs when a deprecated resource type version is used
	DeprecatedAPIUsed = "DeprecatedAPIUsed"
)

// DeprecatedAPIs implements Auditable
type DeprecatedAPIs struct {
	CurrentVersion  *Version
	TargetedVersion *Version
}

func New(config Config) (*DeprecatedAPIs, error) {
	currentVersion, err := config.GetCurrentVersion()
	if err != nil {
		return nil, fmt.Errorf("error creating DeprecatedAPIs auditor: %w", err)
	}

	targetedVersion, err := config.GetTargetedVersion()
	if err != nil {
		return nil, fmt.Errorf("error creating DeprecatedAPIs auditor: %w", err)
	}

	return &DeprecatedAPIs{
		CurrentVersion:  currentVersion,
		TargetedVersion: targetedVersion,
	}, nil
}

type apiLifecycleDeprecated interface {
	APILifecycleDeprecated() (major, minor int)
}

type apiLifecycleRemoved interface {
	APILifecycleRemoved() (major, minor int)
}

type apiLifecycleReplacement interface {
	APILifecycleReplacement() schema.GroupVersionKind
}

type apiLifecycleIntroduced interface {
	APILifecycleIntroduced() (major, minor int)
}

// Audit checks that the resource API version is not deprecetated
func (deprecatedAPIs *DeprecatedAPIs) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult
	lastApplied, ok := k8s.GetAnnotations(resource)[v1.LastAppliedConfigAnnotation]
	if ok && len(lastApplied) > 0 {
		resource, _ = k8sinternal.DecodeResource([]byte(lastApplied))
	}
	deprecated, isDeprecated := resource.(apiLifecycleDeprecated)
	if isDeprecated {
		deprecatedMajor, deprecatedMinor := deprecated.APILifecycleDeprecated()
		if deprecatedMajor == 0 && deprecatedMinor == 0 {
			return nil, fmt.Errorf("Version not found %s (%d.%d)", deprecated, deprecatedMajor, deprecatedMinor)
		} else {
			severity := kubeaudit.Warn
			metadata := kubeaudit.Metadata{
				"DeprecatedMajor": strconv.Itoa(deprecatedMajor),
				"DeprecatedMinor": strconv.Itoa(deprecatedMinor),
			}
			if deprecatedAPIs.CurrentVersion != nil && (deprecatedAPIs.CurrentVersion.Major < deprecatedMajor || deprecatedAPIs.CurrentVersion.Major == deprecatedMajor && deprecatedAPIs.CurrentVersion.Minor < deprecatedMinor) {
				severity = kubeaudit.Info
			}
			gvk := resource.GetObjectKind().GroupVersionKind()
			if gvk.Empty() {
				return nil, fmt.Errorf("GroupVersionKind not found %s", resource)
			} else {
				deprecationMessage := fmt.Sprintf("%s %s is deprecated in v%d.%d+", gvk.GroupVersion().String(), gvk.Kind, deprecatedMajor, deprecatedMinor)
				if removed, hasRemovalInfo := resource.(apiLifecycleRemoved); hasRemovalInfo {
					removedMajor, removedMinor := removed.APILifecycleRemoved()
					if removedMajor != 0 || removedMinor != 0 {
						deprecationMessage = deprecationMessage + fmt.Sprintf(", unavailable in v%d.%d+", removedMajor, removedMinor)
						metadata["RemovedMajor"] = strconv.Itoa(removedMajor)
						metadata["RemovedMinor"] = strconv.Itoa(removedMinor)
					}
					if deprecatedAPIs.TargetedVersion != nil && deprecatedAPIs.TargetedVersion.Major >= removedMajor && deprecatedAPIs.TargetedVersion.Minor >= removedMinor {
						severity = kubeaudit.Error
					}
				}

				if replaced, hasReplacement := resource.(apiLifecycleReplacement); hasReplacement {
					replacement := replaced.APILifecycleReplacement()
					if !replacement.Empty() {
						deprecationMessage = deprecationMessage + fmt.Sprintf("; use %s %s", replacement.GroupVersion().String(), replacement.Kind)
						metadata["ReplacementGroup"] = replacement.GroupVersion().String()
						metadata["ReplacementKind"] = replacement.Kind
					}
				}
				auditResult := &kubeaudit.AuditResult{
					Name:     DeprecatedAPIUsed,
					Severity: severity,
					Message:  deprecationMessage,
					Metadata: metadata,
				}
				auditResults = append(auditResults, auditResult)
			}
		}

	}
	return auditResults, nil
}

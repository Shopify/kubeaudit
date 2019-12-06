package apparmor

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/fix"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const (
	// AppArmorAnnotationMissing occurs when the apparmor annotation is missing
	AppArmorAnnotationMissing = "AppArmorAnnotationMissing"
	// AppArmorDisabled occurs when the apparmor annotation is set to a bad value
	AppArmorDisabled = "AppArmorDisabled"
)

// As of Oct 1, 2018 these constants are not in the K8s API package, but once they are they should be replaced
// https://github.com/kubernetes/kubernetes/blob/7f23a743e8c23ac6489340bbb34fa6f1d392db9d/pkg/security/apparmor/helpers.go#L25
const (
	// The prefix to an annotation key specifying a container profile.
	ContainerAnnotationKeyPrefix = "container.apparmor.security.beta.kubernetes.io/"

	// The profile specifying the runtime default.
	ProfileRuntimeDefault = "runtime/default"
	// The prefix for specifying profiles loaded on the node.
	ProfileNamePrefix = "localhost/"
)

// AppArmor implements Auditable
type AppArmor struct{}

func New() *AppArmor {
	return &AppArmor{}
}

// Audit checks that AppArmor is enabled for all containers
func (a *AppArmor) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container, resource)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditContainer(container *k8stypes.ContainerV1, resource k8stypes.Resource) *kubeaudit.AuditResult {
	annotations := k8s.GetAnnotations(resource)
	containerAnnotation := getContainerAnnotation(container)

	if isAppArmorAnnotationMissing(containerAnnotation, annotations) {
		return &kubeaudit.AuditResult{
			Name:     AppArmorAnnotationMissing,
			Severity: kubeaudit.Error,
			Message:  fmt.Sprintf("AppArmor annotation missing. The annotation '%s' should be added.", containerAnnotation),
			Metadata: kubeaudit.Metadata{
				"Container":         container.Name,
				"MissingAnnotation": containerAnnotation,
			},
			PendingFix: &fix.ByAddingPodAnnotation{
				Key:   containerAnnotation,
				Value: ProfileRuntimeDefault,
			},
		}
	}

	if isAppArmorDisabled(containerAnnotation, annotations) {
		return &kubeaudit.AuditResult{
			Name:     AppArmorDisabled,
			Message:  fmt.Sprintf("AppArmor is disabled. The apparmor annotation should be set to '%s' or start with '%s'.", ProfileRuntimeDefault, ProfileNamePrefix),
			Severity: kubeaudit.Error,
			Metadata: kubeaudit.Metadata{
				"Container":       container.Name,
				"Annotation":      containerAnnotation,
				"AnnotationValue": getProfileName(containerAnnotation, annotations),
			},
			PendingFix: &fix.BySettingPodAnnotation{
				Key:   containerAnnotation,
				Value: ProfileRuntimeDefault,
			},
		}
	}

	return nil
}

func isAppArmorAnnotationMissing(apparmorAnnotation string, annotations map[string]string) bool {
	_, ok := annotations[apparmorAnnotation]
	return !ok
}

func isAppArmorDisabled(apparmorAnnotation string, annotations map[string]string) bool {
	profileName, ok := annotations[apparmorAnnotation]
	return !ok || profileName != ProfileRuntimeDefault && !strings.HasPrefix(profileName, ProfileNamePrefix)
}

func getContainerAnnotation(container *k8stypes.ContainerV1) string {
	return ContainerAnnotationKeyPrefix + container.Name
}

func getProfileName(apparmorAnnotation string, annotations map[string]string) string {
	profileName, _ := annotations[apparmorAnnotation]
	return profileName
}

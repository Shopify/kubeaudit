package all

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/asat"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/hostns"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"github.com/Shopify/kubeaudit/auditors/netpols"
	"github.com/Shopify/kubeaudit/auditors/nonroot"
	"github.com/Shopify/kubeaudit/auditors/privesc"
	"github.com/Shopify/kubeaudit/auditors/privileged"
	"github.com/Shopify/kubeaudit/auditors/rootfs"
	"github.com/Shopify/kubeaudit/auditors/seccomp"
	"github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fixtureDir = "../../internal/test/fixtures/all_resources"

func TestAuditAll(t *testing.T) {
	allErrors := []string{
		apparmor.AppArmorAnnotationMissing,
		asat.AutomountServiceAccountTokenTrueAndDefaultSA,
		capabilities.CapabilityNotDropped,
		hostns.NamespaceHostNetworkTrue,
		hostns.NamespaceHostIPCTrue,
		hostns.NamespaceHostPIDTrue,
		image.ImageTagMissing,
		limits.LimitsNotSet,
		netpols.MissingDefaultDenyIngressAndEgressNetworkPolicy,
		nonroot.RunAsNonRootPSCNilCSCNil,
		privesc.AllowPrivilegeEscalationNil,
		privileged.PrivilegedNil,
		rootfs.ReadOnlyRootFilesystemNil,
		seccomp.SeccompAnnotationMissing,
	}

	allAuditors, err := Auditors(config.KubeauditConfig{})
	require.NoError(t, err)

	for _, file := range test.GetAllFileNames(t, fixtureDir) {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		file := file
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			test.AuditMultiple(t, fixtureDir, file, allAuditors, allErrors, "", test.MANIFEST_MODE)
			test.AuditMultiple(t, fixtureDir, file, allAuditors, allErrors, strings.Split(file, ".")[0], test.LOCAL_MODE)
		})
	}
}

func TestFixAll(t *testing.T) {
	allAuditors, err := Auditors(config.KubeauditConfig{})
	require.NoError(t, err)

	for _, file := range test.GetAllFileNames(t, fixtureDir) {
		t.Run(file, func(t *testing.T) {
			_, report := test.FixSetupMultiple(t, fixtureDir, file, allAuditors)
			for _, result := range report.Results() {
				for _, auditResult := range result.GetAuditResults() {
					if !assert.NotEqual(t, kubeaudit.Error, auditResult.Severity) {
						return
					}
				}
			}
		})
	}
}

// Test all auditors with config
func TestAllWithConfig(t *testing.T) {
	enabledAuditors := []string{
		apparmor.Name, seccomp.Name,
	}
	expectedErrors := []string{
		apparmor.AppArmorAnnotationMissing,
		seccomp.SeccompAnnotationMissing,
	}

	conf := config.KubeauditConfig{
		EnabledAuditors: enabledAuditorsToMap(enabledAuditors),
	}
	auditors, err := Auditors(conf)
	require.NoError(t, err)

	for _, file := range test.GetAllFileNames(t, fixtureDir) {
		t.Run(file, func(t *testing.T) {
			test.AuditMultiple(t, fixtureDir, file, auditors, expectedErrors, "", test.MANIFEST_MODE)
		})
	}
}

func enabledAuditorsToMap(enabledAuditors []string) map[string]bool {
	enabledAuditorMap := map[string]bool{}
	for _, auditorName := range AuditorNames {
		enabledAuditorMap[auditorName] = false
	}
	for _, auditorName := range enabledAuditors {
		enabledAuditorMap[auditorName] = true
	}
	return enabledAuditorMap
}

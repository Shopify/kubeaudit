package all

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/asat"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"github.com/Shopify/kubeaudit/auditors/nonroot"
	"github.com/Shopify/kubeaudit/auditors/privesc"
	"github.com/Shopify/kubeaudit/auditors/privileged"
	"github.com/Shopify/kubeaudit/auditors/rootfs"
	"github.com/Shopify/kubeaudit/auditors/seccomp"
	"github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "fixtures"

func TestAuditAll(t *testing.T) {
	files := []string{"audit_all_v1.yml", "audit_all_v1beta1.yml"}
	allErrors := []string{
		privesc.AllowPrivilegeEscalationNil,
		asat.AutomountServiceAccountTokenTrueAndDefaultSA,
		capabilities.CapabilityNotDropped,
		image.ImageTagMissing,
		privileged.PrivilegedNil,
		rootfs.ReadOnlyRootFilesystemNil,
		limits.LimitsNotSet,
		nonroot.RunAsNonRootPSCNilCSCNil,
		apparmor.AppArmorAnnotationMissing,
		seccomp.SeccompAnnotationMissing,
	}

	allAuditors, err := Auditors(config.KubeauditConfig{})
	assert.NoError(t, err)

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			test.AuditMultiple(t, fixtureDir, file, allAuditors, allErrors)
		})
	}
}

// Test that fixing all fixtures in auditors/* results in manifests that pass all audits
func TestAllForRegression(t *testing.T) {
	auditorDirs, err := ioutil.ReadDir("../")
	if !assert.Nil(t, err) {
		return
	}

	allAuditors, err := Auditors(config.KubeauditConfig{})
	assert.NoError(t, err)

	for _, auditorDir := range auditorDirs {
		if !auditorDir.IsDir() {
			continue
		}

		fixturesDirPath := filepath.Join("..", auditorDir.Name(), "fixtures")
		fixtureFiles, err := ioutil.ReadDir(fixturesDirPath)
		if os.IsNotExist(err) {
			continue
		}
		if !assert.Nil(t, err) {
			return
		}

		for _, fixture := range fixtureFiles {
			t.Run(filepath.Join(fixturesDirPath, fixture.Name()), func(t *testing.T) {
				_, report := test.FixSetupMultiple(t, fixturesDirPath, fixture.Name(), allAuditors)
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
}

// Test all auditors with config
func TestAllWithConfig(t *testing.T) {
	files := []string{"audit_all_v1.yml", "audit_all_v1beta1.yml"}
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
	assert.NoError(t, err)

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			test.AuditMultiple(t, fixtureDir, file, auditors, expectedErrors)
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

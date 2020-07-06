package kubeaudit_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that fixing all fixtures in auditors/* results in manifests that pass all audits
func TestFix(t *testing.T) {
	auditorDirs, err := ioutil.ReadDir("auditors")
	if !assert.Nil(t, err) {
		return
	}

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	require.NoError(t, err)

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

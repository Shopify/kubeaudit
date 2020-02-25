package kubeaudit_test

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	assert.NoError(err)

	auditor, err := kubeaudit.New(allAuditors)
	assert.Nil(err)
	assert.NotNil(auditor)

	_, err = kubeaudit.New(nil)
	assert.NotNil(err)

}

func TestAuditLocal(t *testing.T) {
	assert := assert.New(t)

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	assert.NoError(err)

	auditor, err := kubeaudit.New(allAuditors)
	assert.Nil(err)

	_, err = auditor.AuditLocal("path")
	assert.NotNil(err)
}

func TestAuditCluster(t *testing.T) {
	assert := assert.New(t)

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	assert.NoError(err)

	auditor, err := kubeaudit.New(allAuditors)
	assert.Nil(err)

	_, err = auditor.AuditCluster("")
	assert.NotNil(err)
}

func TestUnknownResource(t *testing.T) {
	// Make sure we produce only warning results for resources kubeaudit doesn't know how to audit
	files := []string{"unknown_type_v1.yml", "custom_resource_definition.yml"}

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	assert.NoError(t, err)

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			_, report := test.FixSetupMultiple(t, "internal/test/fixtures", file, allAuditors)
			require.NotNil(t, report)
			for _, result := range report.Results() {
				for _, auditResult := range result.GetAuditResults() {
					assert.Equal(t, kubeaudit.Warn, auditResult.Severity)
				}
			}
		})
	}
}

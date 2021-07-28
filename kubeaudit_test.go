package kubeaudit_test

import (
	"os"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require := require.New(t)

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	require.NoError(err)

	auditor, err := kubeaudit.New(allAuditors)
	require.NoError(err)
	assert.NotNil(t, auditor)

	_, err = kubeaudit.New(nil)
	require.NotNil(err)

}

func TestAuditLocal(t *testing.T) {
	if os.Getenv("USE_KIND") == "false" {
		return
	}

	require := require.New(t)

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	require.NoError(err)

	auditor, err := kubeaudit.New(allAuditors)
	require.NoError(err)

	_, err = auditor.AuditLocal("", k8sinternal.ClientOptions{})
	require.NoError(err)

	_, err = auditor.AuditLocal("invalid_path", k8sinternal.ClientOptions{})
	require.NotNil(err)
}

func TestAuditCluster(t *testing.T) {
	require := require.New(t)

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	require.NoError(err)

	auditor, err := kubeaudit.New(allAuditors)
	require.NoError(err)

	_, err = auditor.AuditCluster(k8sinternal.ClientOptions{})
	require.NotNil(err)
}

func TestUnknownResource(t *testing.T) {
	// Make sure we produce only warning results for resources kubeaudit doesn't know how to audit
	files := []string{"unknown_resource_type.yml", "custom_resource_definition.yml"}

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	require.NoError(t, err)

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

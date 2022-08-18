package deprecatedapis

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fixtureDir = "fixtures"

func TestAuditDeprecatedAPIs(t *testing.T) {
	cases := []struct {
		file             string
		currentVersion   string
		targetedVersion  string
		expectedSeverity kubeaudit.SeverityLevel
	}{
		{"cronjob.yml", "", "", kubeaudit.Warn},          // Warn is the serverity by default
		{"cronjob.yml", "1.20", "1.21", kubeaudit.Info},  // Info, not yet deprecated in the current version
		{"cronjob.yml", "1.21", "1.22", kubeaudit.Warn},  // Warn, deprecated in the current version
		{"cronjob.yml", "1.22", "1.25", kubeaudit.Error}, // Error, not available in the targeted version
		{"cronjob.yml", "1.20", "1.25", kubeaudit.Error}, // Error, not yet deprecated in the current version but not available in the targeted version
		{"cronjob.yml", "1.20", "", kubeaudit.Info},      // Info, not yet deprecated in the current version and no targeted version defined
		{"cronjob.yml", "1.21", "", kubeaudit.Warn},      // Warn, deprecated in the current version
		{"cronjob.yml", "", "1.20", kubeaudit.Warn},      // Warn is the serverity by default if no current version
		{"cronjob.yml", "", "1.25", kubeaudit.Error},     // Error, not available in the targeted version
	}

	message := "batch/v1beta1 CronJob is deprecated in v1.21+, unavailable in v1.25+, introduced in v1.8+; use batch/v1 CronJob"
	metadata := kubeaudit.Metadata{
		"DeprecatedMajor":  "1",
		"DeprecatedMinor":  "21",
		"RemovedMajor":     "1",
		"RemovedMinor":     "25",
		"IntroducedMajor":  "1",
		"IntroducedMinor":  "8",
		"ReplacementGroup": "batch/v1",
		"ReplacementKind":  "CronJob",
	}

	for i, tc := range cases {
		// These lines are needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		i := i
		t.Run(tc.file+"-"+tc.currentVersion+"-"+tc.targetedVersion, func(t *testing.T) {
			t.Parallel()
			auditor, err := New(Config{CurrentVersion: tc.currentVersion, TargetedVersion: tc.targetedVersion})
			require.Nil(t, err)
			report := test.AuditManifest(t, fixtureDir, tc.file, auditor, []string{DeprecatedAPIUsed})
			assertReport(t, report, tc.expectedSeverity, message, metadata)

			report = test.AuditLocal(t, fixtureDir, tc.file, auditor, fmt.Sprintf("%s-%d", strings.Split(tc.file, ".")[0], i), []string{DeprecatedAPIUsed})

			if report != nil {
				assertReport(t, report, tc.expectedSeverity, message, metadata)
			}
		})
	}
}

func assertReport(t *testing.T, report *kubeaudit.Report, expectedSeverity kubeaudit.SeverityLevel, message string, metadata map[string]string) {
	assert.Equal(t, 1, len(report.Results()))
	for _, result := range report.Results() {
		assert.Equal(t, 1, len(result.GetAuditResults()))
		for _, auditResult := range result.GetAuditResults() {
			require.Equal(t, expectedSeverity, auditResult.Severity)
			require.Equal(t, message, auditResult.Message)
			require.Equal(t, metadata, auditResult.Metadata)
		}
	}
}

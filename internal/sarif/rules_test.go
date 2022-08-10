package sarif

import (
	"testing"

	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/stretchr/testify/assert"
)

func TestAuditorsLengthAndDescription(t *testing.T) {
	// if new auditors are created
	// make sure they're added with a matching description
	for _, auditorName := range all.AuditorNames {
		description, ok := allAuditors[auditorName]
		assert.Truef(t, ok && description != "", "missing description for auditor %s", auditorName)
	}
}

func TestViolationToRules(t *testing.T) {
	// if new rules are added to any of the auditors
	// they should be captured in the mapping
	cases := []struct {
		auditorName   string
		expectedCount int
	}{
		{
			"apparmor",
			3,
		},
		{
			"asat",
			2,
		},
		{
			"capabilities",
			3,
		},
		{
			"deprecatedapis",
			1,
		},
		{
			"hostns",
			3,
		},
		{
			"image",
			3,
		},
		{
			"limits",
			5,
		},
		{
			"mounts",
			1,
		},
		{
			"netpols",
			5,
		},
		{
			"nonroot",
			5,
		},
		{
			"privesc",
			2,
		},
		{
			"privileged",
			2,
		},
		{
			"rootfs",
			2,
		},
		{
			"seccomp",
			5,
		},
	}

	assert.Len(t, cases, len(all.AuditorNames))

	for _, c := range cases {
		var totalCount int

		for _, v := range violationsToRules {
			if v == c.auditorName {
				totalCount += 1
			}
		}

		assert.Equal(t, c.expectedCount, totalCount)
	}
}

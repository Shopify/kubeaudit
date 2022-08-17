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

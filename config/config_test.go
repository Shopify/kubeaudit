package config_test

import (
	"os"
	"testing"

	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that the sample config includes all auditors
func TestConfig(t *testing.T) {
	configFile := "config.yaml"
	reader, err := os.Open(configFile)
	require.NoError(t, err)

	conf, err := config.New(reader)
	require.NoError(t, err)

	assert.Equal(t, len(all.AuditorNames), len(conf.GetEnabledAuditors()), "Config is missing auditors")
}

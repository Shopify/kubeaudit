package kubeaudit_test

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithLogger(t *testing.T) {
	assert := assert.New(t)

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	require.NoError(t, err)

	formatter := logrus.Formatter(&logrus.JSONFormatter{})
	_, err = kubeaudit.New(allAuditors, kubeaudit.WithLogger(formatter))
	assert.NoError(err)
	assert.Equal(formatter, logrus.StandardLogger().Formatter)

	formatter = logrus.Formatter(&logrus.TextFormatter{})
	assert.NotEqual(formatter, logrus.StandardLogger().Formatter)

	_, err = kubeaudit.New(allAuditors, kubeaudit.WithLogger(formatter))
	assert.NoError(err)
	assert.Equal(formatter, logrus.StandardLogger().Formatter)
}

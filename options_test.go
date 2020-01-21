package kubeaudit_test

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWithLogger(t *testing.T) {
	assert := assert.New(t)

	formatter := logrus.Formatter(&logrus.JSONFormatter{})
	_, err := kubeaudit.New(all.Auditors(), kubeaudit.WithLogger(formatter))
	assert.Nil(err)
	assert.Equal(formatter, logrus.StandardLogger().Formatter)

	formatter = logrus.Formatter(&logrus.TextFormatter{})
	assert.NotEqual(formatter, logrus.StandardLogger().Formatter)

	_, err = kubeaudit.New(all.Auditors(), kubeaudit.WithLogger(formatter))
	assert.Nil(err)
	assert.Equal(formatter, logrus.StandardLogger().Formatter)
}

package cmd

import (
	"testing"

	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestShouldLog(t *testing.T) {
	errorCodes := []int{ErrorServiceAccountTokenDeprecated, InfoImageCorrect, ErrorImageTagIncorrect,
		ErrorImageTagMissing, ErrorResourcesLimitsCPUExceeded, ErrorResourcesLimitsMemoryExceeded}
	for _, err := range errorCodes {
		members := shouldLog(err)
		switch err {
		case ErrorServiceAccountTokenDeprecated:
			assert.Equal(t, len(members), 5)
		case InfoImageCorrect:
			assert.Equal(t, len(members), 3)
		case ErrorImageTagMissing:
			assert.Equal(t, len(members), 3)
		case ErrorImageTagIncorrect:
			assert.Equal(t, len(members), 5)
		case ErrorResourcesLimitsCPUExceeded:
			assert.Equal(t, len(members), 5)
		case ErrorResourcesLimitsMemoryExceeded:
			assert.Equal(t, len(members), 5)
		}
	}
}

func TestCreateFields(t *testing.T) {
	rootConfig.manifest = "../fixtures/run_as_non_root_psc_false_csc_nil_multiple_cont_v1.yml"
	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	assert.Nil(t, err)
	results := getResults(resources, auditRunAsNonRoot)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 1, len(results[0].Occurrences))
	fields := createFields(results[0], results[0].Occurrences[0])
	assert.Equal(t, 5, len(fields))
}

func TestPrint(t *testing.T) {
	rootConfig.manifest = "../fixtures/run_as_non_root_psc_false_csc_nil_multiple_cont_v1.yml"
	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	assert.Nil(t, err)
	results := getResults(resources, auditRunAsNonRoot)
	hook := logTest.NewGlobal()
	results[0].Print()
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	hook.Reset()
	rootConfig.manifest = "../fixtures/unknown_type_v1.yml"
	resources, err = getKubeResourcesManifest(rootConfig.manifest)
	results = getResults(resources, auditNetworkPolicies)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	hook.Reset()
}

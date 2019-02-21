package cmd

import (
	"testing"

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
	rootConfig.manifest = "../fixtures/autofix-all-resources_v1.yml"
	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	assert.Nil(t, err)
	results := getResults(resources, auditRunAsNonRoot)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, len(results[0].Occurrences))
	assert.Equal(t, 1, len(results[1].Occurrences))
	fields := createFields(results[0], results[0].Occurrences[0])
	assert.Equal(t, 4, len(fields))
	fields = createFields(results[1], results[1].Occurrences[0])
	assert.Equal(t, 4, len(fields))
}

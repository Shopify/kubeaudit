package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareTextFilesV1(t *testing.T) {
	file1 := "../fixtures/unknown_type_v1.yml"
	file2 := "../fixtures/audit_all_v1.yml"
	assert := assert.New(t)
	filesEqual := compareTextFiles(file1, file2)
	assert.False(filesEqual)
}

func TestCompareFakeFilesV1(t *testing.T) {
	file1 := "fakeFile1.yml"
	file2 := "../fixtures/audit_all_v1.yml"
	assert := assert.New(t)
	filesEqual := compareTextFiles(file1, file2)
	assert.False(filesEqual)
}

func TestCompareFakeFilesV2(t *testing.T) {
	file1 := "./fixtures/audit_all_v1.yml"
	file2 := "fakeFile2.yml"
	assert := assert.New(t)
	filesEqual := compareTextFiles(file1, file2)
	assert.False(filesEqual)
}

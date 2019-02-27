package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutofixNetworkPoliciesFix(t *testing.T) {
	origFilename := "../fixtures/autofix-namespace-missing-default-deny-netpol.yml"
	expectedFilename := "../fixtures/autofix-namespace-missing-default-deny-netpol-fixed.yml"
	assert := assert.New(t)

	// Copy original yaml to a temp file because autofix modifies the input file
	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix_test")
	tmpFilename := tmpFile.Name()
	assert.Nil(err)
	defer os.Remove(tmpFilename)
	origFile, err := os.Open(origFilename)
	assert.Nil(err)
	_, err = io.Copy(tmpFile, origFile)
	assert.Nil(err)
	tmpFile.Close()
	origFile.Close()

	rootConfig.manifest = tmpFilename
	autofix(nil, nil)

	assert.True(compareTextFiles(expectedFilename, tmpFilename))
}

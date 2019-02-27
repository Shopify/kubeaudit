package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkPoliciesFix(t *testing.T) {
	origFilename := "../fixtures/namespace_missing_default_deny_netpol.yml"
	expectedFilename := "../fixtures/namespace_missing_default_deny_netpol-fixed.yml"
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

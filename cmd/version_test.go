package cmd

import (
	"bytes"
	"regexp"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	cmd := RootCmd
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatal("Version command returned an error.")
	}

	type exp struct {
		str string
		msg string
	}

	exps := []exp{
		{`.*level=info\s+msg="Kubeaudit version"\s+BuildDate=\S+\s+Commit=[[:xdigit:]]+\s+Version=\d+\.\d+\.\d+.*`,
			"invalid kubeaudit version"},
	}

	for _, e := range exps {
		assert.Regexp(t, regexp.MustCompile(e.str), buf, e.msg)
	}
}

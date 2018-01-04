package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCapSetFromArray(t *testing.T) {
	assert := assert.New(t)
	capArray := []Capability{"AUDIT_WRITE", "CHOWN"}
	capSet := CapSet{"AUDIT_WRITE": true, "CHOWN": true}
	assert.Equal(NewCapSetFromArray(capArray), capSet)
}

func TestMergeCapSets(t *testing.T) {
	assert := assert.New(t)
	set1 := CapSet{"AUDIT_WRITE": true, "CHOWN": true}
	set2 := CapSet{"CHOWN": true, "DAC_OVERRIDE": true}
	set3 := CapSet{"AUDIT_WRITE": true, "CHOWN": true, "DAC_OVERRIDE": true}
	assert.Equal(mergeCapSets(set1, set2), set3)
}

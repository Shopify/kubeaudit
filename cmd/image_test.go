package cmd

import "testing"

func TestImageTagMissing(t *testing.T) {
	runAuditTest(t, "image_tag_missing.yml", auditImages, []int{ErrorImageTagMissing}, "fakeContainerImg:1.6")
}

func TestImageTagIncorrect(t *testing.T) {
	runAuditTest(t, "image_tag_present.yml", auditImages, []int{ErrorImageTagIncorrect}, "fakeContainerImg:1.6")
}

func TestImageTagCorrect(t *testing.T) {
	runAuditTest(t, "image_tag_present.yml", auditImages, []int{InfoImageCorrect}, "fakeContainerImg:1.5")
}

package cmd

import "testing"

func TestImageTagMissingV1(t *testing.T) {
	runAuditTest(t, "image_tag_missing_v1.yml", auditImages, []int{ErrorImageTagMissing}, "fakeContainerImg:1.6")
}

func TestImageTagIncorrectV1(t *testing.T) {
	runAuditTest(t, "image_tag_present_v1.yml", auditImages, []int{ErrorImageTagIncorrect}, "fakeContainerImg:1.6")
}

func TestImageTagCorrectV1(t *testing.T) {
	runAuditTest(t, "image_tag_present_v1.yml", auditImages, []int{InfoImageCorrect}, "fakeContainerImg:1.5")
}

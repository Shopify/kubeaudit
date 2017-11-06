package cmd

import "testing"

func TestImageTagMissing(t *testing.T) {
	runTest(t, "image_tag_missing.yml", auditImages, ErrorImageTagMissing, "fakeContainerImg:1.6")
}

func TestImageTagIncorrect(t *testing.T) {
	runTest(t, "image_tag_present.yml", auditImages, ErrorImageTagIncorrect, "fakeContainerImg:1.6")
}

func TestImageTagCorrect(t *testing.T) {
	runTest(t, "image_tag_present.yml", auditImages, InfoImageCorrect, "fakeContainerImg:1.5")
}

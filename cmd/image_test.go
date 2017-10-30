package cmd

import "testing"

func TestImageTagMissing(t *testing.T) {
	runImageTest(t, "image_tag_missing.yml", auditImages, "fakeContainerImg:1.6", ErrorImageTagMissing)
}

func TestImageTagIncorrect(t *testing.T) {
	runImageTest(t, "image_tag_present.yml", auditImages, "fakeContainerImg:1.6", ErrorImageTagIncorrect)
}

func TestImageTagCorrect(t *testing.T) {
	runImageTest(t, "image_tag_present.yml", auditImages, "fakeContainerImg:1.5", InfoImageCorrect)
}

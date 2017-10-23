package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakePodImg")
	fakeaudit.CreateFakePodImg("fakePodImg")
}

func TestPodImg(t *testing.T) {
	fakePods := fakeaudit.GetPods("fakePodImg")
	wg.Add(2)
	image := imgFlags{img: "fakeContainerImg:1.5"}
	image.splitImageString()
	results := auditImages(image, kubeAuditPods{list: fakePods})

	if len(results) != 2 {
		t.Error("Test 1: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.Name == "fakePodImg1" && result.Occurrences[0].id != ErrorImageTagMissing {
			t.Error("Test 2: Failed to identify that image tag is missing. Refer: fakePodImg1.yml")
		}

		if result.Name == "fakePodImg2" && result.ImageTag != "1.5" {
			t.Error("Test 3: Failed to identify the correct image tag which is present. Refer: akePodImg2.yml")
		}
	}

	image = imgFlags{img: "fakeContainerImg:1.6"}
	image.splitImageString()
	results = auditImages(image, kubeAuditPods{list: fakePods})

	if len(results) != 2 {
		t.Error("Test 4: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.Name == "fakePodImg1" && result.Occurrences[0].id != ErrorImageTagMissing {
			t.Error("Test 5: Failed to identify that image tag is missing. Refer: fakePodImg1")
		}

		if result.Name == "fakePodImg2" && result.Occurrences[0].id != ErrorImageTagIncorrect {
			t.Error("Test 6: Failed to identify wrong image tag. Refer: fakePodImg2")
		}
	}

}

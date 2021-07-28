package image

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
)

const Name = "image"

const (
	// ImageTagMissing occurs when the container image tag is missing
	ImageTagMissing = "ImageTagMissing"
	// ImageTagIncorrect occurs when the container image tag does not match the user-provided value
	ImageTagIncorrect = "ImageTagIncorrect"
	// ImageCorrect occurs when the container image tag is correct
	ImageCorrect = "ImageCorrect"
)

// Image implements Auditable
type Image struct {
	image string
}

func New(config Config) *Image {
	return &Image{
		image: config.GetImage(),
	}
}

// Audit checks that the container image matches the provided image
func (image *Image) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container, image.image)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditContainer(container *k8s.ContainerV1, image string) *kubeaudit.AuditResult {
	name, tag := splitImageString(image)
	containerName, containerTag := splitImageString(container.Image)

	if isImageTagMissing(containerTag) {
		return &kubeaudit.AuditResult{
			Name:     ImageTagMissing,
			Severity: kubeaudit.Warn,
			Message:  "Image tag is missing.",
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	if isImageTagIncorrect(name, tag, containerName, containerTag) {
		return &kubeaudit.AuditResult{
			Name:     ImageTagIncorrect,
			Severity: kubeaudit.Error,
			Message:  fmt.Sprintf("Container tag is incorrect. It should be set to '%s'.", tag),
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	if isImageCorrect(name, tag, containerName, containerTag) {
		return &kubeaudit.AuditResult{
			Name:     ImageCorrect,
			Severity: kubeaudit.Info,
			Message:  "Image tag is correct",
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	return nil
}

func isImageTagMissing(tag string) bool {
	return tag == ""
}

func isImageTagIncorrect(name, tag, containerName, containerTag string) bool {
	return containerName == name && containerTag != tag
}

func isImageCorrect(name, tag, containerName, containerTag string) bool {
	return containerName == name && containerTag == tag
}

func splitImageString(image string) (name, tag string) {
	tokens := strings.Split(image, ":")
	if len(tokens) > 0 {
		name = tokens[0]
	}
	if len(tokens) > 1 {
		tag = tokens[1]
	}
	return
}

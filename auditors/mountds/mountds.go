package mountds

import (
	"fmt"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
)

const Name = "mountds"

const (
	// DockerSocketMounted occurs when a container has Docker socket mounted
	DockerSocketMounted = "DockerSocketMounted"
)

// DockerSocketPath is the mount path of the Docker socket
const DockerSocketPath = "/var/run/docker.sock"

// DockerSockMounted implements Auditable
type DockerSockMounted struct{}

func New() *DockerSockMounted {
	return &DockerSockMounted{}
}

// Audit checks that the container does not have the Docker socket mounted
func (limits *DockerSockMounted) Audit(resource k8stypes.Resource, _ []k8stypes.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		auditResult := auditContainer(container)
		if auditResult != nil {
			auditResults = append(auditResults, auditResult)
		}
	}

	return auditResults, nil
}

func auditContainer(container *k8stypes.ContainerV1) *kubeaudit.AuditResult {
	if isDockerSocketMounted(container) {
		return &kubeaudit.AuditResult{
			Name:     DockerSocketMounted,
			Severity: kubeaudit.Warn,
			Message:  fmt.Sprintf("Docker socket is mounted. '%s' should be removed from the container's volume mount list.", DockerSocketPath),
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
	}

	return nil
}

func isDockerSocketMounted(container *k8stypes.ContainerV1) bool {
	if container.VolumeMounts == nil {
		return false
	}

	for _, mount := range container.VolumeMounts {
		if mount.MountPath == DockerSocketPath {
			return true
		}
	}

	return false
}

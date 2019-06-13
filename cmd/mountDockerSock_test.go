package cmd

import "testing"

func TestDockerSockMounted(t *testing.T) {
	runAuditTest(t, "docker_sock_mounted.yml", auditMountDockerSock, []int{ErrorDockerSockMounted})
}

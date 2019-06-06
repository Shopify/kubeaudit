package cmd

import "testing"

func TestDockerSockMounted(t *testing.T) {
	runAuditTest(t, "docker_sock_mounted.yml", auditMountDockerSock, []int{ErrorDockerSockMounted})
}

func TestDockerSockMountingAllowed(t *testing.T) {
	runAuditTest(t, "allow_docker_sock_mounted.yml", auditMountDockerSock, []int{ErrorDockerSockMounted, ErrorDockerSockMountAllowed})
}

func TestDockerSockMountingAllowedPod(t *testing.T) {
	runAuditTest(t, "allow_docker_sock_mounted_pod.yml", auditMountDockerSock, []int{ErrorDockerSockMountAllowed})
}

func TestAllowMountDockerSockFromConfig(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_mount_docker_sock_from_config.yml"
	runAuditTest(t, "docker_sock_mounted.yml", auditMountDockerSock, []int{ErrorDockerSockMountAllowed})
	rootConfig.auditConfig = ""
}

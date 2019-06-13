package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DockerSockPath is the mount path of the Docker socket
const DockerSockPath = "/var/run/docker.sock"

func checkMountDockerSock(container ContainerV1, result *Result) {
	if container.VolumeMounts != nil {
		for _, mount := range container.VolumeMounts {
			if mount.MountPath == DockerSockPath {		
				occ := Occurrence{
					container: container.Name,
					id:        ErrorDockerSockMounted,
					kind:      Warn,
					message:   "/var/run/docker.sock is being mounted, please avoid this practice.",
				}
				result.Occurrences = append(result.Occurrences, occ)
			}
		}
	} 
	return
}

func auditMountDockerSock(resource Resource) (results []Result) {
	for _, container := range getContainers(resource) {
		result, err, warn := newResultFromResource(resource)
		if warn != nil {
			log.Warn(warn)
			return
		}
		if err != nil {
			log.Error(err)
			return
		}

		checkMountDockerSock(container, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
		}
	}
	return
}

var mountdsCmd = &cobra.Command{
	Use:   "mountds",
	Short: "Audit containers that mount /var/run/docker.sock",
	Long: `This command determines which containers in a kubernetes cluster
mount /var/run/docker.sock. 

A PASS is given when a container does not mount /var/run/docker.sock
A FAIL is generated when a container mounts /var/run/docker.sock

Example usage:
kubeaudit mountds`,
	Run: runAudit(auditMountDockerSock),
}

func init() {
	RootCmd.AddCommand(mountdsCmd)
}

package cmd

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkReadOnlyRootFS(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if container.SecurityContext.ReadOnlyRootFilesystem == nil {
		occ := Occurrence{id: ErrorReadOnlyRootFilesystemNIL, kind: Error, message: "ReadOnlyRootFilesystem not set which results in a writable rootFS, please set to true"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if !*container.SecurityContext.ReadOnlyRootFilesystem {
		occ := Occurrence{id: ErrorReadOnlyRootFilesystemFalse, kind: Error, message: "ReadOnlyRootFilesystem set to false, please set to true"}
		result.Occurrences = append(result.Occurrences, occ)
	}
}

func auditReadOnlyRootFS(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkReadOnlyRootFS(container, result)
			if result != nil && len(result.Occurrences) > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	for _, result := range results {
		result.Print()
	}
	return
}

var readonlyfsCmd = &cobra.Command{
	Use:   "rootfs",
	Short: "Audit containers with read only root filesystems",
	Long: `This command determines which containers in a kubernetes cluster
have their filesystems set to read only.

A PASS is given when a container has a read only root filesystem
A FAIL is given when a container does not have a read only root filesystem

Example usage:
kubeaudit runAsNonRoot`,
	Run: func(cmd *cobra.Command, args []string) {
		if rootConfig.json {
			log.SetFormatter(&log.JSONFormatter{})
		}
		var resources []Items

		if rootConfig.manifest != "" {
			var err error
			resources, err = getKubeResourcesManifest(rootConfig.manifest)
			if err != nil {
				log.Error(err)
			}
		} else {
			kube, err := kubeClient(rootConfig.kubeConfig)
			if err != nil {
				log.Error(err)
			}
			resources = getKubeResources(kube)
		}

		var wg sync.WaitGroup
		wg.Add(len(resources))

		for _, resource := range resources {
			go func(items Items) {
				auditReadOnlyRootFS(items)
				wg.Done()
			}(resource)
		}

		wg.Wait()
	},
}

func init() {
	securityContextCmd.AddCommand(readonlyfsCmd)
}

package cmd

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkRunAsNonRoot(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if container.SecurityContext.RunAsNonRoot == nil {
		occ := Occurrence{id: ErrorRunAsNonRootNIL, kind: Error, message: "RunAsNonRoot is not set, which results in root user being allowed!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if *container.SecurityContext.RunAsNonRoot == false {
		occ := Occurrence{id: ErrorRunAsNonRootFalse, kind: Error, message: "RunAsNonRoot is set to false (root user allowed), please set to true!"}
		result.Occurrences = append(result.Occurrences, occ)
	}
}

func auditRunAsNonRoot(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkRunAsNonRoot(container, result)
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

// runAsNonRootCmd represents the runAsNonRoot command
var runAsNonRootCmd = &cobra.Command{
	Use:   "nonroot",
	Short: "Audit containers running as root",
	Long: `This command determines which containers in a kubernetes cluster
are running as root (uid=0).

A PASS is given when a container runs as a uid greater than 0
A FAIL is generated when a container runs as root

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
				auditRunAsNonRoot(items)
				wg.Done()
			}(resource)
		}

		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(runAsNonRootCmd)
}

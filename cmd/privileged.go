package cmd

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkPrivileged(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if container.SecurityContext.Privileged == nil {
		// TODO find out what this exactly means
		occ := Occurrence{id: ErrorPrivilegedNIL, kind: Warn, message: "Privileged defaults to false, which results in non privileged, which is okay."}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if *container.SecurityContext.Privileged == true {
		occ := Occurrence{id: ErrorPrivilegedTrue, kind: Error, message: "Privileged set to true! Please change it to false!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
}

func auditPrivileged(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkPrivileged(container, result)
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
var privileged = &cobra.Command{
	Use:   "privileged",
	Short: "Audit containers running as root",
	Long: `This command determines which containers in a kubernetes cluster
are running as privileged.

A PASS is given when a container runs in a non-privileged mode
A FAIL is generated when a container runs in a privileged mode

Example usage:
kubeaudit privileged`,
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
			go func() {
				auditPrivileged(resource)
				wg.Done()
			}()
		}

		wg.Wait()
	},
}

func init() {
	securityContextCmd.AddCommand(privileged)
}

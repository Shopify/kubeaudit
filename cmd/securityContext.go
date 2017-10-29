package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkSecurityContext(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if container.SecurityContext.Capabilities == nil {
		occ := Occurrence{id: ErrorCapabilitiesNIL, kind: Error, message: "Capabilites field not defined!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if container.SecurityContext.Capabilities.Add != nil {
		result.CapsAdded = container.SecurityContext.Capabilities.Add
		occ := Occurrence{id: ErrorCapabilitiesAdded, kind: Error, message: "Capabilities were added!"}
		result.Occurrences = append(result.Occurrences, occ)
	}

	if container.SecurityContext.Capabilities.Drop == nil {
		occ := Occurrence{id: ErrorCapabilitiesNoneDropped, kind: Error, message: "No capabilities were dropped!"}
		result.Occurrences = append(result.Occurrences, occ)
	}

	if container.SecurityContext.Capabilities.Drop != nil {
		// TODO need a check for which caps have been dropped and whether that's an
		// error because not enough have been dropped
		result.CapsDropped = container.SecurityContext.Capabilities.Drop
		occ := Occurrence{id: ErrorCapabilitiesSomeDropped, kind: Error, message: "Not all of the capabilities were dropped!"}
		result.Occurrences = append(result.Occurrences, occ)
	}
}

func auditSecurityContext(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkSecurityContext(container, result)
			if result != nil && len(result.Occurrences) > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	for _, result := range results {
		result.Print()
	}
	defer wg.Done()
	return
}

var securityContextCmd = &cobra.Command{
	Use:   "sc",
	Short: "Audit container security contexts",
	Long: `This command determines which containers in a kubernetes cluster
are running as root.
An INFO log is given when a container has a securityContext
An ERROR log is generated when a container does not have a defined securityContext
A WARN log is generated when some linux capabilities are added or not dropped
This command is also a root command, check kubeaudit sc --help
Example usage:
kubeaudit sc
kubeaudit sc nonroot
kubeaudit sc rootfs`,
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

		count := len(resources)
		wg.Add(count)
		for _, resource := range resources {
			go auditSecurityContext(resource)
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(securityContextCmd)
}

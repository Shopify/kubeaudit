package cmd

import (
	"io/ioutil"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type capsDropList struct {
	Drop []string `yaml:"capabilitiesToBeDropped"`
}

func recommendedCapabilitiesToBeDropped() (dropList []Capability, err error) {
	yamlFile, err := ioutil.ReadFile("config/capabilities-drop-list.yml")
	if err != nil {
		return
	}
	caps := capsDropList{}
	err = yaml.Unmarshal(yamlFile, &caps)
	if err != nil {
		return
	}
	for _, drop := range caps.Drop {
		dropList = append(dropList, Capability(drop))
	}
	return
}

func capsNotDropped(dropped []Capability) (notDropped []Capability, err error) {
	toBeDropped, err := recommendedCapabilitiesToBeDropped()
	if err != nil {
		return
	}
	for _, toBeDroppedCap := range toBeDropped {
		found := false
		for _, droppedCap := range dropped {
			if toBeDroppedCap == droppedCap {
				found = true
			}
		}
		if found == false {
			notDropped = append(notDropped, toBeDroppedCap)
		}
	}
	return
}

func checkSecurityContext(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if container.SecurityContext.Capabilities == nil {
		occ := Occurrence{id: ErrorCapabilitiesNIL, kind: Error, message: "Capabilities field not defined!"}
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
		capsNotDropped, err := capsNotDropped(container.SecurityContext.Capabilities.Drop)
		if err != nil {
			occ := Occurrence{id: KubeauditInternalError, kind: Error, message: "This should not have happened, if you are on kubeaudit master please consider to report: " + err.Error()}
			result.Occurrences = append(result.Occurrences, occ)
			return
		}
		if len(capsNotDropped) > 0 {
			result.CapsNotDropped = capsNotDropped
			occ := Occurrence{id: ErrorCapabilitiesSomeDropped, kind: Error, message: "Not all of the recommended capabilities were dropped! Please drop the mentioned capabiliites."}
			result.Occurrences = append(result.Occurrences, occ)
		}
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

		var wg sync.WaitGroup
		wg.Add(len(resources))

		for _, resource := range resources {
			go func() {
				auditSecurityContext(resource)
				wg.Done()
			}()
		}

		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(securityContextCmd)
}

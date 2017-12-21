package cmd

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
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

func checkCapabilities(container Container, result *Result) {
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

func auditCapabilities(resource k8sRuntime.Object) (results []Result) {
	for _, container := range getContainers(resource) {
		result := newResultFromResource(resource)
		checkCapabilities(container, &result)
		if len(result.Occurrences) > 0 {
			results = append(results, result)
			break
		}
	}
	return
}

var capabilitiesCmd = &cobra.Command{
	Use:   "caps",
	Short: "Audit container for capabilities",
	Run:   runAudit(auditCapabilities),
}

func init() {
	RootCmd.AddCommand(capabilitiesCmd)
}

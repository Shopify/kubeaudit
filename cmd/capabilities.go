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

func recommendedCapabilitiesToBeDropped() (dropCapSet CapSet, err error) {
	yamlFile, err := ioutil.ReadFile("config/capabilities-drop-list.yml")
	if err != nil {
		return
	}
	caps := capsDropList{}
	err = yaml.Unmarshal(yamlFile, &caps)
	if err != nil {
		return
	}
	dropCapSet = make(CapSet)
	for _, drop := range caps.Drop {
		dropCapSet[Capability(drop)] = true
	}
	return
}

func checkCapabilities(container Container, result *Result) {
	added := CapSet{}
	dropped := CapSet{}
	if container.SecurityContext != nil && container.SecurityContext.Capabilities != nil {
		added = NewCapSetFromArray(container.SecurityContext.Capabilities.Add)
		dropped = NewCapSetFromArray(container.SecurityContext.Capabilities.Drop)
	}

	allowedMap := result.allowedCaps()
	allowed := make(CapSet)
	for k := range allowedMap {
		allowed[k] = true
	}

	toBeDropped, err := recommendedCapabilitiesToBeDropped()
	if err != nil {
		occ := Occurrence{
			container: container.Name,
			id:        KubeauditInternalError,
			kind:      Error,
			message:   "This should not have happened, if you are on kubeaudit master please consider to report: " + err.Error(),
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	for _, cap := range sortCapSet(mergeCapSets(toBeDropped, dropped, allowed, added)) {
		if !allowed[cap] && !dropped[cap] && toBeDropped[cap] {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorCapabilityNotDropped,
				kind:      Error,
				message:   "Capability not dropped",
				metadata:  Metadata{"CapName": string(cap)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if !allowed[cap] && added[cap] {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorCapabilityAdded,
				kind:      Error,
				message:   "Capability added",
				metadata:  Metadata{"CapName": string(cap)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if allowed[cap] && (toBeDropped[cap] && !dropped[cap] || added[cap]) {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorCapabilityAllowed,
				kind:      Warn,
				message:   "Capability allowed",
				metadata: Metadata{
					"CapName": string(cap),
					"Reason":  prettifyReason(allowedMap[cap]),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if allowed[cap] && !(toBeDropped[cap] && !dropped[cap] || added[cap]) {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorMisconfiguredKubeauditAllow,
				kind:      Warn,
				message:   "Capability allowed but not present",
				metadata: Metadata{
					"CapName": string(cap),
					"Reason":  allowedMap[cap],
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}
}

func auditCapabilities(resource k8sRuntime.Object) (results []Result) {
	for _, container := range getContainers(resource) {
		result, err := newResultFromResource(resource)
		if err == ErrResourceTypeNotSupported {
			continue
		} else if err != nil {
			panic(err)
		}

		checkCapabilities(container, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
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

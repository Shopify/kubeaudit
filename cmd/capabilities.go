package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

type capsDropList struct {
	Drop []string `yaml:"capabilitiesToBeDropped"`
}

type CapSet map[Capability]bool

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

func allowedCaps(result *Result) (allowed map[Capability]string) {
	allowed = make(map[Capability]string)
	for k, v := range result.Labels {
		if strings.Contains(k, "kubeaudit.allow.capability.") {
			allowed[Capability(strings.ToUpper(strings.TrimPrefix(k, "kubeaudit.allow.capability.")))] = v
		}
	}
	return
}

func arrayToCapSet(array []Capability) (set CapSet) {
	set = make(CapSet)
	for _, cap := range array {
		set[cap] = true
	}
	return
}

func mergeCapSets(sets ...CapSet) (merged CapSet) {
	merged = make(CapSet)
	for _, set := range sets {
		for k, v := range set {
			merged[k] = v
		}
	}
	return
}

func checkCapabilities(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{
			id:      ErrorSecurityContextNIL,
			kind:    Error,
			message: "SecurityContext not set, please set it!",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if container.SecurityContext.Capabilities == nil {
		occ := Occurrence{
			id:      ErrorCapabilitiesNIL,
			kind:    Error,
			message: "Capabilities field not defined!",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	added := arrayToCapSet(container.SecurityContext.Capabilities.Add)
	dropped := arrayToCapSet(container.SecurityContext.Capabilities.Drop)
	allowedMap := allowedCaps(result)
	allowed := make(CapSet)
	for k := range allowedMap {
		allowed[k] = true
	}
	toBeDropped, err := recommendedCapabilitiesToBeDropped()
	if err != nil {
		occ := Occurrence{
			id:      KubeauditInternalError,
			kind:    Error,
			message: "This should not have happened, if you are on kubeaudit master please consider to report: " + err.Error(),
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	for cap := range mergeCapSets(toBeDropped, dropped, allowed, added) {
		if !allowed[cap] && !dropped[cap] && toBeDropped[cap] {
			occ := Occurrence{
				id:       ErrorCapabilityNotDropped,
				kind:     Error,
				message:  "Capability not dropped",
				metadata: Metadata{"CapName": string(cap)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if !allowed[cap] && added[cap] {
			occ := Occurrence{
				id:       ErrorCapabilityAdded,
				kind:     Error,
				message:  "Capability added",
				metadata: Metadata{"CapName": string(cap)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if allowed[cap] && (toBeDropped[cap] && !dropped[cap] || added[cap]) {
			occ := Occurrence{
				id:      ErrorCapabilityAllowed,
				kind:    Warn,
				message: "Capability allowed",
				metadata: Metadata{
					"CapName": string(cap),
					"Reason":  prettifyReason(allowedMap[cap]),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if allowed[cap] && !(toBeDropped[cap] && !dropped[cap] || added[cap]) {
			occ := Occurrence{
				id:      ErrorMisconfiguredKubeauditAllow,
				kind:    Warn,
				message: "Capability allowed but not present",
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

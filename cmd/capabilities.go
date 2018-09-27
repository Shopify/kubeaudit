package cmd

import (
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

type capsDropList struct {
	Drop []string `yaml:"capabilitiesToBeDropped"`
}

const defaultDropCapConfig = `
# SANE DEFAULTS:
capabilitiesToBeDropped:
  # https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities
  - SETPCAP #Modify process capabilities.
  - MKNOD #Create special files using mknod(2).
  - AUDIT_WRITE #Write records to kernel auditing log.
  - CHOWN #Make arbitrary changes to file UIDs and GIDs (see chown(2)).
  - NET_RAW #Use RAW and PACKET sockets.
  - DAC_OVERRIDE #Bypass file read, write, and execute permission checks.
  - FOWNER #Bypass permission checks on operations that normally require the file system UID of the process to match the UID of the file.
  - FSETID #Don’t clear set-user-ID and set-group-ID permission bits when a file is modified.
  - KILL #Bypass permission checks for sending signals.
  - SETGID #Make arbitrary manipulations of process GIDs and supplementary GID list.
  - SETUID #Make arbitrary manipulations of process UIDs.
  - NET_BIND_SERVICE #Bind a socket to internet domain privileged ports (port numbers less than 1024).
  - SYS_CHROOT #Use chroot(2), change root directory.
  - SETFCAP #Set file capabilities.
`

func recommendedCapabilitiesToBeDropped() (dropCapSet CapSet, err error) {
	yamlFile := []byte(defaultDropCapConfig)
	if rootConfig.dropCapConfig != "" {
		if _, err = os.Stat(rootConfig.dropCapConfig); err != nil {
			return
		}
		yamlFile, err = ioutil.ReadFile(rootConfig.dropCapConfig)
		if err != nil {
			return
		}
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
	allCapsDrop := false
	if container.SecurityContext != nil && container.SecurityContext.Capabilities != nil {
		added = NewCapSetFromArray(container.SecurityContext.Capabilities.Add)
		if len(container.SecurityContext.Capabilities.Drop) != 0 &&
			strings.ToLower(string(container.SecurityContext.Capabilities.Drop[0])) == "all" {
			allCapsDrop = true
		} else {
			dropped = NewCapSetFromArray(container.SecurityContext.Capabilities.Drop)
		}
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
	if allCapsDrop {
		dropped = toBeDropped
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
		if err != nil {
			log.Error(err)
			return
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

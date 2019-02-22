package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkRunAsNonRootCSC(container ContainerV1, result *Result) {
	if labelExists, reason := getContainerOverrideLabelReason(result, container, "allow-run-as-root"); labelExists {
		if container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil || *container.SecurityContext.RunAsNonRoot == false {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorRunAsNonRootFalseAllowed,
				kind:      Warn,
				message:   "Allowed setting RunAsNonRoot to false in ContainerSecurityContext",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorMisconfiguredKubeauditAllow,
				kind:      Warn,
				message:   "Allowed setting RunAsNonRoot to false, but it is set to true",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorRunAsNonRootPSCNilCSCNil,
			kind:      Error,
			message:   "RunAsNonRoot is not set in ContainerSecurityContext, which results in root user being allowed!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else if *container.SecurityContext.RunAsNonRoot == false {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorRunAsNonRootPSCTrueFalseCSCFalse,
			kind:      Error,
			message:   "RunAsNonRoot is set to false (root user allowed) in ContainerSecurityContext, please set to true!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

// Checks the PodSecurityContext for RANR

func checkRunAsNonRootPSC(podSpec PodSpecV1, container ContainerV1, result *Result) {
	if labelExists, reason := getContainerOverrideLabelReason(result, container, "allow-run-as-root"); labelExists {
		if podSpec.SecurityContext == nil || podSpec.SecurityContext.RunAsNonRoot == nil || *podSpec.SecurityContext.RunAsNonRoot == false {
			occ := Occurrence{
				container: container.Name,
				podHost:   podSpec.Hostname,
				id:        ErrorRunAsNonRootFalseAllowed,
				kind:      Warn,
				message:   "Allowed setting RunAsNonRoot to false",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				container: container.Name,
				podHost:   podSpec.Hostname,
				id:        ErrorMisconfiguredKubeauditAllow,
				kind:      Warn,
				message:   "Allowed setting RunAsNonRoot to false, but it is set to true",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if *podSpec.SecurityContext.RunAsNonRoot == false {
		occ := Occurrence{
			container: container.Name,
			podHost:   podSpec.Hostname,
			id:        ErrorRunAsNonRootPSCFalseCSCNil,
			kind:      Error,
			message:   "RunAsNonRoot is set to false (root user allowed) in PodsSecurityContext and not set in ContainerSecurityContext, please set to true!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditRunAsNonRoot(resource Resource) (results []Result) {
	// get PodSpec for PodSecurityContext
	podSpec := getPodSpecs(resource)
	for _, container := range getContainers(resource) {
		result, err := newResultFromResource(resource)
		if err != nil {
			log.Error(err)
			return
		}

		// check if ContainerSecurityContext is defined properly, else audit the PodSecurityContext
		if shouldAuditCSC(podSpec, container) {
			checkRunAsNonRootCSC(container, result)
		} else {
			checkRunAsNonRootPSC(podSpec, container, result)
		}
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
		}
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
kubeaudit nonroot`,
	Run: runAudit(auditRunAsNonRoot),
}

func init() {
	RootCmd.AddCommand(runAsNonRootCmd)
}

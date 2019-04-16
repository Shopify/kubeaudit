package cmd

import (
	"github.com/spf13/cobra"
)

// Checks the PodSecurityContext for NIX
func checkNamespaces(podSpec PodSpecV1, result *Result) {
	if labelExists, reason := getPodOverrideLabelReason(result, "allow-namespace-host-network"); labelExists {
		if podSpec.HostNetwork {
			occ := Occurrence{
				podHost:  podSpec.Hostname,
				id:       ErrorNamespaceHostNetworkTrueAllowed,
				kind:     Warn,
				message:  "Allowed setting hostNetwork to true",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				podHost:  podSpec.Hostname,
				id:       ErrorMisconfiguredKubeauditAllow,
				kind:     Warn,
				message:  "Allowed setting hostNetwork to true, but it is set to false",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if podSpec.HostNetwork {
		occ := Occurrence{
			podHost: podSpec.Hostname,
			id:      ErrorNamespaceHostNetworkTrue,
			kind:    Error,
			message: "hostNetwork is set to true in podSpec, please set to false!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	if labelExists, reason := getPodOverrideLabelReason(result, "allow-namespace-host-IPC"); labelExists {
		if podSpec.HostIPC {
			occ := Occurrence{
				podHost:  podSpec.Hostname,
				id:       ErrorNamespaceHostIPCTrueAllowed,
				kind:     Warn,
				message:  "Allowed setting hostIPC to true",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				podHost:  podSpec.Hostname,
				id:       ErrorMisconfiguredKubeauditAllow,
				kind:     Warn,
				message:  "Allowed setting hostIPC to true, but it is set to false",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if podSpec.HostIPC {
		occ := Occurrence{
			podHost: podSpec.Hostname,
			id:      ErrorNamespaceHostIPCTrue,
			kind:    Error,
			message: "hostIPC is set to true  in podSpec, please set to false!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	if labelExists, reason := getPodOverrideLabelReason(result, "allow-namespace-host-PID"); labelExists {
		if podSpec.HostPID {
			occ := Occurrence{
				podHost:  podSpec.Hostname,
				id:       ErrorNamespaceHostPIDTrueAllowed,
				kind:     Warn,
				message:  "Allowed setting hostPID to true",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				podHost:  podSpec.Hostname,
				id:       ErrorMisconfiguredKubeauditAllow,
				kind:     Warn,
				message:  "Allowed setting hostPID to true, but it is set to false",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if podSpec.HostPID {
		occ := Occurrence{
			podHost: podSpec.Hostname,
			id:      ErrorNamespaceHostPIDTrue,
			kind:    Error,
			message: "hostPID is set to true  in podSpec, please set to false!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditNamespaces(resource Resource) (results []Result) {
	switch kubeType := resource.(type) {
	case *PodV1:
		podSpec := kubeType.Spec
		result, _, _ := newResultFromResource(resource)
		checkNamespaces(podSpec, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
		}
	}
	return
}

// runAsNonRootCmd represents the runAsNonRoot command
var namespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "Audit Pods for hostNetwork, hostIPC and hostPID",
	Long: `This command determines which pods in a kubernetes cluster
are running with hostNetwork, hostIPC or hostPID set to true.
	
A PASS is given when a pod has hostNetwork, hostIPC and hostPID set to false or not set
A FAIL is generated when a pod has at least one of hostNetwork, hostIPC or hostPID set to true

kubeaudit namespaces`,
	Run: runAudit(auditNamespaces),
}

func init() {
	RootCmd.AddCommand(namespacesCmd)
}

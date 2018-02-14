package cmd

import (
	"github.com/spf13/cobra"

	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func checkPrivileged(container Container, result *Result) {
	if container.SecurityContext == nil || container.SecurityContext.Privileged == nil {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorPrivilegedNIL,
			kind:      Warn,
			message:   "Privileged defaults to false, which results in non privileged, which is okay.",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else if allowed("privileged", container.Name, result.Labels) {
		if *container.SecurityContext.Privileged == true {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorPrivilegedTrueAllowed,
				kind:      Warn,
				message:   "Allowed setting privileged to true",
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorMisconfiguredKubeauditAllow,
				kind:      Warn,
				message:   "Allowed setting privileged to true, but privileged is false or nil",
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if *container.SecurityContext.Privileged == true {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorPrivilegedTrue,
			kind:      Error,
			message:   "Privileged set to true! Please change it to false!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditPrivileged(resource k8sRuntime.Object) (results []Result) {
	for _, container := range getContainers(resource) {
		result := newResultFromResource(resource)
		checkPrivileged(container, &result)
		if len(result.Occurrences) > 0 {
			results = append(results, result)
			break
		}
	}
	return
}

var privileged = &cobra.Command{
	Use:   "priv",
	Short: "Audit containers running as privileged",
	Long: `This command determines which containers in a kubernetes cluster
are running as privileged.

A PASS is given when a container runs in a non-privileged mode
A FAIL is generated when a container runs in a privileged mode

Example usage:
kubeaudit privileged`,
	Run: runAudit(auditPrivileged),
}

func init() {
	RootCmd.AddCommand(privileged)
}

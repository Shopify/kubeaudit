package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// As of Oct 1, 2018 these constants are not in the K8s API package, but once they are they should be replaced
// https://github.com/kubernetes/kubernetes/blob/7f23a743e8c23ac6489340bbb34fa6f1d392db9d/pkg/security/apparmor/helpers.go#L25
const (
	// The prefix to an annotation key specifying a container profile.
	ContainerAnnotationKeyPrefix = "container.apparmor.security.beta.kubernetes.io/"

	// The profile specifying the runtime default.
	ProfileRuntimeDefault = "runtime/default"
	// The prefix for specifying profiles loaded on the node.
	ProfileNamePrefix = "localhost/"
)

func checkAppArmor(resource Resource, result *Result) {
	annotations := getPodAnnotations(resource)

	for _, container := range getContainers(resource) {
		containerAnnotation := ContainerAnnotationKeyPrefix + container.Name
		containerProfile, ok := annotations[containerAnnotation]

		if !ok {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorAppArmorAnnotationMissing,
				kind:      Error,
				message:   fmt.Sprintf("AppArmor annotation missing."),
			}
			result.Occurrences = append(result.Occurrences, occ)
			continue
		}

		if badAppArmorProfileName(containerProfile) {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorAppArmorDisabled,
				kind:      Error,
				message:   fmt.Sprintf("AppArmor disabled."),
				metadata: Metadata{
					"Annotation": containerAnnotation,
					"Reason":     prettifyReason(containerProfile),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}
}

func badAppArmorProfileName(profileName string) bool {
	return profileName != ProfileRuntimeDefault && !strings.HasPrefix(profileName, ProfileNamePrefix)
}

func auditAppArmor(resource Resource) (results []Result) {
	result, err, warn := newResultFromResource(resource)
	if warn != nil {
		log.Warn(warn)
		return
	}
	if err != nil {
		log.Error(err)
		return
	}

	checkAppArmor(resource, result)

	if len(result.Occurrences) > 0 {
		results = append(results, *result)
	}

	return
}

var appArmor = &cobra.Command{
	Use:   "apparmor",
	Short: "Audit containers running without AppArmor",
	Long: `This command determines which containers in a kubernetes cluster
are running without AppArmor enabled.

A PASS is given when all containers have AppArmor enabled.
A FAIL is generated when a container has AppArmor disabled or misconfigured.

Example usage:
kubeaudit apparmor`,
	Run: runAudit(auditAppArmor),
}

func init() {
	RootCmd.AddCommand(appArmor)
}

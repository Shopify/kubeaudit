package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	apiv1 "k8s.io/api/core/v1"
)

func checkSeccomp(resource Resource, result *Result) {
	annotations := getPodAnnotations(resource)
	podAnnotation := apiv1.SeccompPodAnnotationKey
	podProfile, podOk := annotations[podAnnotation]

	if podOk {
		if badSeccompProfileName(podProfile) {
			occ := Occurrence{
				container: "",
				id:        ErrorSeccompDisabledPod,
				kind:      Error,
				message:   fmt.Sprintf("Seccomp disabled for pod."),
				metadata: Metadata{
					"Annotation": podAnnotation,
					"Reason":     prettifyReason(podProfile),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
			return
		}

		if podProfile == apiv1.DeprecatedSeccompProfileDockerDefault {
			occ := Occurrence{
				container: "",
				id:        ErrorSeccompDeprecatedPod,
				kind:      Warn,
				message:   fmt.Sprintf("Seccomp annotation for pod set to a deprecated value."),
				metadata: Metadata{
					"Annotation": podAnnotation,
					"Reason":     prettifyReason(podProfile),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}

	for _, container := range getContainers(resource) {
		containerAnnotation := apiv1.SeccompContainerAnnotationKeyPrefix + container.Name
		containerProfile, containerOk := annotations[containerAnnotation]

		if !containerOk && !podOk {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorSeccompAnnotationMissing,
				kind:      Error,
				message:   fmt.Sprintf("Seccomp annotation missing."),
			}
			result.Occurrences = append(result.Occurrences, occ)
			continue
		}

		if containerOk && badSeccompProfileName(containerProfile) {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorSeccompDisabled,
				kind:      Error,
				message:   fmt.Sprintf("Seccomp disabled for container."),
				metadata: Metadata{
					"Annotation": containerAnnotation,
					"Reason":     prettifyReason(containerProfile),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
			continue
		}

		if containerOk && containerProfile == apiv1.DeprecatedSeccompProfileDockerDefault {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorSeccompDeprecated,
				kind:      Warn,
				message:   fmt.Sprintf("Seccomp annotation for container set to a deprecated value."),
				metadata: Metadata{
					"Annotation": containerAnnotation,
					"Reason":     prettifyReason(containerProfile),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}
}

func badSeccompProfileName(profileName string) bool {
	switch {
	case profileName == apiv1.SeccompProfileRuntimeDefault:
		return false
	case profileName == apiv1.DeprecatedSeccompProfileDockerDefault:
		return false
	case strings.HasPrefix(profileName, ProfileNamePrefix):
		return false
	default:
		return true
	}
}

func auditSeccomp(resource Resource) (results []Result) {
	result, err := newResultFromResource(resource)
	if err != nil {
		log.Error(err)
		return
	}

	checkSeccomp(resource, result)

	if len(result.Occurrences) > 0 {
		results = append(results, *result)
	}

	return
}

var seccomp = &cobra.Command{
	Use:   "seccomp",
	Short: "Audit containers running without Seccomp",
	Long: `This command determines which containers in a kubernetes cluster
are running without Seccomp enabled.

A PASS is given when all containers have Seccomp enabled.
A FAIL is generated when a container has Seccomp disabled or misconfigured.

Example usage:
kubeaudit seccomp`,
	Run: runAudit(auditSeccomp, getResources),
}

func init() {
	RootCmd.AddCommand(seccomp)
}

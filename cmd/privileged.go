package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
)

func printResultPrivileged(results []Result) {
	for _, result := range results {
		if result.err > 0 {
			log.WithField("type", result.kubeType).Error(result.namespace, "/", result.name)
		}
	}
}

func checkPrivileged(container apiv1.Container, result *Result) {
	if container.SecurityContext != nil {
		if container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			result.err = 1
		}
	} else {
		result.err = 2
	}

	return
}

func auditPrivileged(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkPrivileged(container, result)
			if result != nil && result.err > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	printResultPrivileged(results)
	defer wg.Done()
	return
}

// runAsNonRootCmd represents the runAsNonRoot command
var privileged = &cobra.Command{
	Use:   "privileged",
	Short: "Audit containers running as root",
	Long: `This command determines which containers in a kubernetes cluster
are running as privileged.

A PASS is given when a container runs in a non-privileged mode
A FAIL is generated when a container runs in a privileged mode

Example usage:
kubeaudit privileged`,
	Run: func(cmd *cobra.Command, args []string) {
		if rootConfig.json {
			log.SetFormatter(&log.JSONFormatter{})
		}

		if rootConfig.manifest != "" {
			resources, err := getKubeResources(rootConfig.manifest)
			if err != nil {
				log.Error(err)
			}
			count := len(resources)
			wg.Add(count)
			for _, resource := range resources {
				go auditSecurityContext(resource)
			}
			wg.Wait()
		} else {
			kube, err := kubeClient(rootConfig.kubeConfig)
			if err != nil {
				log.Error(err)
			}

			// fetch deployments, statefulsets, daemonsets
			// and pods which do not belong to another abstraction
			deployments := getDeployments(kube)
			statefulSets := getStatefulSets(kube)
			daemonSets := getDaemonSets(kube)
			pods := getPods(kube)
			replicationControllers := getReplicationControllers(kube)

			wg.Add(5)
			go auditPrivileged(kubeAuditStatefulSets{list: statefulSets})
			go auditPrivileged(kubeAuditDaemonSets{list: daemonSets})
			go auditPrivileged(kubeAuditPods{list: pods})
			go auditPrivileged(kubeAuditReplicationControllers{list: replicationControllers})
			go auditPrivileged(kubeAuditDeployments{list: deployments})
			wg.Wait()
		}
	},
}

func init() {
	securityContextCmd.AddCommand(privileged)
}

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
)

func printResultNR(results []Result) {
	for _, result := range results {
		if result.err > 0 {
			log.WithField("type", result.kubeType).Error(result.namespace, "/", result.name)
		}
	}
}

func checkRunAsNonRoot(container apiv1.Container, result *Result) {
	if container.SecurityContext != nil {
		if container.SecurityContext.RunAsNonRoot == nil {
			result.err = 1
		} else if !*container.SecurityContext.RunAsNonRoot {
			result.err = 2
		}
	} else {
		result.err = 3
	}

	return
}

func auditRunAsNonRoot(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkRunAsNonRoot(container, result)
			if result != nil && result.err > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	printResultNR(results)
	defer wg.Done()
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
kubeaudit runAsNonRoot`,
	Run: func(cmd *cobra.Command, args []string) {
		if rootConfig.json {
			log.SetFormatter(&log.JSONFormatter{})
		}

		if rootConfig.manifest != "" {
			wg.Add(1)
			resource := getKubeResource(rootConfig.manifest)
			auditSecurityContext(resource)
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
			go auditRunAsNonRoot(kubeAuditStatefulSets{list: statefulSets})
			go auditRunAsNonRoot(kubeAuditDaemonSets{list: daemonSets})
			go auditRunAsNonRoot(kubeAuditPods{list: pods})
			go auditRunAsNonRoot(kubeAuditReplicationControllers{list: replicationControllers})
			go auditRunAsNonRoot(kubeAuditDeployments{list: deployments})
			wg.Wait()
		}
	},
}

func init() {
	securityContextCmd.AddCommand(runAsNonRootCmd)
}

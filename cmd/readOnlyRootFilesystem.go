package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
)

func printResultRFS(results []Result) {
	for _, result := range results {
		switch result.err {
		case ErrorSecurityContextNIL:
			log.WithFields(log.Fields{
				"type": result.kubeType,
			}).Error("SecurityContext not set, please set it!", result.namespace, "/", result.name)
		case ErrorReadOnlyRootFilesystemNIL:
			log.WithFields(log.Fields{
				"type": result.kubeType,
			}).Error("ReadOnlyFilesystem not set, defaults to nil", result.namespace, "/", result.name)
		case ErrorReadOnlyRootFilesystemFalse:
			log.WithFields(log.Fields{
				"type": result.kubeType,
			}).Error("ReadOnlyFilesystem set to false", result.namespace, "/", result.name)
		}
	}
}

func checkReadOnlyRootFS(container apiv1.Container, result *Result) {
	if container.SecurityContext == nil {
		result.err = ErrorSecurityContextNIL
		return
	}
	if container.SecurityContext.ReadOnlyRootFilesystem == nil {
		result.err = ErrorReadOnlyRootFilesystemNIL
		return
	}
	if !*container.SecurityContext.ReadOnlyRootFilesystem {
		result.err = ErrorReadOnlyRootFilesystemFalse
		return
	}
	return
}

func auditReadOnlyRootFS(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkReadOnlyRootFS(container, result)
			if result != nil && result.err > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	printResultRFS(results)
	defer wg.Done()
	return
}

var readonlyfsCmd = &cobra.Command{
	Use:   "rootfs",
	Short: "Audit containers with read only root filesystems",
	Long: `This command determines which containers in a kubernetes cluster
have their filesystems set to read only.

A PASS is given when a container has a read only root filesystem
A FAIL is given when a container does not have a read only root filesystem

Example usage:
kubeaudit runAsNonRoot`,
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
				go auditReadOnlyRootFS(resource)
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
			go auditReadOnlyRootFS(kubeAuditStatefulSets{list: statefulSets})
			go auditReadOnlyRootFS(kubeAuditDaemonSets{list: daemonSets})
			go auditReadOnlyRootFS(kubeAuditPods{list: pods})
			go auditReadOnlyRootFS(kubeAuditReplicationControllers{list: replicationControllers})
			go auditReadOnlyRootFS(kubeAuditDeployments{list: deployments})
			wg.Wait()
		}
	},
}

func init() {
	securityContextCmd.AddCommand(readonlyfsCmd)
}

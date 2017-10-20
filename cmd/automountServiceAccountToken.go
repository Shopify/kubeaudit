package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func printResultASAT(results []Result) {
	for _, result := range results {
		switch result.err {
		case ErrorServiceAccountTokenDeprecated:
			log.WithFields(log.Fields{
				"type":               result.kubeType,
				"namespace":          result.namespace,
				"name":               result.name,
				"serviceAccount":     result.dsa,
				"serviceAccountName": result.sa,
			}).Warn("deprecated serviceAccount detected (sub for serviceAccountName)")
		case ErrorServiceAccountTokenTrueAndNoName:
			log.WithFields(log.Fields{
				"type":      result.kubeType,
				"namespace": result.namespace,
				"name":      result.name,
			}).Error("automountServiceAccountToken = true with no serviceAccountName")
		case ErrorServiceAccountTokenNILAndNoName:
			log.WithFields(log.Fields{
				"type":      result.kubeType,
				"namespace": result.namespace,
				"name":      result.name,
			}).Error("automountServiceAccountToken nil (mounted by default) with no serviceAccountName")
		}
	}
}

func checkAutomountServiceAccountToken(result *Result) {
	// Check for use of deprecated service account name
	if result.dsa != "" {
		result.err = ErrorServiceAccountTokenDeprecated
		return
	}

	if result.token != nil && *result.token && result.sa == "" {
		// automountServiceAccountToken = true, and serviceAccountName is blank (default: default)
		result.err = ErrorServiceAccountTokenTrueAndNoName
		return
	}

	if result.token == nil && result.sa == "" {
		// automountServiceAccountToken = nil (default: true), and serviceAccountName is blank (default: default)
		result.err = ErrorServiceAccountTokenNILAndNoName
		return
	}
}

func auditAutomountServiceAccountToken(items Items) (results []Result) {
	for _, item := range items.Iter() {
		result := ServiceAccountIter(item)
		checkAutomountServiceAccountToken(result)

		if result.err > 0 {
			results = append(results, *result)
		}
	}

	printResultASAT(results)
	defer wg.Done()
	return
}

// satCmd represents the sat command
var satCmd = &cobra.Command{
	Use:   "sat",
	Short: "Audit automountServiceAccountToken = true pods against an empty (default) service account",
	Long: `This command determines which pods are running with
autoMountServiceAcccountToken = true and default service account names.

An ERROR log is generated when a container matches one of the fol:
  automountServiceAccountToken = true and serviceAccountName is blank (default: default)
  automountServiceAccountToken = nil  and serviceAccountName is blank (default: default)

A WARN log is generated when a pod is found using Pod.Spec.DeprecatedServiceAccount
Fix this by updating serviceAccount to serviceAccountName in your .yamls

Example usage:
kubeaudit rbac sat`,
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
				go auditAutomountServiceAccountToken(resource)
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
			go auditAutomountServiceAccountToken(kubeAuditStatefulSets{list: statefulSets})
			go auditAutomountServiceAccountToken(kubeAuditDaemonSets{list: daemonSets})
			go auditAutomountServiceAccountToken(kubeAuditPods{list: pods})
			go auditAutomountServiceAccountToken(kubeAuditReplicationControllers{list: replicationControllers})
			go auditAutomountServiceAccountToken(kubeAuditDeployments{list: deployments})
			wg.Wait()
		}
	},
}

func init() {
	rbacCmd.AddCommand(satCmd)
}

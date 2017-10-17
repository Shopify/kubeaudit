package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkAutomountServiceAccountToken(result *Result) {
	// Check for use of deprecated service account name
	if result.DSA != "" {
		occ := Occurrence{id: ErrorServiceAccountTokenDeprecated, kind: Warn, message: "serviceAccount is a depreciated alias for ServiceAccountName, use that one instead"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if result.Token != nil && *result.Token && result.SA == "" {
		// automountServiceAccountToken = true, and serviceAccountName is blank (default: default)
		occ := Occurrence{id: ErrorServiceAccountTokenTrueAndNoName, kind: Error, message: "Default serviceAccount with token mounted. Please set AutomountServiceAccountToken to false"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if result.Token == nil && result.SA == "" {
		// automountServiceAccountToken = nil (default: true), and serviceAccountName is blank (default: default)
		occ := Occurrence{id: ErrorServiceAccountTokenNILAndNoName, kind: Error, message: "Default serviceAccount with token mounted. Please set AutomountServiceAccountToken to false"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
}

func auditAutomountServiceAccountToken(items Items) (results []Result) {
	for _, item := range items.Iter() {
		result := ServiceAccountIter(item)
		checkAutomountServiceAccountToken(result)

		if result != nil && len(result.Occurrences) > 0 {
			results = append(results, *result)
		}
	}
	for _, result := range results {
		result.Print()
	}
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
			resources, Err := getKubeResources(rootConfig.manifest)
			if Err != nil {
				log.Error(Err)
			}
			count := len(resources)
			wg.Add(count)
			for _, resource := range resources {
				go auditAutomountServiceAccountToken(resource)
			}
			wg.Wait()
		} else {
			kube, Err := kubeClient(rootConfig.kubeConfig)
			if Err != nil {
				log.Error(Err)
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

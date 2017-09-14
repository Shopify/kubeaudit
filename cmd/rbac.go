package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rbacCmd = &cobra.Command{
	Use:   "rbac",
	Short: "Audit RBAC things",
	Long: `Example usage:
kubeaudit rbac sat`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Example usage:\n  kubeaudit rbac sat\n")
	},
}

func init() {
	RootCmd.AddCommand(rbacCmd)
}

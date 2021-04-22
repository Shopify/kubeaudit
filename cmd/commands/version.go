package commands

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed VERSION
var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current kubeaudit version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.TrimSpace(version))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

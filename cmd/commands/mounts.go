package commands

import (
	"bytes"
	"fmt"
	"github.com/Shopify/kubeaudit/auditors/mounts"
	"github.com/spf13/cobra"
	"strings"
)

const sensitivePathsFlagName = "denyPathsList"

var mountsConfig mounts.Config

var mountsCmd = &cobra.Command{
	Use:   "mounts",
	Short: "Audit containers that mount sensitive paths",
	Long: fmt.Sprintf(`This command determines which containers mount sensitive host paths. If no paths list is provided, the following 
paths are used:
%s

A WARN result is generated when a container mounts one or more paths specified with the '--denyPathsList' argument.

Example usage:
kubeaudit mounts --denyPathsList "%s"`, formatPathsList(), strings.Join(mounts.DefaultSensitivePaths[:3], ",")),
	Run: func(cmd *cobra.Command, args []string) {
		runAudit(mounts.New(mountsConfig))(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(mountsCmd)
	setPathsFlags(mountsCmd)
}

func setPathsFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(&mountsConfig.SensitivePaths, sensitivePathsFlagName, "d", mounts.DefaultSensitivePaths,
		"List of sensitive paths that shouldn't be mounted")
}

func formatPathsList() string {
	var buffer bytes.Buffer
	for _, path := range mounts.DefaultSensitivePaths {
		buffer.WriteString("\n- ")
		buffer.WriteString(path)
	}
	return buffer.String()
}

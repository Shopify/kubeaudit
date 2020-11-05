package commands

import (
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/spf13/cobra"
)

var imageConfig image.Config

const imageFlagName = "image"

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Audit containers not using a specified image:tag",
	Long: `This command audits a container against a given image:tag.

An ERROR result is generated when a container does not match the image:tag

An INFO result is generated when a container has a matching image:tag.

This command is also a root command, check 'kubeaudit image --help'.

Example usage:
kubeaudit image --image gcr.io/google_containers/echoserver:1.7
kubeaudit image -i gcr.io/google_containers/echoserver:1.7`,
	Run: func(cmd *cobra.Command, args []string) {
		runAudit(image.New(imageConfig))(cmd, args)
	},
}

func setImageFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&imageConfig.Image, imageFlagName, "i", "", "Image to check against")
}

func init() {
	RootCmd.AddCommand(imageCmd)
	setImageFlags(imageCmd)
}

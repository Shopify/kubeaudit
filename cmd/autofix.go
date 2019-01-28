package cmd

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getAuditFunctions() []interface{} {
	return []interface{}{
		auditAllowPrivilegeEscalation, auditReadOnlyRootFS, auditRunAsNonRoot,
		auditAutomountServiceAccountToken, auditPrivileged, auditCapabilities,
		auditAppArmor, auditSeccomp,
	}
}

// The fix function does not preserve comments (because kubernetes resources do not support comments) so we convert
// both the original manifest file and the fixed manifest file into MapSlices (an array representation of a map which
// preserves the order of the keys) using the Shopify/yaml fork of go-yaml/yaml (the fork adds comment support) and
// then merge the fixed MapSlice back into the original MapSlice so that we get the comments and original order back.
func autofix(*cobra.Command, []string) {
	resources, unsupportedResources, err := getKubeResourcesManifest(rootConfig.manifest)
	if err != nil {
		log.Error(err)
	}

	fixedResources := fix(resources)

	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(tmpFile.Name())

	err = writeManifestFile(fixedResources, unsupportedResources, tmpFile.Name())
	if err != nil {
		log.Error(err)
	}

	fixedYaml, err := mergeYAML(rootConfig.manifest, tmpFile.Name())
	if err != nil {
		log.Error(err)
	}

	err = ioutil.WriteFile(rootConfig.manifest, fixedYaml, 0644)
	if err != nil {
		log.Error(err)
	}
}x

var autofixCmd = &cobra.Command{
	Use:   "autofix",
	Short: "Automagically fixes a manifest to be secure",
	Long: `"autofix" will examine a manifest file and automagically fill in the blanks to leave your yaml file more secure than it found it

Example usage:
kubeaudit autofix -f /path/to/yaml`,
	Run: autofix,
}

func init() {
	RootCmd.AddCommand(autofixCmd)
}

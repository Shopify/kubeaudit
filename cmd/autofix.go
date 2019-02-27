package cmd

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

// The fix function does not preserve comments (because kubernetes resources do not support comments) so we convert
// both the original manifest file and the fixed manifest file into MapSlices (an array representation of a map which
// preserves the order of the keys) using the Shopify/yaml fork of go-yaml/yaml (the fork adds comment support) and
// then merge the fixed MapSlice back into the original MapSlice so that we get the comments and original order back.
func autofix(*cobra.Command, []string) {

	var toAppend = false

	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	if err != nil {
		log.Error(err)
		return
	}

	fixedResources, extraResources := fix(resources)

	tmpFixedFile, err := ioutil.TempFile("", "kubeaudit_autofix_fixed")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(tmpFixedFile.Name())
	tmpOrigFile, err := ioutil.TempFile("", "kubeaudit_autofix_orig")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(tmpOrigFile.Name())
	finalFile, err := ioutil.TempFile("", "kubeaudit_autofix_final")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(finalFile.Name())

	splitResources, toAppend, err := splitYamlResources(rootConfig.manifest, finalFile.Name())
	if err != nil {
		log.Error(err)
	}

	for index := range fixedResources {
		err = WriteToFile(fixedResources[index], tmpFixedFile.Name())
		if err != nil {
			log.Error(err)
		}
		err := ioutil.WriteFile(tmpOrigFile.Name(), splitResources[index], 0644)
		if err != nil {
			log.Error(err)
		}
		fixedYaml, err := mergeYAML(tmpOrigFile.Name(), tmpFixedFile.Name())
		if err != nil {
			log.Error(err)
		}
		err = writeManifestFile(fixedYaml, finalFile.Name(), toAppend)
		if err != nil {
			log.Error(err)
		}
		toAppend = true
	}

	for index := range extraResources {
		info, _ := k8sRuntime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), "application/yaml")
		groupVersion := schema.GroupVersion{Group: extraResources[index].GetObjectKind().GroupVersionKind().Group, Version: extraResources[index].GetObjectKind().GroupVersionKind().Version}
		encoder := scheme.Codecs.EncoderForVersion(info.Serializer, groupVersion)
		fixedData, err := k8sRuntime.Encode(encoder, extraResources[index])
		if err != nil {
			log.Error(err)
		}
		err = writeManifestFile(fixedData, finalFile.Name(), toAppend)
		if err != nil {
			log.Error(err)
		}
		toAppend = true
	}

	finalData, err := ioutil.ReadFile(finalFile.Name())
	if err != nil {
		log.Error(err)
	}
	err = os.Truncate(rootConfig.manifest, 0)
	if err != nil {
		log.Error(err)
	}
	err = writeManifestFile(finalData, rootConfig.manifest, !isFirstLineSeparatorOrComment(finalFile.Name()))
	if err != nil {
		log.Error(err)
	}
}

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

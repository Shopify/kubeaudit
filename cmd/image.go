package cmd

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var imgConfig imgFlags

type imgFlags struct {
	img  string
	name string
	tag  string
}

func (image *imgFlags) splitImageString() {
	tokens := strings.Split(image.img, ":")
	if len(tokens) > 0 {
		image.name = tokens[0]
	}
	if len(tokens) > 1 {
		image.tag = tokens[1]
	}
}

func checkImage(container Container, image imgFlags, result *Result) {
	image.splitImageString()
	contImage := imgFlags{img: container.Image}
	contImage.splitImageString()
	result.ImageName = contImage.name
	result.ImageTag = contImage.tag

	if len(contImage.tag) == 0 {
		occ := Occurrence{id: ErrorImageTagMissing, kind: Warn, message: "Image tag was missing"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if contImage.name == image.name && contImage.tag != image.tag {
		occ := Occurrence{id: ErrorImageTagIncorrect, kind: Error, message: "Image tag was incorrect"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if contImage.name == image.name && contImage.tag == image.tag {
		occ := Occurrence{id: InfoImageCorrect, kind: Info, message: "Image tag was correct"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
}

func auditImages(image imgFlags, items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkImage(container, image, result)
			if result != nil && len(result.Occurrences) > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	for _, result := range results {
		result.Print()
	}
	return
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Audit container images",
	Long: `This command audits a container against a given image:tag
An INFO log is given when a container has a matching image:tag
An ERROR log is generated when a container does not match the image:tag
This command is also a root command, check kubeaudit sc --help
Example usage:
kubeaudit image --image gcr.io/google_containers/echoserver:1.7
kubeaudit image -i gcr.io/google_containers/echoserver:1.7`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(imgConfig.img) == 0 {
			log.Error("Empty image name. Are you missing the image flag?")
			return
		}
		imgConfig.splitImageString()
		if len(imgConfig.tag) == 0 {
			log.Error("Empty image tag. Are you missing the image tag?")
			return
		}

		if rootConfig.json {
			log.SetFormatter(&log.JSONFormatter{})
		}
		var resources []Items

		if rootConfig.manifest != "" {
			var err error
			resources, err = getKubeResourcesManifest(rootConfig.manifest)
			if err != nil {
				log.Error(err)
			}
		} else {
			kube, err := kubeClient(rootConfig.kubeConfig)
			if err != nil {
				log.Error(err)
			}
			resources = getKubeResources(kube)
		}

		var wg sync.WaitGroup
		wg.Add(len(resources))

		for _, resource := range resources {
			go func(items Items) {
				auditImages(imgConfig, items)
				wg.Done()
			}(resource)
		}

		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(imageCmd)
	imageCmd.Flags().StringVarP(&imgConfig.img, "image", "i", "", "image to check against")
}

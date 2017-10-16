package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
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

func printResultImg(results []Result) {
	for _, result := range results {
		switch result.err {
		case KubeAuditInfo:
			log.WithFields(log.Fields{
				"type":      result.kubeType,
				"tag":       result.imgTag,
				"namespace": result.namespace,
				"name":      result.name,
			}).Info(result.namespace, "/", result.name)
		case ErrorImageTagMissing:
			log.WithFields(log.Fields{
				"type":      result.kubeType,
				"tag":       result.imgTag,
				"namespace": result.namespace,
				"name":      result.name,
			}).Error("Image tag was missing")
		case ErrorImageTagIncorrect:
			log.WithFields(log.Fields{
				"type":      result.kubeType,
				"tag":       result.imgTag,
				"namespace": result.namespace,
				"name":      result.name,
			}).Error("Image tag was incorrect")
		}
	}
}

func checkImage(container apiv1.Container, image imgFlags, result *Result) {
	image.splitImageString()
	contImage := imgFlags{img: container.Image}
	contImage.splitImageString()

	if len(contImage.tag) == 0 {
		if image.name == contImage.name {
			// Image name was proper but image tag was missing
			result.err = ErrorImageTagMissing
		}
		return
	}

	if contImage.name == image.name && contImage.tag != image.tag {
		result.err = ErrorImageTagIncorrect
		result.imgName = contImage.name
		result.imgTag = contImage.tag
	}
}

func auditImages(image imgFlags, items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkImage(container, image, result)
			if result != nil && result.err > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	printResultImg(results)
	defer wg.Done()
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

		if rootConfig.manifest != "" {
			resources, err := getKubeResources(rootConfig.manifest)
			if err != nil {
				log.Error(err)
			}
			count := len(resources)
			wg.Add(count)
			for _, resource := range resources {
				go auditImages(imgConfig, resource)
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
			replicationControllers := getReplicationControllers(kube)
			pods := getPods(kube)

			wg.Add(5)
			go auditImages(imgConfig, kubeAuditStatefulSets{list: statefulSets})
			go auditImages(imgConfig, kubeAuditDaemonSets{list: daemonSets})
			go auditImages(imgConfig, kubeAuditPods{list: pods})
			go auditImages(imgConfig, kubeAuditReplicationControllers{list: replicationControllers})
			go auditImages(imgConfig, kubeAuditDeployments{list: deployments})
			wg.Wait()
		}
	},
}

func init() {
	RootCmd.AddCommand(imageCmd)
	imageCmd.Flags().StringVarP(&imgConfig.img, "image", "i", "", "image to check against")
}

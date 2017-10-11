package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
)

var imgConfig imgFlags

type imgFlags struct {
	img string
}

func printResultImg(results []Result) {
	for _, result := range results {
		switch result.err {
		case 0:
			log.WithFields(log.Fields{
				"type": result.kubeType,
				"tag":  result.img_tag}).Info(result.namespace,
				"/", result.name)
		case 1:
			log.WithFields(log.Fields{
				"type": result.kubeType,
				"tag":  result.img_tag}).Error("Image tag was missing ", result.namespace,
				"/", result.name)
		case 2:
			log.WithFields(log.Fields{
				"type": result.kubeType,
				"tag":  result.img_tag}).Error("Image tag was incorrect ", result.namespace,
				"/", result.name)
		}
	}
}

func checkImage(container apiv1.Container, image string, result *Result) {
	var (
		contImg string
		contTag string
	)

	imageLabel := strings.Split(image, ":")

	if len(imageLabel) < 2 {
		log.Error("Image tag is missing!")
		os.Exit(1)
	}

	compImg := imageLabel[0]
	compTag := imageLabel[1]

	contImgLabel := strings.Split(container.Image, ":")
	contImg = contImgLabel[0]

	if len(contImgLabel) < 2 {
		if compImg == contImg {
			// Image name was proper but image tag was missing
			result.err = 1
		}
		return
	}

	contTag = contImgLabel[1]

	if contImg == compImg && contTag != compTag {
		result.err = 2
		result.img_name = contImg
		result.img_tag = contTag
	}
	return
}

func auditImages(image string, items Items) (results []Result) {
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
				go auditSecurityContext(resource)
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
			go auditImages(imgConfig.img, kubeAuditStatefulSets{list: statefulSets})
			go auditImages(imgConfig.img, kubeAuditDaemonSets{list: daemonSets})
			go auditImages(imgConfig.img, kubeAuditPods{list: pods})
			go auditImages(imgConfig.img, kubeAuditReplicationControllers{list: replicationControllers})
			go auditImages(imgConfig.img, kubeAuditDeployments{list: deployments})
			wg.Wait()
		}
	},
}

func init() {
	RootCmd.AddCommand(imageCmd)
	imageCmd.Flags().StringVarP(&imgConfig.img, "image", "i", "", "image to check against")
}

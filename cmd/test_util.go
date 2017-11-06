package cmd

import (
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var path = "../fixtures/"

func runTest(t *testing.T, file string, function interface{}, errCode int, imageStr ...string) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	var image imgFlags
	switch function.(type) {
	case (func(imgFlags, Items) []Result):
		if len(imageStr) != 1 {
			log.Fatal("Incorrect number of images specified")
		}
		image = imgFlags{img: imageStr[0]}
		image.splitImageString()
	}

	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	var results []Result
	for _, resource := range resources {
		var currentResults []Result
		switch f := function.(type) {
		case (func(Items) []Result):
			currentResults = f(resource)
		case (func(imgFlags, Items) []Result):
			currentResults = f(image, resource)
		default:
			log.Fatal("Invalid function provided")
		}
		for _, currentResult := range currentResults {
			results = append(results, currentResult)
		}
	}
	var errors []int
	for _, result := range results {
		for _, occurrence := range result.Occurrences {
			errors = append(errors, occurrence.id)
		}
	}
	assert.Contains(errors, errCode)
}

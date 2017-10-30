package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var path = "../fixtures/"

func runTest(t *testing.T, file string, function func(Items) []Result, errCode int) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	var results []Result
	for _, resource := range resources {
		current_results := function(resource)
		for _, current_result := range current_results {
			results = append(results, current_result)
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

func runImageTest(t *testing.T, file string, function func(imgFlags, Items) []Result, image_str string, errCode int) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	image := imgFlags{img: image_str}
	image.splitImageString()
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	var results []Result
	for _, resource := range resources {
		current_results := function(image, resource)
		for _, current_result := range current_results {
			results = append(results, current_result)
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

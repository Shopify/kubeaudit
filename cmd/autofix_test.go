package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
)

func TestFixV1(t *testing.T) {
	file := "../fixtures/autofix_v1.yml"
	fileFixed := "../fixtures/autofix-fixed_v1.yml"
	rootConfig.manifest = file
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources, _ := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assertEqualWorkloads(assert, correctlyFixedResources, fixedResources)
}

func TestAllResourcesFixV1(t *testing.T) {
	file := "../fixtures/autofix-all-resources_v1.yml"
	fileFixedResources := "../fixtures/autofix-fixed_v1.yml"
	fileExtraResources := "../fixtures/autofix-extra-resources-fixed_v1.yml"
	rootConfig.manifest = file
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources, extraResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixedResources)
	assert.Nil(err)
	correctlyFixedExtraResources, err := getKubeResourcesManifest(fileExtraResources)
	assertEqualWorkloads(assert, correctlyFixedResources, fixedResources)
	assertEqualWorkloads(assert, correctlyFixedExtraResources, extraResources)
}

func TestExtraResourcesFixV1(t *testing.T) {
	file := "../fixtures/autofix-extra-resources_v1.yml"
	fileFixed := "../fixtures/autofix-extra-resources-fixed_v1.yml"
	rootConfig.manifest = file
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	_, extraResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assertEqualWorkloads(assert, correctlyFixedResources, extraResources)
}

func TestExtraResourcesEgressFixV1(t *testing.T) {
	file := "../fixtures/autofix-extra-resources-egress_v1.yml"
	fileFixed := "../fixtures/autofix-extra-resources-egress-fixed_v1.yml"
	rootConfig.manifest = file
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	_, extraResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assertEqualWorkloads(assert, correctlyFixedResources, extraResources)

}

func TestExtraResourcesIngressFixV1(t *testing.T) {
	file := "../fixtures/autofix-extra-resources-ingress_v1.yml"
	fileFixed := "../fixtures/autofix-extra-resources-ingress-fixed_v1.yml"
	rootConfig.manifest = file
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	_, extraResources := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assertEqualWorkloads(assert, correctlyFixedResources, extraResources)
}

func TestFixV1Beta1(t *testing.T) {
	file := "../fixtures/autofix_v1beta1.yml"
	fileFixed := "../fixtures/autofix-fixed_v1beta1.yml"
	assert := assert.New(t)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	fixedResources, _ := fix(resources)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assertEqualWorkloads(assert, correctlyFixedResources, fixedResources)
}

func TestFixV1Beta2(t *testing.T) {
	origFilename := "../fixtures/autofix-all-resources_v1.yml"
	expectedFilename := "../fixtures/autofix-all-resources-fixed_v1.yml"
	assert := assert.New(t)

	// Copy original yaml to a temp file because autofix modifies the input file
	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix_test")
	tmpFilename := tmpFile.Name()
	assert.Nil(err)
	defer os.Remove(tmpFilename)
	origFile, err := os.Open(origFilename)
	assert.Nil(err)
	_, err = io.Copy(tmpFile, origFile)
	assert.Nil(err)
	tmpFile.Close()
	origFile.Close()

	rootConfig.manifest = tmpFilename
	autofix(nil, nil)

	assert.True(compareTextFiles(expectedFilename, tmpFilename))
}

func TestInvalidManifests(t *testing.T) {
	file := "../configs/capSetConfig.yaml"
	assert := assert.New(t)
	_, err := getKubeResourcesManifest(file)
	assert.NotNil(err)
}

func TestPreserveComments(t *testing.T) {
	origFilename := "../fixtures/autofix_v1.yml"
	expectedFilename := "../fixtures/autofix-fixed_v1.yml"
	assert := assert.New(t)

	// Copy original yaml to a temp file because autofix modifies the input file
	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix_test")
	tmpFilename := tmpFile.Name()
	assert.Nil(err)
	defer os.Remove(tmpFilename)
	origFile, err := os.Open(origFilename)
	assert.Nil(err)
	_, err = io.Copy(tmpFile, origFile)
	assert.Nil(err)
	tmpFile.Close()
	origFile.Close()

	rootConfig.manifest = tmpFilename
	autofix(nil, nil)

	assert.True(compareTextFiles(expectedFilename, tmpFilename))
}

func TestPreserveCommentsV2(t *testing.T) {
	origFilename := "../fixtures/preserve_comments_v2.yml"
	expectedFilename := "../fixtures/preserve_comments-fixed_v2.yml"
	assert := assert.New(t)
	// Copy original yaml to a temp file because autofix modifies the input file
	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix_test")
	tmpFilename := tmpFile.Name()
	assert.Nil(err)
	defer os.Remove(tmpFilename)
	origFile, err := os.Open(origFilename)
	assert.Nil(err)
	_, err = io.Copy(tmpFile, origFile)
	assert.Nil(err)
	tmpFile.Close()
	origFile.Close()

	rootConfig.manifest = tmpFilename
	autofix(nil, nil)

	assert.True(compareTextFiles(expectedFilename, tmpFilename))
}

func TestUnsupportedResources(t *testing.T) {
	origFilename := "../fixtures/autofix-unsupported_v1.yml"
	expectedFilename := "../fixtures/autofix-unsupported-fixed_v1.yml"
	assert := assert.New(t)

	// Copy original yaml to a temp file because autofix modifies the input file
	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix_test")
	tmpFilename := tmpFile.Name()
	assert.Nil(err)
	defer os.Remove(tmpFilename)
	origFile, err := os.Open(origFilename)
	assert.Nil(err)
	_, err = io.Copy(tmpFile, origFile)
	assert.Nil(err)
	tmpFile.Close()
	origFile.Close()

	rootConfig.manifest = tmpFilename
	autofix(nil, nil)

	assert.True(compareTextFiles(expectedFilename, tmpFilename))
}

func TestSingleUnsupportedResource(t *testing.T) {

	origFilename := "../fixtures/autofix-single-unsupported_v1.yml"
	expectedFilename := "../fixtures/autofix-single-unsupported-fixed_v1.yml"
	assert := assert.New(t)

	// Copy original yaml to a temp file because autofix modifies the input file
	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix_test")
	tmpFilename := tmpFile.Name()
	assert.Nil(err)
	defer os.Remove(tmpFilename)
	origFile, err := os.Open(origFilename)
	assert.Nil(err)
	_, err = io.Copy(tmpFile, origFile)
	assert.Nil(err)
	tmpFile.Close()
	origFile.Close()

	rootConfig.manifest = tmpFilename
	autofix(nil, nil)

	assert.True(compareTextFiles(expectedFilename, tmpFilename))

}

var testData = []struct {
	orig   *yaml.Node
	fixed  *yaml.Node
	merged *yaml.Node
}{
	// update
	{
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "K", HeadComment: "Hi"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "V", LineComment: "Bye"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "K"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "V2"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "K", HeadComment: "Hi"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "V2", LineComment: "Bye"},
		}},
	},
	// add
	{
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "Comment"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "Comment"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
		}},
	},
	// remove
	{
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "Comment 1"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment 2"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment 2"},
		}},
	},
	// preserve order
	{
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k", HeadComment: "Comment 1"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "EOL Comment 1"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2", HeadComment: "Comment 2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "EOL Comment 2", FootComment: "Foot Comment"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "knew"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "vnew"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2new"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "vnew"},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k", HeadComment: "Comment 1"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "vnew", LineComment: "EOL Comment 1"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2", HeadComment: "Comment 2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2new", LineComment: "EOL Comment 2", FootComment: "Foot Comment"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "knew"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "vnew"},
		}},
	},
	// sequence of strings
	{
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "values"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v1", LineComment: "EOL Comment 1"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", HeadComment: "Comment 1", LineComment: "EOL Comment 2"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v3", HeadComment: "EOL Comment 3"},
			}},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "values"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v3"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v4"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v1"},
			}},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "values"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v1", LineComment: "EOL Comment 1"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v3", HeadComment: "EOL Comment 3"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "v4"},
			}},
		}},
	},
	// complex mapslice (nested sequences and mapslices)
	{
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "containers"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "container1", LineComment: "Container 1"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "image"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "image1", LineComment: "Image 1"},
				}},
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "container2", LineComment: "Container 2"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "ports"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port 1"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "protocol"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "TCP"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "containerPort"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "3000", LineComment: "Port 3000"},
						}},
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port 2"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "protocol", HeadComment: "Comment"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "UDP"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "containerPort"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "6000", LineComment: "Port 6000"},
						}},
					}},
				}},
			}},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "subsets"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "addresses"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "hostname"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "hname", LineComment: "Comment 1"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "ip"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "10.0.0.1"},
						}},
					}},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "ports"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "pname", LineComment: "Comment 2"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "3000"},
						}},
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "pname2", LineComment: "Comment 3"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "8000"},
						}},
					}},
				}},
			}},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "containers"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "container3"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "image"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "image1"},
				}},
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "container2"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "ports"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port 1"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "containerPort"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "6000"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "protocol"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "TCP"},
						}},
					}},
				}},
			}},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "subsets"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "addresses"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "hostname"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "hname"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "ip"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "10.0.0.1"},
						}},
					}},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "ports"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "pname"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "6000"},
						}},
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "newname"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "8000"},
						}},
					}},
				}},
			}},
		}},
		&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "containers"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "container2", LineComment: "Container 2"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "ports"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port 1"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "protocol", HeadComment: "Comment"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "TCP"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "containerPort"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "6000", LineComment: "Port 6000"},
						}},
					}},
				}},
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "container3"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "image"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "image1"},
				}},
			}},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "subsets"},
			{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "addresses"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "hostname"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "hname", LineComment: "Comment 1"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "ip"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "10.0.0.1"},
						}},
					}},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "ports"},
					{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "newname", LineComment: "Comment 3"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "8000"},
						}},
						{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "pname"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "port"},
							{Kind: yaml.ScalarNode, Tag: strTag, Value: "6000"},
						}},
					}},
				}},
			}},
		}},
	},
}

func TestMergeYAML(t *testing.T) {
	assert := assert.New(t)

	for _, test := range testData {
		merged := mergeMaps(test.orig, test.fixed)
		assert.True(deepEqual(test.merged, merged))
		assert.Equal(test.merged, merged)
	}
}

func TestFindItemInSequence(t *testing.T) {
	assert := assert.New(t)

	// same string
	seq := &yaml.Node{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", HeadComment: "Comment"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment 2"},
	}}
	item := &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"}
	index := findItemInSequence("", item, seq)
	assert.Equal(1, index)

	// different string
	seq = &yaml.Node{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", HeadComment: "Comment"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment 2"},
	}}
	item = &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v3"}
	index = findItemInSequence("", item, seq)
	assert.Equal(-1, index)

	// matching mapslice
	seq = &yaml.Node{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
		{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "port1"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "protocol"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "TCP"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "containerPort"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "3000"},
		}},
		{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "port2"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "protocol"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "UDP", HeadComment: "Comment"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "containerPort"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "6000", LineComment: "Comment 2"},
		}},
	}}
	item = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "port1"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "protocol"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "TCP"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "containerPort"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "6000"},
	}}
	index = findItemInSequence("ports", item, seq)
	assert.Equal(1, index)
}

func TestFindItemInMapSlice(t *testing.T) {
	assert := assert.New(t)

	// key present
	m := &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k", HeadComment: "Comment"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "Comment 2"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
	}}
	item, index := findValInMap("k2", m)
	assert.Equal(3, index)
	assert.True(deepEqual(m.Content[3], item))

	// key not present
	m = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k", HeadComment: "Comment"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "Comment 2"},
	}}
	item, index = findValInMap("k2", m)
	assert.Equal(-1, index)
}

func TestEqualValueForKey(t *testing.T) {
	assert := assert.New(t)

	// same string
	m1 := &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment 2"},
	}}
	m2 := &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k2"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
	}}
	assert.True(equalValueForKey("k2", m1, m2))

	// different string
	m1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v"},
	}}
	m2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
	}}
	assert.False(equalValueForKey("k2", m1, m2))

	// string and number
	m1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "2"},
	}}
	m2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
		{Kind: yaml.ScalarNode, Tag: "!!int", Value: "2"},
	}}
	assert.False(equalValueForKey("k", m1, m2))
}

func TestSequenceItemMatch(t *testing.T) {
	assert := assert.New(t)

	// same strings
	item1 := &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment"}
	item2 := &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"}
	assert.True(sequenceItemMatch("", item1, item2))

	// different strings
	item1 = &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "Comment"}
	item2 = &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment"}
	assert.False(sequenceItemMatch("", item1, item2))

	// string and mapslice
	item1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, LineComment: "Comment"},
	}, LineComment: "Comment"}
	item2 = &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"}
	assert.False(sequenceItemMatch("", item1, item2))

	// mapslices with same value for identifying key
	item1 = &yaml.Node{Kind: yaml.MappingNode, LineComment: "Comment", Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "containername", LineComment: "line"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image1"},
	}}
	item2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image2"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "containername"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "imagePullPolicy"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "IfNotPresent"},
	}}
	assert.True(sequenceItemMatch("containers", item1, item2))

	// mapslices with different value for identifying key
	item1 = &yaml.Node{Kind: yaml.MappingNode, LineComment: "Comment", Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "containername"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image1"},
	}}
	item2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "containername2"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "image1"},
	}}
	assert.False(sequenceItemMatch("containers", item1, item2))

	// Container.envFrom
	item1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "configMapRef"},
		{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "n"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "optional"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "true"},
		}},
	}}
	item2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "configMapRef"},
		{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "n"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "optional"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "true"},
		}},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "prefix"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "pre"},
	}}
	assert.True(sequenceItemMatch("envFrom", item1, item2))
	item1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "configMapRef"},
		{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "n"},
		}},
	}}
	item2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "secretRef"},
		{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "name"},
			{Kind: yaml.ScalarNode, Tag: strTag, Value: "n"},
		}},
	}}
	assert.False(sequenceItemMatch("envFrom", item1, item2))
}

func TestDeepEqual(t *testing.T) {
	assert := assert.New(t)

	var v1 *yaml.Node
	var v2 *yaml.Node

	// maps should be equal, regardless of comments
	v1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k", HeadComment: "head", LineComment: "line", FootComment: "foot"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", HeadComment: "head", LineComment: "line", FootComment: "foot"},
	}}
	v2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "k"},
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v"},
	}}
	assert.True(deepEqual(v1, v2))

	// sequences should be equal, regardless of comments
	v1 = &yaml.Node{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", HeadComment: "head", LineComment: "line", FootComment: "foot"},
	}}
	v2 = &yaml.Node{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "v"},
	}}
	assert.True(deepEqual(v1, v2))

	v1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "matchExpressions"},
		{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
			{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "key"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "labelkey"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "operator"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "In"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "values"},
				{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value1", LineComment: "Comment 1"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value2"},
				}},
			}},
		}},
	}}
	v2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "matchExpressions"},
		{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
			{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "key"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "labelkey"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "values"},
				{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value1"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value2"},
				}},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "operator"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "In"},
			}},
		}},
	}}
	assert.True(deepEqual(v1, v2))

	v1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "matchExpressions"},
		{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
			{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "key"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "labelkey"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "operator"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "In"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "values"},
				{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value1", LineComment: "Comment 1"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value2"},
				}},
			}},
		}},
	}}
	v1 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: strTag, Value: "matchExpressions"},
		{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
			{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "key"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "labelkey"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "operator"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "In"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "values"},
				{Kind: yaml.SequenceNode, Tag: seqTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value1", LineComment: "Comment 1"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "newvalue"},
				}},
			}},
		}},
	}}
	assert.False(deepEqual(v1, v2))
}

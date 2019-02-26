package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Shopify/yaml"
	"github.com/stretchr/testify/assert"
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
	orig   yaml.MapSlice
	fixed  yaml.MapSlice
	merged yaml.MapSlice
}{
	// update
	{
		yaml.MapSlice{{Key: "k", Value: "v", Comment: "Comment"}},
		yaml.MapSlice{{Key: "k", Value: "v2"}},
		yaml.MapSlice{{Key: "k", Value: "v2", Comment: "Comment"}},
	},
	// add
	{
		yaml.MapSlice{{Key: "k", Value: "v", Comment: "Comment"}},
		yaml.MapSlice{{Key: "k", Value: "v"}, {Key: "k2", Value: "v2"}},
		yaml.MapSlice{{Key: "k", Value: "v", Comment: "Comment"}, {Key: "k2", Value: "v2"}},
	},
	// remove
	{
		yaml.MapSlice{{Key: "k", Value: "v", Comment: "Comment 1"}, {Key: "k2", Value: "v2", Comment: "Comment 2"}},
		yaml.MapSlice{{Key: "k2", Value: "v2"}},
		yaml.MapSlice{{Key: "k2", Value: "v2", Comment: "Comment 2"}},
	},
	// preserve order
	{
		yaml.MapSlice{
			{Comment: "Comment 1"},
			{Comment: "Comment 2"},
			{Key: "k", Value: "v", Comment: "EOL Comment 1"},
			{Comment: "Comment 3"},
			{Key: "k2", Value: "v2", Comment: "EOL Comment 2"},
			{Comment: "Comment 4"},
		},
		yaml.MapSlice{
			{Key: "knew", Value: "vnew"},
			{Key: "k2", Value: "v2new"},
			{Key: "k", Value: "vnew"},
		},
		yaml.MapSlice{
			{Comment: "Comment 1"},
			{Comment: "Comment 2"},
			{Key: "k", Value: "vnew", Comment: "EOL Comment 1"},
			{Comment: "Comment 3"},
			{Key: "k2", Value: "v2new", Comment: "EOL Comment 2"},
			{Comment: "Comment 4"},
			{Key: "knew", Value: "vnew"},
		},
	},
	// sequence of strings
	{
		yaml.MapSlice{{Key: "values", Value: []yaml.SequenceItem{
			{Value: "v1", Comment: "EOL Comment 1"},
			{Comment: "Comment 1"},
			{Value: "v2", Comment: "EOL Comment 2"},
			{Value: "v3", Comment: "EOL Comment 3"},
		}}},
		yaml.MapSlice{{Key: "values", Value: []yaml.SequenceItem{
			{Value: "v3"},
			{Value: "v4"},
			{Value: "v1"},
		}}},
		yaml.MapSlice{{Key: "values", Value: []yaml.SequenceItem{
			{Value: "v1", Comment: "EOL Comment 1"},
			{Comment: "Comment 1"},
			{Value: "v3", Comment: "EOL Comment 3"},
			{Value: "v4"},
		}}},
	},
	// complex mapslice (nested sequences and mapslices)
	{
		yaml.MapSlice{
			{Key: "containers", Value: []yaml.SequenceItem{
				{Value: yaml.MapSlice{
					{Key: "name", Value: "container1", Comment: "Container 1"},
					{Key: "image", Value: "image1", Comment: "Image 1"},
				}},
				{Value: yaml.MapSlice{
					{Key: "name", Value: "container2", Comment: "Container 2"},
					{Key: "ports", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "name", Value: "port1"},
							{Key: "protocol", Value: "TCP"},
							{Key: "containerPort", Value: "3000", Comment: "Port 3000"},
						}},
						{Value: yaml.MapSlice{
							{Key: "name", Value: "port2"},
							{Comment: "Comment"},
							{Key: "protocol", Value: "UDP"},
							{Key: "containerPort", Value: "6000", Comment: "Port 6000"},
						}},
					}},
				}},
			}},
			{Key: "subsets", Value: []yaml.SequenceItem{
				{Value: yaml.MapSlice{
					{Key: "addresses", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "hostname", Value: "hname", Comment: "Comment 1"},
							{Key: "ip", Value: "10.0.0.1"},
						}},
					}},
					{Key: "ports", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "name", Value: "pname", Comment: "Comment 2"},
							{Key: "port", Value: "3000"},
						}},
						{Value: yaml.MapSlice{
							{Key: "name", Value: "pname2", Comment: "Comment 3"},
							{Key: "port", Value: "8000"},
						}},
					}},
				}},
			}},
		},
		yaml.MapSlice{
			{Key: "containers", Value: []yaml.SequenceItem{
				{Value: yaml.MapSlice{
					{Key: "name", Value: "container3"},
					{Key: "image", Value: "image1"},
				}},
				{Value: yaml.MapSlice{
					{Key: "name", Value: "container2"},
					{Key: "ports", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "name", Value: "port1"},
							{Key: "containerPort", Value: "6000"},
							{Key: "protocol", Value: "TCP"},
						}},
					}},
				}},
			}},
			{Key: "subsets", Value: []yaml.SequenceItem{
				{Value: yaml.MapSlice{
					{Key: "addresses", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "hostname", Value: "hname"},
							{Key: "ip", Value: "10.0.0.1"},
						}},
					}},
					{Key: "ports", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "name", Value: "pname"},
							{Key: "port", Value: "6000"},
						}},
						{Value: yaml.MapSlice{
							{Key: "name", Value: "newname"},
							{Key: "port", Value: "8000"},
						}},
					}},
				}},
			}},
		},
		yaml.MapSlice{
			{Key: "containers", Value: []yaml.SequenceItem{
				{Value: yaml.MapSlice{
					{Key: "name", Value: "container2", Comment: "Container 2"},
					{Key: "ports", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "name", Value: "port1"},
							{Comment: "Comment"},
							{Key: "protocol", Value: "TCP"},
							{Key: "containerPort", Value: "6000", Comment: "Port 6000"},
						}},
					}},
				}},
				{Value: yaml.MapSlice{
					{Key: "name", Value: "container3"},
					{Key: "image", Value: "image1"},
				}},
			}},
			{Key: "subsets", Value: []yaml.SequenceItem{
				{Value: yaml.MapSlice{
					{Key: "addresses", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "hostname", Value: "hname", Comment: "Comment 1"},
							{Key: "ip", Value: "10.0.0.1"},
						}},
					}},
					{Key: "ports", Value: []yaml.SequenceItem{
						{Value: yaml.MapSlice{
							{Key: "name", Value: "newname", Comment: "Comment 3"},
							{Key: "port", Value: "8000"},
						}},
						{Value: yaml.MapSlice{
							{Key: "name", Value: "pname"},
							{Key: "port", Value: "6000"},
						}},
					}},
				}},
			}},
		},
	},
}

func TestMergeYAML(t *testing.T) {
	assert := assert.New(t)

	for _, test := range testData {
		assert.Equal(test.merged, mergeMapSlices(test.orig, test.fixed))
	}
}

func TestFindItemInSequence(t *testing.T) {
	assert := assert.New(t)

	// same string
	seq := []yaml.SequenceItem{
		{Comment: "Comment"},
		{Value: "v"},
		{Value: "v2", Comment: "Comment 2"},
	}
	item := yaml.SequenceItem{Value: "v2"}
	seqItem, index := findItemInSequence("", item, seq)
	assert.Equal(2, index)
	assert.True(deepEqual(seq[2], seqItem))

	// different string
	seq = []yaml.SequenceItem{
		{Comment: "Comment"},
		{Value: "v"},
		{Value: "v2", Comment: "Comment 2"},
	}
	item = yaml.SequenceItem{Value: "v3"}
	_, index = findItemInSequence("", item, seq)
	assert.Equal(-1, index)

	// matching mapslice
	seq = []yaml.SequenceItem{
		{Value: yaml.MapSlice{
			{Key: "name", Value: "port1"},
			{Key: "protocol", Value: "TCP"},
			{Key: "containerPort", Value: "3000"},
		}},
		{Value: yaml.MapSlice{
			{Key: "name", Value: "port2"},
			{Comment: "Comment"},
			{Key: "protocol", Value: "UDP"},
			{Key: "containerPort", Value: "6000", Comment: "Comment 2"},
		}},
	}
	item = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "name", Value: "port1"},
		{Key: "protocol", Value: "TCP"},
		{Key: "containerPort", Value: "6000"},
	}}
	seqItem, index = findItemInSequence("ports", item, seq)
	assert.Equal(1, index)
	assert.True(deepEqual(seq[1], seqItem))
}

func TestFindItemInMapSlice(t *testing.T) {
	assert := assert.New(t)

	// key present
	m := yaml.MapSlice{
		{Comment: "Comment"},
		{Key: "k", Value: "v", Comment: "Comment 2"},
		{Key: "k2", Value: "v2"},
	}
	item, index := findItemInMapSlice("k2", m)
	assert.Equal(2, index)
	assert.True(deepEqual(m[2], item))

	// key not present
	m = yaml.MapSlice{
		{Comment: "Comment"},
	}
	item, index = findItemInMapSlice("k2", m)
	assert.Equal(-1, index)
}

func TestMapPairMatch(t *testing.T) {
	assert := assert.New(t)

	// same string
	m1 := yaml.MapSlice{{Key: "k2", Value: "v2", Comment: "Comment 2"}}
	m2 := yaml.MapSlice{
		{Key: "k", Value: "v"},
		{Key: "k2", Value: "v2"},
	}
	assert.True(mapPairMatch("k2", m1, m2))

	// different string
	m1 = yaml.MapSlice{{Key: "k", Value: "v"}}
	m2 = yaml.MapSlice{{Key: "k", Value: "v2"}}
	assert.False(mapPairMatch("k2", m1, m2))

	// string and number
	m1 = yaml.MapSlice{{Key: "k", Value: "2"}}
	m2 = yaml.MapSlice{{Key: "k", Value: 2}}
	assert.False(mapPairMatch("k", m1, m2))
}

func TestSequenceItemMatch(t *testing.T) {
	assert := assert.New(t)

	// same strings
	item1 := yaml.SequenceItem{Value: "v2", Comment: "Comment"}
	item2 := yaml.SequenceItem{Value: "v2"}
	assert.True(sequenceItemMatch("", item1, item2))

	// different strings
	item1 = yaml.SequenceItem{Value: "v", Comment: "Comment"}
	item2 = yaml.SequenceItem{Value: "v2", Comment: "Comment"}
	assert.False(sequenceItemMatch("", item1, item2))

	// string and mapslice
	item1 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Comment: "Comment"},
	}, Comment: "Comment"}
	item2 = yaml.SequenceItem{Value: "v2"}
	assert.False(sequenceItemMatch("", item1, item2))

	// mapslices with same value for identifying key
	item1 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Comment: "Comment"},
		{Key: "name", Value: "containername", Comment: "Comment 2"},
		{Key: "image", Value: "image1"},
	}, Comment: "Comment"}
	item2 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "image", Value: "image2"},
		{Key: "name", Value: "containername"},
		{Key: "imagePullPolicy", Value: "IfNotPresent"},
	}}
	assert.True(sequenceItemMatch("containers", item1, item2))

	// mapslices with different value for identifying key
	item1 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "name", Value: "containername"},
		{Key: "image", Value: "image1"},
	}, Comment: "Comment"}
	item2 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "name", Value: "containername2"},
		{Key: "image", Value: "image1"},
	}}
	assert.False(sequenceItemMatch("containers", item1, item2))

	// Container.envFrom
	item1 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "configMapRef", Value: yaml.MapSlice{
			{Key: "name", Value: "n"},
			{Key: "optional", Value: true},
		}},
	}}
	item2 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "configMapRef", Value: yaml.MapSlice{
			{Key: "name", Value: "n"},
			{Key: "optional", Value: false},
		}},
		{Key: "prefix", Value: "pre"},
	}}
	assert.True(sequenceItemMatch("envFrom", item1, item2))
	item1 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "configMapRef", Value: yaml.MapSlice{{Key: "name", Value: "n"}}},
	}}
	item2 = yaml.SequenceItem{Value: yaml.MapSlice{
		{Key: "secretRef", Value: yaml.MapSlice{{Key: "name", Value: "n"}}},
	}}
	assert.False(sequenceItemMatch("envFrom", item1, item2))

}

func TestDeepEqual(t *testing.T) {
	assert := assert.New(t)

	var v1 interface{}
	var v2 interface{}

	v1 = yaml.MapSlice{{Key: "k", Value: "v", Comment: "comment"}}
	v2 = yaml.MapSlice{{Key: "k", Value: "v"}}
	assert.True(deepEqual(v1, v2))

	v1 = []yaml.SequenceItem{{Value: "v", Comment: "comment"}}
	v2 = []yaml.SequenceItem{{Value: "v"}}
	assert.True(deepEqual(v1, v2))

	v1 = yaml.MapSlice{{Key: "matchExpressions", Value: []yaml.SequenceItem{
		{Value: yaml.MapSlice{
			{Key: "key", Value: "labelkey"},
			{Key: "operator", Value: "In"},
			{Key: "values", Value: []yaml.SequenceItem{
				{Value: "value1", Comment: "Comment 1"},
				{Value: "value2"},
			}},
		}},
	}}}
	v2 = yaml.MapSlice{{Key: "matchExpressions", Value: []yaml.SequenceItem{
		{Value: yaml.MapSlice{
			{Key: "key", Value: "labelkey"},
			{Key: "operator", Value: "In"},
			{Key: "values", Value: []yaml.SequenceItem{
				{Value: "value1", Comment: "Comment 1"},
				{Value: "value2"},
			}},
		}},
	}}}
	assert.True(deepEqual(v1, v2))

	v1 = yaml.MapSlice{{Key: "matchExpressions", Value: []yaml.SequenceItem{
		{Value: yaml.MapSlice{
			{Key: "key", Value: "labelkey"},
			{Key: "operator", Value: "In"},
			{Key: "values", Value: []yaml.SequenceItem{
				{Value: "value1", Comment: "Comment 1"},
				{Value: "value2"},
			}},
		}},
	}}}
	v2 = yaml.MapSlice{{Key: "matchExpressions", Value: []yaml.SequenceItem{
		{Value: yaml.MapSlice{
			{Key: "key", Value: "labelkey"},
			{Key: "operator", Value: "In"},
			{Key: "values", Value: []yaml.SequenceItem{
				{Value: "value1"},
				{Value: "newvalue"},
			}},
		}},
	}}}
	assert.False(deepEqual(v1, v2))
}

package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
)

func TestMerge(t *testing.T) {
	assert := assert.New(t)

	// empty yaml
	_, err := Merge(nil, nil)
	assert.NotNil(err)
	_, err = Merge([]byte{}, []byte{})
	assert.NotNil(err)
	_, err = Merge([]byte(""), []byte(""))
	assert.NotNil(err)

	// invalid yaml
	_, err = Merge([]byte("a: b: c"), nil)
	assert.NotNil(err)
	_, err = Merge([]byte("a:"), nil)
	assert.NotNil(err)
	_, err = Merge([]byte("a: b"), []byte("a: b: c"))
	assert.NotNil(err)

	// non-map root node
	_, err = Merge([]byte("- a"), nil)
	assert.NotNil(err)
	_, err = Merge([]byte("a: b"), []byte("- a"))
	assert.NotNil(err)

	// valid yaml
	merged, err := Merge([]byte("a: b"), []byte("a: b"))
	assert.NoError(err)
	assert.Equal("a: b\n", string(merged))
}

func TestMergeMaps(t *testing.T) {
	cases := []struct {
		testName string
		orig     *yaml.Node
		fixed    *yaml.Node
		merged   *yaml.Node
	}{
		{
			"Update",
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
		{
			"Add",
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
		{
			"Remove",
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
		{
			"Preserve order",
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
		{
			"Sequence of strings",
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
		{
			"Map of maps",
			&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "a"},
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "b", LineComment: "EOL Comment 1"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "c", HeadComment: "Comment 1", LineComment: "EOL Comment 2"},
				}},
			}},
			&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "a"},
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "b"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "d"},
				}},
			}},
			&yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "a"},
				{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "b", LineComment: "EOL Comment 1"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "d", HeadComment: "Comment 1", LineComment: "EOL Comment 2"},
				}},
			}},
		},
		{
			"Complex mapslice (nested sequence and mapslices)",
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

	for _, test := range cases {
		t.Run(test.testName, func(t *testing.T) {
			merged := mergeMaps(test.orig, test.fixed)
			assert.True(t, deepEqual(test.merged, merged))
			assert.Equal(t, test.merged, merged)
		})
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
	_, index = findValInMap("k2", m)
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
	cases := []struct {
		testName    string
		item1       *yaml.Node
		item2       *yaml.Node
		sequenceKey string
		expected    bool
	}{
		{
			testName:    "Same strings",
			item1:       &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment"},
			item2:       &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
			sequenceKey: "",
			expected:    true,
		},
		{
			testName:    "Different strings",
			item1:       &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v", LineComment: "Comment"},
			item2:       &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2", LineComment: "Comment"},
			sequenceKey: "",
			expected:    false,
		},
		{
			testName: "String and mapslice",
			item1: &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, LineComment: "Comment"},
			}, LineComment: "Comment"},
			item2:       &yaml.Node{Kind: yaml.ScalarNode, Tag: strTag, Value: "v2"},
			sequenceKey: "",
			expected:    false,
		},
	}

	for _, test := range cases {
		t.Run(test.testName, func(t *testing.T) {
			assert.Equal(t, test.expected, sequenceItemMatch(test.sequenceKey, test.item1, test.item2))
		})
	}

	cases2 := []struct {
		testName    string
		sequenceKey string
		mapKey      string
	}{
		{"EndpointSubset.addresses : EndpointAddress.hostname", "addresses", "hostname"},
		{"EndpointSubset.addresses : EndpointAddress.ip", "addresses", "ip"},
		{"EndpointSubset.notReadyAddresses : EndpointAddress.hostname", "notReadyAddresses", "hostname"},
		{"EndpointSubset.notReadyAddresses : EndpointAddress.ip", "notReadyAddresses", "ip"},
		{"NetworkPolicySpec.ingress : NetworkPolicyIngressRule.ports", "ingress", "ports"},
		{"NetworkPolicySpec.ingress : NetworkPolicyIngressRule.from", "ingress", "from"},
		{"ConfigMapProjection.items : KeyToPath.key", "items", "key"},
		{"DownwardAPIVolumeSource.items : DownwardAPIVolumeFile.path", "items", "path"},
		{"NodeSelector.nodeSelectorTerms : NodeSelectorTerm.matchExpressions", "nodeSelectorTerms", "matchExpressions"},
		{"NodeSelector.nodeSelectorTerms : NodeSelectorTerm.matchFields", "nodeSelectorTerms", "matchFields"},
		{"ObjectMeta.ownerReferences : OwnerReference.uid", "ownerReferences", "uid"},
		{"ObjectMeta.ownerReferences : OwnerReference.name", "ownerReferences", "name"},
		{"NodeAffinity.preferredDuringSchedulingIgnoredDuringExecution : PreferredSchedulingTerm.preference", "preferredDuringSchedulingIgnoredDuringExecution", "preference"},
		{"PodAffinity.preferredDuringSchedulingIgnoredDuringExecution : WeightedPodAffinityTerm.podAffinityTerm", "preferredDuringSchedulingIgnoredDuringExecution", "podAffinityTerm"},
		{"ClusterRole.rules : PolicyRule.resources", "rules", "resources"},
		{"IngressSpec.rules : IngressRule.host", "rules", "host"},
		{"IngressSpec.rules : IngressRule.host", "rules", "host"},
		{"IngressSpec.tls : IngressTLS.secretName", "tls", "secretName"},
		{"IngressSpec.tls : IngressTLS.hosts", "tls", "hosts"},
	}

	for _, test := range cases2 {
		t.Run(test.testName, func(t *testing.T) {
			item1 := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: test.mapKey},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "value"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "randkey"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "randval"},
			}}
			item2 := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: test.mapKey},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "value"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "randkey"},
				{Kind: yaml.ScalarNode, Tag: strTag, Value: "otherval"},
			}}
			assert.True(t, sequenceItemMatch(test.sequenceKey, item1, item2))
			item2.Content[1].Value = "othervalue"
			assert.False(t, sequenceItemMatch(test.sequenceKey, item1, item2))
		})
	}

	// nested maps
	cases3 := []struct {
		testName           string
		sequenceKey        string
		intermediateMapKey string
		mapKey             string
	}{
		{"ProjectedVolumeSource.sources : VolumeProjection.configMap.name", "sources", "configMap", "name"},
		{"ProjectedVolumeSource.sources : VolumeProjection.downwardAPI.items", "sources", "downwardAPI", "items"},
		{"ProjectedVolumeSource.sources : VolumeProjection.secret.name", "sources", "secret", "name"},
		{"ProjectedVolumeSource.sources : VolumeProjection.serviceAccountToken.name", "sources", "serviceAccountToken", "path"},
		{"Container.envFrom : EnvFromSource.configMapRef.name", "envFrom", "configMapRef", "name"},
		{"Container.envFrom : EnvFromSource.secretRef.name", "envFrom", "secretRef", "name"},
	}

	for _, test := range cases3 {
		t.Run(test.testName, func(t *testing.T) {
			item1 := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: test.intermediateMapKey},
				{Kind: yaml.MappingNode, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: test.mapKey},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "randkey"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "randval"},
				}},
			}}
			item2 := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: strTag, Value: test.intermediateMapKey},
				{Kind: yaml.MappingNode, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: strTag, Value: test.mapKey},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "value"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "randkey"},
					{Kind: yaml.ScalarNode, Tag: strTag, Value: "otherval"},
				}},
			}}
			assert.True(t, sequenceItemMatch(test.sequenceKey, item1, item2))
			item2.Content[1].Content[1].Value = "othervalue"
			assert.False(t, sequenceItemMatch(test.sequenceKey, item1, item2))
			item2.Content[0].Value = "bla"
			assert.False(t, sequenceItemMatch(test.sequenceKey, item1, item2))
		})
	}
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
	v2 = &yaml.Node{Kind: yaml.MappingNode, Tag: mapTag, Content: []*yaml.Node{
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

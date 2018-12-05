package cmd

import "sort"

// CapSet represents a set of capabilities.
type CapSet map[CapabilityV1]bool

// NewCapSetFromArray converts an array of capabilities into a CapSet.
func NewCapSetFromArray(array []CapabilityV1) (set CapSet) {
	set = make(CapSet)
	for _, cap := range array {
		set[cap] = true
	}
	return
}

func mergeCapSets(sets ...CapSet) (merged CapSet) {
	merged = make(CapSet)
	for _, set := range sets {
		for k, v := range set {
			merged[k] = v
		}
	}
	return
}

func sortCapSet(capSet CapSet) (sorted []CapabilityV1) {
	keys := []string{}
	for key := range capSet {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)

	for _, key := range keys {
		sorted = append(sorted, CapabilityV1(key))
	}
	return
}

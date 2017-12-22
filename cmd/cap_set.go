package cmd

import "sort"

type CapSet map[Capability]bool

func NewCapSetFromArray(array []Capability) (set CapSet) {
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

func sortCapSet(capSet CapSet) (sorted []Capability) {
	keys := []string{}
	for key := range capSet {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)

	for _, key := range keys {
		sorted = append(sorted, Capability(key))
	}
	return
}

package cmd

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

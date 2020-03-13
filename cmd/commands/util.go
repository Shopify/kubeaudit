package commands

import "strings"

// Implements pflag.Value
// Custom type so slice values can be represented as a space-separated list of strings in the flag value
//
// Example:
//
//     var vals []string
//     defaultVals := []string{"default1", "default2"}
//
//     Flags().VarP(newStringSliceValue(strings.Join(defaultVals, " "), &vals), "flagname", "f", "Flag description")
//
type stringSliceValue []string

func newStringSliceValue(val string, p *[]string) *stringSliceValue {
	*p = strings.Split(val, " ")
	return (*stringSliceValue)(p)
}

func (s *stringSliceValue) Set(val string) error {
	*s = stringSliceValue(strings.Split(val, " "))
	return nil
}
func (s *stringSliceValue) Type() string {
	return "stringSlice"
}

func (s *stringSliceValue) String() string { return strings.Join(*s, " ") }

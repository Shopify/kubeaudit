package cmd

import "testing"

func TestNamespaceMissingDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "namespace_missing_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyNetworkPolicy})
}

func TestNamespaceHasDefaulDenyNetPol(t *testing.T) {
	expectedMessage := "Namespace has a default deny NetworkPolicy"
	results := runAuditTest(t, "namespace_has_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyNetworkPolicy})
	for _, result := range results {
		if result.Occurrences[0].message != expectedMessage {
			t.Logf("Expected: %s\nGot: %s", expectedMessage, result.Occurrences[0].message)
			t.Fail()
		}
	}
}

func TestNamespaceHasDefaulDenyAndAllowAllNetPol(t *testing.T) {
	runAuditTest(t, "namespace_has_default_deny_and_allow_all_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyNetworkPolicy, ErrorMissingDefaultDenyNetworkPolicy})
}

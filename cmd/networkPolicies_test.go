package cmd

import "testing"

func TestNamespaceMissingDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "namespace_missing_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyNetworkPolicy})
}

func TestNamespaceHasDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "namespace_has_default_deny_netpol.yml", auditNetworkPolicies, []int{InfoDefaultDenyNetworkPolicyExists})
}

func TestNamespaceHasDefaulDenyAndAllowAllNetPol(t *testing.T) {
	runAuditTest(t, "namespace_has_default_deny_and_allow_all_netpol.yml", auditNetworkPolicies, []int{InfoDefaultDenyNetworkPolicyExists, WarningAllowAllIngressNetworkPolicyExists})
}

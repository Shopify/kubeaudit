package cmd

import "testing"

func TestNamespaceMissingDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "namespace_missing_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressNetworkPolicy, ErrorMissingDefaultDenyEgressNetworkPolicy})
}

func TestNamespaceMissingDefaultDenyEgressNetPol(t *testing.T) {
	runAuditTest(t, "namespace_missing_default_deny_egress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyEgressNetworkPolicy})
}

func TestNamespaceMissingDefaultDenyIngressNetPol(t *testing.T) {
	t.Skip()
	runAuditTest(t, "namespace_missing_default_deny_ingress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressNetworkPolicy})
}

func TestNamespaceHasDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "namespace_has_default_deny_netpol.yml", auditNetworkPolicies, []int{InfoDefaultDenyNetworkPolicyExists})
}

func TestNamespaceHasDefaulDenyAndAllowAllNetPol(t *testing.T) {
	runAuditTest(t, "namespace_has_default_deny_and_allow_all_netpol.yml", auditNetworkPolicies, []int{InfoDefaultDenyNetworkPolicyExists, WarningAllowAllIngressNetworkPolicyExists, WarningAllowAllEgressNetworkPolicyExists})
}

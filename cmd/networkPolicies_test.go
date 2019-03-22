package cmd

import "testing"

func TestNamespaceMissingDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "namespace_missing_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressAndEgressNetworkPolicy})
}

func TestAllowAuditNamespaceMissingDefaulDenyNetPolFromConfig(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_audit_from_config.yml"
	runAuditTest(t, "namespace_missing_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressAndEgressNetworkPolicy})
}

func TestNamespaceMissingDefaultDenyEgressNetPol(t *testing.T) {
	runAuditTest(t, "namespace_missing_default_deny_egress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyEgressNetworkPolicy})
}

func TestNamespaceMissingDefaultDenyIngressNetPol(t *testing.T) {
	runAuditTest(t, "namespace_missing_default_deny_ingress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressNetworkPolicy})
}

func TestNamespaceHasDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "namespace_has_default_deny_netpol.yml", auditNetworkPolicies, []int{InfoDefaultDenyNetworkPolicyExists})
}

func TestNamespaceHasDefaulDenyAndAllowAllNetPol(t *testing.T) {
	runAuditTest(t, "namespace_has_default_deny_and_allow_all_netpol.yml", auditNetworkPolicies, []int{InfoDefaultDenyNetworkPolicyExists, WarningAllowAllIngressNetworkPolicyExists, WarningAllowAllEgressNetworkPolicyExists})
}

func TestAllowedNamespaceMissingDefaulDenyNetPol(t *testing.T) {
	runAuditTest(t, "allowed_namespace_missing_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressAndEgressNetworkPolicyAllowed})
}

func TestAllowedNamespaceMissingDefaultDenyEgressNetPol(t *testing.T) {
	runAuditTest(t, "allowed_namespace_missing_default_deny_egress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyEgressNetworkPolicyAllowed})
}

func TestAllowedNamespaceMissingDefaultDenyIngressNetPol(t *testing.T) {
	runAuditTest(t, "allowed_namespace_missing_default_deny_ingress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressNetworkPolicyAllowed})
}

func TestAllowedNamespaceMissingDefaulDenyNetPolFromConfig(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_namespace_missing_default_deny_net_pol_from_config.yml"
	runAuditTest(t, "namespace_missing_default_deny_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressAndEgressNetworkPolicyAllowed})
	rootConfig.auditConfig = ""
}

func TestAllowedNamespaceMissingDefaultDenyEgressNetPolFromConfig(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_namespace_missing_default_deny_egress_net_pol_from_config.yml"
	runAuditTest(t, "namespace_missing_default_deny_egress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyEgressNetworkPolicyAllowed})
	rootConfig.auditConfig = ""
}

func TestAllowedNamespaceMissingDefaultDenyIngressNetPolFromConfig(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_namespace_missing_default_deny_ingress_net_pol_from_config.yml"
	runAuditTest(t, "namespace_missing_default_deny_ingress_netpol.yml", auditNetworkPolicies, []int{ErrorMissingDefaultDenyIngressNetworkPolicyAllowed})
	rootConfig.auditConfig = ""
}

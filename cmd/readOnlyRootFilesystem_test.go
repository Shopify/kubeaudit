package cmd

import "testing"

func TestSecurityContextNilRORFV1(t *testing.T) {
	runAuditTest(t, "security_context_nil_v1.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemNil})
}

func TestReadOnlyRootFilesystemNilV1(t *testing.T) {
	runAuditTest(t, "read_only_root_filesystem_nil_v1.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemNil})
}

func TestReadOnlyRootFilesystemFalseV1(t *testing.T) {
	runAuditTest(t, "read_only_root_filesystem_false_v1.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemFalse})
}

func TestReadOnlyRootFilesystemFalseAllowedV1(t *testing.T) {
	runAuditTest(t, "read_only_root_filesystem_false_allowed_v1.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemFalseAllowed})
}

func TestReadOnlyRootFilesystemMisconfiguredAllowV1(t *testing.T) {
	runAuditTest(t, "read_only_root_filesystem_misconfigured_allow_v1.yml", auditReadOnlyRootFS, []int{ErrorMisconfiguredKubeauditAllow})
}

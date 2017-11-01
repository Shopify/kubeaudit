package cmd

import "testing"

func TestSecurityContextNIL_RORF(t *testing.T) {
	runTest(t, "security_context_nil.yml", auditReadOnlyRootFS, ErrorSecurityContextNIL)
}

func TestReadOnlyRootFilesystemNIL(t *testing.T) {
	runTest(t, "read_only_root_filesystem_nil.yml", auditReadOnlyRootFS, ErrorReadOnlyRootFilesystemNIL)
}

func TestReadOnlyRootFilesystemFalse(t *testing.T) {
	runTest(t, "read_only_root_filesystem_false.yml", auditReadOnlyRootFS, ErrorReadOnlyRootFilesystemFalse)
}

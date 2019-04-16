package cmd

// KubeauditConfig sets up config for kubeaudit from flag `config`
type KubeauditConfig struct {
	APIVersion string               `yaml:"apiVersion"`
	Kind       string               `yaml:"kind"`
	Spec       *KubeauditConfigSpec `yaml:"spec"`
	Audit      bool                 `yaml:"audit"`
}

// KubeauditConfigSpec contains Config Spec
type KubeauditConfigSpec struct {
	Manifest     []*KubeauditConfigManifest   `yaml:"manifest"`
	Capabilities *KubeauditConfigCapabilities `yaml:"capabilities"`
	Overrides    *KubeauditConfigOverrides    `yaml:"overrides"`
}

// KubeauditConfigManifest contains path to the manifests to audit
type KubeauditConfigManifest struct {
	Path string `yaml:"path"`
}

// KubeauditConfigCapabilities contains list of capabilities supported
type KubeauditConfigCapabilities struct {
	NetAdmin       string `yaml:"NET_ADMIN"`
	SetPCAP        string `yaml:"SETPCAP"`
	MKNOD          string `yaml:"MKNOD"`
	AuditWrite     string `yaml:"AUDIT_WRITE"`
	Chown          string `yaml:"CHOWN"`
	NetRaw         string `yaml:"NET_RAW"`
	DacOverride    string `yaml:"DAC_OVERRIDE"`
	FOWNER         string `yaml:"FOWNER"`
	FSetID         string `yaml:"FSETID"`
	Kill           string `yaml:"KILL"`
	SetGID         string `yaml:"SETGID"`
	SetUID         string `yaml:"SETUID"`
	NetBindService string `yaml:"NET_BIND_SERVICE"`
	SYSChroot      string `yaml:"SYS_CHROOT"`
	SetFCAP        string `yaml:"SETFCAP"`
}

// KubeauditConfigOverrides contains list of available overrides
type KubeauditConfigOverrides struct {
	PrivilegeEscalation                string `yaml:"privilege-escalation"`
	Privileged                         string `yaml:"privileged"`
	RunAsRoot                          string `yaml:"run-as-root"`
	AutomountServiceAccountToken       string `yaml:"automount-service-account-token"`
	ReadOnlyRootFilesystemFalse        string `yaml:"read-only-root-filesystem-false"`
	NonDefaultDenyIngressNetworkPolicy string `yaml:"non-default-deny-ingress-network-policy"`
	NonDefaultDenyEgressNetworkPolicy  string `yaml:"non-default-deny-egress-network-policy"`
	HostNetwork                        string `yaml:"namespace-host-network"`
	HostPID                            string `yaml:"namespace-host-PID"`
	HostIPC                            string `yaml:"namespace-host-IPC"`
}

func mapOverridesToStructFields(label string) string {
	switch label {
	case "allow-privilege-escalation":
		return "PrivilegeEscalation"
	case "allow-privileged":
		return "Privileged"
	case "allow-run-as-root":
		return "RunAsRoot"
	case "allow-automount-service-account-token":
		return "AutomountServiceAccountToken"
	case "allow-read-only-root-filesystem-false":
		return "ReadOnlyRootFilesystemFalse"
	case "allow-non-default-deny-egress-network-policy":
		return "NonDefaultDenyEgressNetworkPolicy"
	case "allow-non-default-deny-ingress-network-policy":
		return "NonDefaultDenyIngressNetworkPolicy"
	case "allow-namespace-host-network":
		return "HostNetwork"
	case "allow-namespace-host-IPC":
		return "HostIPC"
	case "allow-namespace-host-PID":
		return "HostPID"
	}
	return ""
}

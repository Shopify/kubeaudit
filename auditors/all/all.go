package all

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/asat"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/hostns"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"github.com/Shopify/kubeaudit/auditors/mountds"
	"github.com/Shopify/kubeaudit/auditors/netpols"
	"github.com/Shopify/kubeaudit/auditors/nonroot"
	"github.com/Shopify/kubeaudit/auditors/privesc"
	"github.com/Shopify/kubeaudit/auditors/privileged"
	"github.com/Shopify/kubeaudit/auditors/rootfs"
	"github.com/Shopify/kubeaudit/auditors/seccomp"
)

func Auditors() []kubeaudit.Auditable {
	// An error occurs when the passed in cpu and memory limits can't be parsed, but we are just using defaults so it
	// is safe to ignore the error
	limitsAuditor, _ := limits.New("", "")

	return []kubeaudit.Auditable{
		apparmor.New(),
		asat.New(),
		capabilities.New(nil),
		hostns.New(),
		image.New(""),
		limitsAuditor,
		mountds.New(),
		netpols.New(),
		nonroot.New(),
		privesc.New(),
		privileged.New(),
		rootfs.New(),
		seccomp.New(),
	}
}

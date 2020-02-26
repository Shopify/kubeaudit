package all

import (
	"errors"
	"fmt"

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
	"github.com/Shopify/kubeaudit/config"
)

var ErrUnknownAuditor = errors.New("Unknown auditor")

var AuditorNames = []string{
	apparmor.Name,
	asat.Name,
	capabilities.Name,
	hostns.Name,
	image.Name,
	limits.Name,
	mountds.Name,
	netpols.Name,
	nonroot.Name,
	privesc.Name,
	privileged.Name,
	rootfs.Name,
	seccomp.Name,
}

func Auditors(conf config.KubeauditConfig) ([]kubeaudit.Auditable, error) {
	enabledAuditors := conf.GetEnabledAuditors()
	if len(enabledAuditors) == 0 {
		enabledAuditors = AuditorNames
	}

	auditors := make([]kubeaudit.Auditable, 0, len(enabledAuditors))
	for _, auditorName := range enabledAuditors {
		auditor, err := initAuditor(auditorName, conf)
		if err != nil {
			return nil, err
		}
		auditors = append(auditors, auditor)
	}

	return auditors, nil
}

func initAuditor(name string, conf config.KubeauditConfig) (kubeaudit.Auditable, error) {
	switch name {
	case apparmor.Name:
		return apparmor.New(), nil
	case asat.Name:
		return asat.New(), nil
	case capabilities.Name:
		return capabilities.New(conf.GetAuditorConfigs().Capabilities), nil
	case hostns.Name:
		return hostns.New(), nil
	case image.Name:
		return image.New(conf.GetAuditorConfigs().Image), nil
	case limits.Name:
		return limits.New(conf.GetAuditorConfigs().Limits)
	case mountds.Name:
		return mountds.New(), nil
	case netpols.Name:
		return netpols.New(), nil
	case nonroot.Name:
		return nonroot.New(), nil
	case privesc.Name:
		return privesc.New(), nil
	case privileged.Name:
		return privileged.New(), nil
	case rootfs.Name:
		return rootfs.New(), nil
	case seccomp.Name:
		return seccomp.New(), nil
	}

	return nil, fmt.Errorf("unknown auditor %s: %w", name, ErrUnknownAuditor)
}

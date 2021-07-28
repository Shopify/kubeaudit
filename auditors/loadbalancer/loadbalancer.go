package loadbalancer

import (
	"fmt"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
)

const LoadbalancerType = "LoadBalancer"

const (
	ExposedService = "ExposedService"
)

type Loadbalancer struct{}

func New() *Loadbalancer {
	return &Loadbalancer{}
}

func (lb *Loadbalancer) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	// Get/parse service resource
	serviceMeta := k8s.GetObjectMeta(resource)
	if serviceMeta == nil {
		return nil, fmt.Errorf("Cannot parse resource")
	}

	// Check for type
	serviceSpec := k8s.GetServiceSpec(resource)
	if serviceSpec == nil {
		return nil, fmt.Errorf("Cannot parse service spec")
	}

	// Return result if it's a type load balancer
	if serviceSpec.Type == LoadbalancerType {
		auditResults = append(auditResults, &kubeaudit.AuditResult{
			Name:     ExposedService,
			Severity: kubeaudit.Warn,
			Message:  fmt.Sprintf("Service is exposed to the internet. Service name: %s, namespace: %s\n", serviceMeta.Name, serviceMeta.Namespace),
			Metadata: kubeaudit.Metadata{
				"Type": LoadbalancerType,
			},
		})
	}
	return auditResults, nil
}

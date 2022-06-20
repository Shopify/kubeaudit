package kubeaudit

import (
	"bytes"
	"fmt"

	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"gopkg.in/yaml.v3"
)

func getResourcesFromClient(client k8sinternal.KubeClient, options k8sinternal.ClientOptions) ([]KubeResource, error) {
	var resources []KubeResource

	k8sresources, err := client.GetAllResources(options)
	if err != nil {
		return nil, err
	}
	for _, resource := range k8sresources {
		resources = append(resources, &kubeResource{object: resource})
	}

	return resources, nil
}

func getResourcesFromManifest(data []byte) ([]KubeResource, error) {
	var resources []KubeResource
	bufSlice := bytes.Split(data, []byte("---"))

	for _, b := range bufSlice {
		obj, err := k8sinternal.DecodeResource(b)
		if err == nil && obj != nil {
			source := &kubeResource{
				object: obj,
				bytes:  b,
			}
			resources = append(resources, source)
		} else if err := yaml.Unmarshal(data, &yaml.Node{}); err != nil {
			return nil, fmt.Errorf("Invalid yaml: %w", err)
		} else {
			resources = append(resources, &kubeResource{bytes: b})
		}
	}

	return resources, nil
}

func auditResources(resources []KubeResource, auditable []Auditable) ([]Result, error) {
	var results []Result

	for _, resource := range resources {
		result, err := auditResource(resource, resources, auditable)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func auditResource(resource KubeResource, resources []KubeResource, auditables []Auditable) (Result, error) {
	result := &workloadResult{
		Resource:     resource,
		AuditResults: []*AuditResult{},
	}

	if resource.Object() == nil {
		return result, nil
	}

	for _, auditable := range auditables {
		auditResults, err := auditable.Audit(resource.Object(), unwrapResources(resources))
		if err != nil {
			return nil, err
		}
		result.AuditResults = append(result.AuditResults, auditResults...)
	}

	return result, nil
}

func unwrapResources(resources []KubeResource) []k8s.Resource {
	unwrappedResources := make([]k8s.Resource, 0, len(resources))
	for _, resource := range resources {
		unwrappedResources = append(unwrappedResources, resource.Object())
	}
	return unwrappedResources
}

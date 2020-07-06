package kubeaudit

import (
	"bytes"
	"fmt"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes"
)

func getResourcesFromClientset(clientset kubernetes.Interface, options k8s.ClientOptions) []KubeResource {
	var resources []KubeResource

	for _, resource := range k8s.GetAllResources(clientset, options) {
		resources = append(resources, &kubeResource{object: resource})
	}

	return resources
}

func getResourcesFromManifest(data []byte) ([]KubeResource, error) {
	var resources []KubeResource
	bufSlice := bytes.Split(data, []byte("---"))

	for _, b := range bufSlice {
		obj, err := k8s.DecodeResource(b)
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

	if !k8stypes.IsSupportedResourceType(resource.Object()) {
		resourceInfo := resource.Object().GetObjectKind().GroupVersionKind()
		auditResult := &AuditResult{
			Name:     ErrorUnsupportedResource,
			Severity: Warn,
			Message:  "Resource is not supported.",
			Metadata: Metadata{
				"Kind":    resourceInfo.Kind,
				"Group":   resourceInfo.Group,
				"Version": resourceInfo.Version,
			},
		}
		result.AuditResults = append(result.AuditResults, auditResult)
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

func unwrapResources(resources []KubeResource) []k8stypes.Resource {
	unwrappedResources := make([]k8stypes.Resource, 0, len(resources))
	for _, resource := range resources {
		unwrappedResources = append(unwrappedResources, resource.Object())
	}
	return unwrappedResources
}

func logAuditResult(result *AuditResult, baseLogger *log.Logger) {
	logger := baseLogger.WithFields(getLogFieldsForResult(result))
	switch result.Severity {
	case Info:
		logger.Info(result.Message)
	case Warn:
		logger.Warn(result.Message)
	case Error:
		logger.Error(result.Message)
	}
}

func getLogFieldsForResult(result *AuditResult) log.Fields {
	fields := log.Fields{
		"AuditResultName": result.Name,
	}

	for k, v := range result.Metadata {
		fields[k] = v
	}

	return fields
}

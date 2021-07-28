package kubeaudit

import (
	"bytes"

	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/internal/yaml"
	"github.com/Shopify/kubeaudit/pkg/k8s"
)

func fix(results []Result) ([]byte, error) {
	var outputBytes [][]byte
	var newResources []k8s.Resource

	if len(results) == 0 {
		return []byte{}, nil
	}

	// Fix all the resources
	for _, result := range results {
		for _, auditResult := range result.GetAuditResults() {
			newResources = append(newResources, auditResult.Fix(result.GetResource().Object())...)
		}
	}

	// Convert all the resources to bytes
	for _, result := range results {
		if result.GetResource().Object() == nil {
			outputBytes = append(outputBytes, result.GetResource().Bytes())
			continue
		}

		fixedresourceBytes, err := resourceToBytes(result.GetResource().Object(), result.GetResource().Bytes())
		if err != nil {
			return nil, err
		}

		outputBytes = append(outputBytes, fixedresourceBytes)
	}

	// Convert all the new resources to bytes
	for _, newResource := range newResources {
		fixedresourceBytes, err := resourceToBytes(newResource, nil)
		if err != nil {
			return nil, err
		}
		outputBytes = append(outputBytes, fixedresourceBytes)
	}

	fixedManifest := bytes.Join(outputBytes, []byte("---"))

	return fixedManifest, nil
}

func resourceToBytes(fixedResource k8s.Resource, origResourceBytes []byte) ([]byte, error) {
	fixedresourceBytes, err := k8sinternal.EncodeResource(fixedResource)
	if err != nil {
		return nil, err
	}

	if origResourceBytes == nil {
		// This is a new resource (not in the original manifest)
		// Add  a leading newline
		fixedresourceBytes = append([]byte{'\n'}, fixedresourceBytes...)
	} else {
		fixedresourceBytes, err = yaml.Merge(origResourceBytes, fixedresourceBytes)
		if err != nil {
			return nil, err
		}

		// Add any leading and trailing whitespace that was present in the original
		fixedresourceBytes = bytes.Replace(origResourceBytes, bytes.TrimSpace(origResourceBytes), fixedresourceBytes, 1)

		// Remove the redundant trailing newline
		if fixedresourceBytes[len(fixedresourceBytes)-1] == '\n' {
			fixedresourceBytes = fixedresourceBytes[:len(fixedresourceBytes)-1]
		}
	}

	fixedresourceBytes, err = cleanupManifest(origResourceBytes, fixedresourceBytes)
	if err != nil {
		return nil, err
	}

	return fixedresourceBytes, nil
}

// TODO do this better??
func cleanupManifest(origData, finalData []byte) ([]byte, error) {
	objectMetacreationTs := []byte("\n  creationTimestamp: null\n")
	specTemplatecreationTs := []byte("\n      creationTimestamp: null\n")
	jobSpecTemplatecreationTs := []byte("\n          creationTimestamp: null\n")
	nullStatus := []byte("\nstatus: {}\n")
	nullReplicaStatus := []byte("status:\n  replicas: 0\n")
	nullLBStatus := []byte("status:\n  loadBalancer: {}\n")
	nullMetaStatus := []byte("\n    status: {}\n")

	var hasObjectMetacreationTs, hasSpecTemplatecreationTs, hasJobSpecTemplatecreationTs, hasNullStatus,
		hasNullReplicaStatus, hasNullLBStatus, hasNullMetaStatus bool

	if origData != nil {
		hasObjectMetacreationTs = bytes.Contains(origData, objectMetacreationTs)
		hasSpecTemplatecreationTs = bytes.Contains(origData, specTemplatecreationTs)
		hasJobSpecTemplatecreationTs = bytes.Contains(origData, jobSpecTemplatecreationTs)

		hasNullStatus = bytes.Contains(origData, nullStatus)
		hasNullReplicaStatus = bytes.Contains(origData, nullReplicaStatus)
		hasNullLBStatus = bytes.Contains(origData, nullLBStatus)
		hasNullMetaStatus = bytes.Contains(origData, nullMetaStatus)
	} // null value is false in case of origFile

	if !hasObjectMetacreationTs {
		finalData = bytes.Replace(finalData, objectMetacreationTs, []byte("\n"), -1)
	}
	if !hasSpecTemplatecreationTs {
		finalData = bytes.Replace(finalData, specTemplatecreationTs, []byte("\n"), -1)
	}
	if !hasJobSpecTemplatecreationTs {
		finalData = bytes.Replace(finalData, jobSpecTemplatecreationTs, []byte("\n"), -1)
	}
	if !hasNullStatus {
		finalData = bytes.Replace(finalData, nullStatus, []byte("\n"), -1)
	}
	if !hasNullReplicaStatus {
		finalData = bytes.Replace(finalData, nullReplicaStatus, []byte("\n"), -1)
	}
	if !hasNullLBStatus {
		finalData = bytes.Replace(finalData, nullLBStatus, []byte("\n"), -1)
	}
	if !hasNullMetaStatus {
		finalData = bytes.Replace(finalData, nullMetaStatus, []byte("\n"), -1)
	}

	return finalData, nil
}

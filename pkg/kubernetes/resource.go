package kubernetes

import (
	"bytes"
	"fmt"

	"github.com/fllaca/okay/pkg/utils"
	"sigs.k8s.io/yaml"
)

type Resource map[string]interface{}

func GetApiVersionKind(resource map[string]interface{}) string {
	return utils.JoinNotEmptyStrings("/", GetString(resource, "apiVersion"), GetString(resource, "kind"))
}

func GetMetadata(resource map[string]interface{}) map[string]interface{} {
	value, _ := resource["metadata"].(map[string]interface{})
	return value
}

func GetName(resource map[string]interface{}) string {
	metadata := GetMetadata(resource)
	if metadata == nil {
		return ""
	}
	return GetString(metadata, "name")
}

func GetNamespace(resource map[string]interface{}) string {
	metadata := GetMetadata(resource)
	if metadata == nil {
		return ""
	}
	return GetString(metadata, "namespace")
}

func GetString(input map[string]interface{}, key string) string {
	value, _ := input[key].(string)
	return value
}

func ParseResourcesFromYaml(fileBytes []byte) ([]Resource, error) {
	result := make([]Resource, 0)
	documentsBytes := bytes.Split(fileBytes, []byte("\n---\n"))

	for docIndex, documentBytes := range documentsBytes {
		resource := make(Resource)
		err := yaml.Unmarshal(documentBytes, &resource)
		if err != nil {
			return nil, fmt.Errorf("Error parsing resource from document %d: %s", docIndex, err)
		}
		if len(resource) > 0 {
			result = append(result, resource)
		}
	}
	return result, nil
}

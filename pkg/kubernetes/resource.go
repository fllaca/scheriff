package kubernetes

import (
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

func ParseResource(resourceBytes []byte) (map[string]interface{}, error) {
	var resource map[string]interface{}
	err := yaml.Unmarshal(resourceBytes, &resource)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return resource, nil
}

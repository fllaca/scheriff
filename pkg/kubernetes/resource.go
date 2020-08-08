package kubernetes

import (
	"github.com/fllaca/okay/pkg/utils"
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

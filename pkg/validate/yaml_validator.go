package validate

import (
	"bytes"
	"fmt"
	"sigs.k8s.io/yaml"
)

type YamlFileValidator struct {
	resourceValidator ResourceValidator
}

func NewYamlFileValidator(resourceValidator ResourceValidator) YamlFileValidator {
	return YamlFileValidator{
		resourceValidator: resourceValidator,
	}
}

func (yamlValidator YamlFileValidator) Validate(fileBytes []byte) []ValidationResult {
	result := make([]ValidationResult, 0)
	documentsBytes := bytes.Split(fileBytes, []byte("\n---\n"))
	for docIndex, documentBytes := range documentsBytes {
		k8sResource, err := parseResource(documentBytes)
		if err != nil {
			result = append(result, ValidationResult{
				Message:  fmt.Sprintf("Error parsing k8s resource from document %d: %s\n", docIndex, err),
				Severity: SeverityError,
			})
			continue
		}
		if len(k8sResource) == 0 {
			continue
		}
		result = append(result, yamlValidator.resourceValidator.Validate(k8sResource))
	}
	return result
}

func parseResource(resourceBytes []byte) (map[string]interface{}, error) {
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

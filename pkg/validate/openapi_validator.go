package validate

import (
	"encoding/json"
	"fmt"

	"github.com/fllaca/okay/pkg/kubernetes"
	"github.com/fllaca/okay/pkg/utils"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// extPropsGroupVersionKind holds the data inside the "x-kubernetes-group-version-kind" Extension Properties of Kubernetes schemas
type extPropsGroupVersionKind struct {
	Group   string `json:"group" yaml:"group"`
	Kind    string `json:"kind" yaml:"kind"`
	Version string `json:"version" yaml:"version"`
}

func (extPropsGroupVersionKind extPropsGroupVersionKind) String() string {
	return utils.JoinNotEmptyStrings("/", extPropsGroupVersionKind.Group, extPropsGroupVersionKind.Version, extPropsGroupVersionKind.Kind)
}

// OpenApiValidator validates Kubernetes manifests using OpenApi schemas
type OpenApiValidator struct {
	schemaCache map[string]*openapi3.Schema
}

func NewOpenApi2Validator(openApi2SpecsBytes []byte) (*OpenApiValidator, error) {
	// TODO: set to false when using verbose output
	openapi3.SchemaErrorDetailsDisabled = true
	swagger2 := &openapi2.Swagger{}

	err := json.Unmarshal(openApi2SpecsBytes, swagger2)
	if err != nil {
		return nil, err
	}
	// In kubernetes API specs this field is specifed as "type: string", although integers are also accepted
	/* Alternative: add definition before converting to openapi3
	swagger2.Definitions["io.k8s.apimachinery.pkg.util.intstr.IntOrString"] = &openapi3.SchemaRef{
		Value: openapi3.NewOneOfSchema(
			openapi3.NewStringSchema(),
			openapi3.NewInt32Schema()),
	}
	*/

	swagger3, err := openapi2conv.ToV3Swagger(swagger2)
	swagger3.Components.Schemas["io.k8s.apimachinery.pkg.util.intstr.IntOrString"] = &openapi3.SchemaRef{
		Value: openapi3.NewOneOfSchema(
			openapi3.NewStringSchema(),
			openapi3.NewInt32Schema()),
	}
	// resolve references to the new IntOrString schema
	{
		sl := openapi3.NewSwaggerLoader()
		if err := sl.ResolveRefsIn(swagger3, nil); err != nil {
			return nil, err
		}
	}

	// build schemaCache:
	schemaCache, err := buildSchemaCache(swagger3)

	if err != nil {
		return nil, err
	}

	return &OpenApiValidator{
		schemaCache: schemaCache,
	}, nil
}

func (oeValidator OpenApiValidator) Validate(input map[string]interface{}) ValidationResult {

	kind := kubernetes.GetApiVersionKind(input)
	name := kubernetes.GetName(input)
	namespace := kubernetes.GetNamespace(input)
	schema := oeValidator.schemaCache[kind]

	result := ValidationResult{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
	}

	if schema == nil {
		result.Message = fmt.Sprintf("Kind '%s' not found in schema", kind)
		result.Severity = SeverityWarning
		return result
	}

	err := schema.VisitJSON(input)

	if err != nil {
		result.Message = err.Error()
		result.Severity = SeverityError
		return result
	}

	result.Message = "valid"
	result.Severity = SeverityOK
	return result
}

func buildSchemaCache(swagger3 *openapi3.Swagger) (map[string]*openapi3.Schema, error) {
	schemaCache := make(map[string]*openapi3.Schema)
	for _, schema := range swagger3.Components.Schemas {
		kindDefs, err := getK8sGroupVersionKind(schema)
		if err != nil {
			return nil, fmt.Errorf("Cannot load GroupVersionKind from: %v, %s", schema, err)
		}
		for _, kindDef := range kindDefs {
			schemaCache[kindDef.String()] = schema.Value
		}
	}
	return schemaCache, nil
}

func getK8sGroupVersionKind(schema *openapi3.SchemaRef) ([]extPropsGroupVersionKind, error) {
	var kindDefs []extPropsGroupVersionKind = make([]extPropsGroupVersionKind, 0)
	data := schema.Value.ExtensionProps.Extensions["x-kubernetes-group-version-kind"]
	if data == nil {
		return kindDefs, nil
	}
	if k8sExtensionBytes, ok := data.(json.RawMessage); ok {
		err := json.Unmarshal(k8sExtensionBytes, &kindDefs)
		if err != nil {
			return nil, err
		}
		return kindDefs, nil
	}
	return kindDefs, nil
}

const (
	crdv1beta1ApiVersionKind = "apiextensions.k8s.io/v1beta1/CustomResourceDefinition"
	crdv1ApiVersionKind      = "apiextensions.k8s.io/v1/CustomResourceDefinition"
)

// AddCrdSchemas adds additional schemas from a CustomResourceDefinition that can be used to validate other resources
// TODO:  <09-08-20, @fllaca> // more generic way of handling CRDs? :thinking:
func (oeValidator OpenApiValidator) AddCrdSchemas(crdResource kubernetes.Resource) error {
	apiVersionKind := kubernetes.GetApiVersionKind(crdResource)
	switch apiVersionKind {
	case crdv1ApiVersionKind:
		crdv1 := apiextensionsv1.CustomResourceDefinition{}
		err := convertObject(crdResource, &crdv1)
		if err != nil {
			return err
		}
		for _, version := range crdv1.Spec.Versions {
			schema, err := getSchema(version.Schema.OpenAPIV3Schema)
			if err != nil {
				return err
			}
			kindDef := extPropsGroupVersionKind{
				Group:   crdv1.Spec.Group,
				Version: version.Name,
				Kind:    crdv1.Spec.Names.Kind,
			}
			oeValidator.schemaCache[kindDef.String()] = schema
		}
	case crdv1beta1ApiVersionKind:
		crdv1beta1 := apiextensionsv1beta1.CustomResourceDefinition{}
		err := convertObject(crdResource, &crdv1beta1)
		if err != nil {
			return err
		}
		defaultSchema := openapi3.NewObjectSchema()
		if crdv1beta1.Spec.Validation != nil && crdv1beta1.Spec.Validation.OpenAPIV3Schema != nil {
			defaultSchema, err = getSchema(crdv1beta1.Spec.Validation.OpenAPIV3Schema)
			if err != nil {
				return err
			}
		}
		for _, version := range crdv1beta1.Spec.Versions {
			schema := defaultSchema
			if version.Schema != nil && version.Schema.OpenAPIV3Schema != nil {
				schema, err = getSchema(version.Schema.OpenAPIV3Schema)
				if err != nil {
					return err
				}
			}
			kindDef := extPropsGroupVersionKind{
				Group:   crdv1beta1.Spec.Group,
				Version: version.Name,
				Kind:    crdv1beta1.Spec.Names.Kind,
			}
			oeValidator.schemaCache[kindDef.String()] = schema
		}
	default:
		return fmt.Errorf("Invalid CRD Kind: %s", apiVersionKind)
	}
	return nil
}

func getSchema(jsonSchemaProps interface{}) (*openapi3.Schema, error) {
	schema := openapi3.NewSchema()
	return schema, convertObject(jsonSchemaProps, schema)
}

// convertObject transforms an object into another type by json marshaling/unmarshaling of its properties
func convertObject(source interface{}, target interface{}) error {
	jsonBytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, target)
	if err != nil {
		return err
	}
	return nil
}

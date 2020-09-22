package cmd

import (
	"os"
	"testing"

	"github.com/fllaca/scheriff/pkg/validate"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	wd, _ := os.Getwd()
	t.Log(wd)
	tests := []struct {
		name             string
		opts             validateOptions
		expectedResults  []validate.ValidationResult
		expectedExitCode int
	}{
		{
			name: "test valid deploy",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/deployment_valid.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "v1/Service"},
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "extensions/v1beta1/Deployment"},
			},
		},
		{
			name: "test invalid deployment",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/deployment_invalid.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"valid", "OK", "some-app-envoy", "example", "v1/Service"},
				{"Error at \"/spec/template/spec/containers/0/name\":Property 'name' is missing", "ERROR", "some-app-envoy", "example", "extensions/v1beta1/Deployment"},
			},
		},
		{
			name: "test invalid yaml",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/invalid_yaml.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"Error parsing k8s resource from document 0: error converting YAML to JSON: yaml: line 3: mapping values are not allowed in this context\n", validate.SeverityError, "", "", ""},
			},
		},
		{
			name: "test folder recursive",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/test_recursive"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
				recursive:             true,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "example-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "v1/Service"},
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
			},
		},
		{
			name: "test folder not recursive",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/test_recursive"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "example-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "v1/Service"},
			},
		},
		{
			name: "test valid crd",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/certificate_valid.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "example-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
			},
		},
		{
			name: "test invalid crd",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/certificate_invalid.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"Error at \"/spec/secretName\":Property 'secretName' is missing", validate.SeverityError, "example-invalid-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
			},
		},
		{
			name: "test unknown crd",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/unknown_kind.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
			},
		},
		{
			name: "test unknown crd strict",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/unknown_kind.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                true,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
			},
		},
		{
			name: "test unknown crd then invalid yaml",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/unknown_kind.yaml", "testdata/manifests/invalid_yaml.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
				{"Error parsing k8s resource from document 0: error converting YAML to JSON: yaml: line 3: mapping values are not allowed in this context\n", validate.SeverityError, "", "", ""},
			},
		},
		{
			name: "test unknown crd then invalid yaml strict",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/unknown_kind.yaml", "testdata/manifests/invalid_yaml.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                true,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
				{"Error parsing k8s resource from document 0: error converting YAML to JSON: yaml: line 3: mapping values are not allowed in this context\n", validate.SeverityError, "", "", ""},
			},
		},
		{
			name: "test unknown crd then invalid yaml in same file",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/warn_error.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
				{"Error parsing k8s resource from document 1: error converting YAML to JSON: yaml: line 2: mapping values are not allowed in this context\n", validate.SeverityError, "", "", ""},
			},
		},
		{
			name: "test unknown crd then invalid yaml in same file strict",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/warn_error.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
				strict:                true,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
				{"Error parsing k8s resource from document 1: error converting YAML to JSON: yaml: line 2: mapping values are not allowed in this context\n", validate.SeverityError, "", "", ""},
			},
		},
		{
			name: "test schema with invalid json",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests"},
				openApiSchemaFilename: "testdata/schemas/invalid-json.json",
				crds:                  []string{},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test valid crd v1",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/v1_crontab.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "my-new-cron-object", "", "stable.example.com/v1/CronTab"},
			},
		},
		{
			name: "test crd v1beta1 without schemas",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/without_schemas.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "my-new-cron-object", "", "stable.example.com/v1/CronTab"},
			},
		},
		{
			name: "test crd v1beta1 without default 'spec.validation.schema'",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/crontab_without_default_val.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "my-new-cron-object", "", "stable.example.com/v1/CronTab"},
			},
		},
		{
			name: "test unknown Kind in crd definition",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/unknown_crd_kind.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test invalid crd definition",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/invalid_crd.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test crd with bad yaml syntax",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/invalid_yaml.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test non-existing manifests folder",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/doesnotexist"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/invalid_crd.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test non-existing schema file",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/doesnotexist.json",
				crds:                  []string{"testdata/crds/invalid_crd.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test non-existing crd file",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/doesnotexist.yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test crd invalid yaml",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/crd_v1_crontab.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{"testdata/crds/invalid_yaml"},
				recursive:             false,
				strict:                false,
			},
			expectedExitCode: 1,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name: "test additional properties",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/deployment_additional_properties.yaml", "testdata/manifests/cm_managed_fields.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
			},
			expectedExitCode: 1,
			expectedResults: []validate.ValidationResult{
				{Message: "valid", Severity: "OK", Name: "test", Namespace: "default", Kind: "v1/ConfigMap"},
				{Message: "Error at \"/spec/template/spec\":Property 'unexpectedAdditionalProperty' is unsupported", Severity: "ERROR", Name: "some-app-envoy", Namespace: "example", Kind: "apps/v1/Deployment"},
				{Message: "valid", Severity: "OK", Name: "test-cm", Namespace: "default", Kind: "v1/ConfigMap"},
			},
		},
		{
			name: "test nullable properties",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/ns_nullable_field.yaml"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{Message: "valid", Severity: "OK", Name: "test", Namespace: "", Kind: "v1/Namespace"},
			},
		},
		{
			name: "test validating input files without yaml/yml extension",
			opts: validateOptions{
				filenames:             []string{"testdata/manifests/non_yaml_extension.txt"},
				openApiSchemaFilename: "testdata/schemas/k8s-1.17.0.json",
				crds:                  []string{},
				recursive:             false,
			},
			expectedExitCode: 0,
			expectedResults: []validate.ValidationResult{
				{Message: "valid", Severity: "OK", Name: "test-cm", Namespace: "default", Kind: "v1/ConfigMap"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exitCode, results := runValidate(test.opts)
			assert.Equal(t, test.expectedExitCode, exitCode)
			assert.Equal(t, test.expectedResults, results)
		})
	}
}

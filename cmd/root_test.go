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
		filenames        []string
		recursive        bool
		schema           string
		crds             []string
		expectedResults  []validate.ValidationResult
		expectedExitCode int
	}{
		{
			name:             "test valid deploy",
			filenames:        []string{"testdata/manifests/deployment_valid.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{},
			expectedExitCode: 0,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "v1/Service"},
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "extensions/v1beta1/Deployment"},
			},
		},
		{
			name:             "test invalid deployment",
			filenames:        []string{"testdata/manifests/deployment_invalid.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"valid", "OK", "some-app-envoy", "example", "v1/Service"},
				{"Error at \"/spec/template/spec/containers/0/name\":Property 'name' is missing", "ERROR", "some-app-envoy", "example", "extensions/v1beta1/Deployment"},
			},
		},
		{
			name:             "test invalid yaml",
			filenames:        []string{"testdata/manifests/invalid_yaml.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"Error parsing k8s resource from document 0: error converting YAML to JSON: yaml: line 3: mapping values are not allowed in this context\n", validate.SeverityError, "", "", ""},
			},
		},
		{
			name:             "test folder recursive",
			filenames:        []string{"testdata/manifests/test_recursive"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
			expectedExitCode: 0,
			recursive:        true,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "example-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "v1/Service"},
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
			},
		},
		{
			name:             "test folder not recursive",
			filenames:        []string{"testdata/manifests/test_recursive"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
			expectedExitCode: 0,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "example-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
				{"valid", validate.SeverityOK, "some-app-envoy", "example", "v1/Service"},
			},
		},
		{
			name:             "test valid crd",
			filenames:        []string{"testdata/manifests/certificate_valid.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
			expectedExitCode: 0,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "example-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
			},
		},
		{
			name:             "test invalid crd",
			filenames:        []string{"testdata/manifests/certificate_invalid.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/cert-manager-legacy-v0.15.0.crds.yaml"},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"Error at \"/spec/secretName\":Property 'secretName' is missing", validate.SeverityError, "example-invalid-cert", "example", "cert-manager.io/v1alpha2/Certificate"},
			},
		},
		{
			name:             "test unknown crd",
			filenames:        []string{"testdata/manifests/unknown_kind.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{},
			expectedExitCode: 0,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"Kind 'example.io/v1/UnknownCRD' not found in schema", validate.SeverityWarning, "example-unknown-kind", "example", "example.io/v1/UnknownCRD"},
			},
		},
		{
			name:             "test schema with invalid json",
			filenames:        []string{"testdata/manifests"},
			schema:           "testdata/schemas/invalid-json.json",
			crds:             []string{},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name:             "test valid crd v1",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/v1_crontab.yaml"},
			expectedExitCode: 0,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "my-new-cron-object", "", "stable.example.com/v1/CronTab"},
			},
		},
		{
			name:             "test crd v1beta1 without schemas",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/without_schemas.yaml"},
			expectedExitCode: 0,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "my-new-cron-object", "", "stable.example.com/v1/CronTab"},
			},
		},
		{
			name:             "test crd v1beta1 without default 'spec.validation.schema'",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/crontab_without_default_val.yaml"},
			expectedExitCode: 0,
			recursive:        false,
			expectedResults: []validate.ValidationResult{
				{"valid", validate.SeverityOK, "my-new-cron-object", "", "stable.example.com/v1/CronTab"},
			},
		},
		{
			name:             "test unknown Kind in crd definition",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/unknown_crd_kind.yaml"},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name:             "test invalid crd definition",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/invalid_crd.yaml"},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name:             "test non-existing manifests folder",
			filenames:        []string{"testdata/manifests/doesnotexist"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/invalid_crd.yaml"},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name:             "test non-existing schema file",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/doesnotexist.json",
			crds:             []string{"testdata/crds/invalid_crd.yaml"},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name:             "test non-existing crd file",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/doesnotexist.yaml"},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults:  []validate.ValidationResult{},
		},
		{
			name:             "test crd invalid yaml",
			filenames:        []string{"testdata/manifests/crd_v1_crontab.yaml"},
			schema:           "testdata/schemas/k8s-1.17.0.json",
			crds:             []string{"testdata/crds/invalid_yaml"},
			expectedExitCode: 1,
			recursive:        false,
			expectedResults:  []validate.ValidationResult{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exitCode, results := runValidate(test.filenames, test.schema, test.crds, test.recursive)
			assert.Equal(t, test.expectedExitCode, exitCode)
			assert.Equal(t, test.expectedResults, results)
		})
	}
}

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/fllaca/okay/pkg/fs"
	"github.com/fllaca/okay/pkg/kubernetes"
	"github.com/fllaca/okay/pkg/utils"
	"github.com/fllaca/okay/pkg/validate"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "okay",
		Short: "A Kubernetes manifests validator tool",
		Long:  `A Kubernetes manifests validator tool`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			runValidate(filenames, openApiSchema)
		},
	}
	filenames     = make([]string, 0)
	openApiSchema = ""
	recursive     = false
)

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&filenames, "filename", "f", []string{"."}, "that contains the configuration to be validated")
	rootCmd.PersistentFlags().StringVarP(&openApiSchema, "schema", "s", "", "Kubernetes OpenAPI schema")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "R", false, "Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	rootCmd.MarkPersistentFlagRequired("filename")
	rootCmd.MarkPersistentFlagRequired("schema")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func runValidate(filenames []string, schema string) {
	fmt.Printf("Validating config at %v with schema at %s\n", filenames, openApiSchema)

	opeanApi2SpecsBytes, err := ioutil.ReadFile(openApiSchema)
	if err != nil {
		fmt.Printf("Error loading specs from %s: %s\n", openApiSchema, err)
	}
	validator, err := validate.NewOpenApi2Validator(opeanApi2SpecsBytes)
	if err != nil {
		fmt.Printf("Error loading specs from %s: %s\n", openApiSchema, err)
	}
	// TODO implement extraSchemas from CRDs

	for _, filename := range filenames {
		err := fs.ApplyToPathWithFilter(filename, recursive, func(file string) error {
			fmt.Printf("Validating manifests in %s\n", file)
			fileBytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file, err)
				// continue processing other files in input
				return nil
			}

			documentsBytes := bytes.Split(fileBytes, []byte("\n---\n"))
			for docIndex, documentBytes := range documentsBytes {
				k8sResource, err := kubernetes.ParseResource(documentBytes)
				if err != nil {
					fmt.Printf("\t - Error parsing k8s resource from document %d of %s: %s\n", docIndex, file, err)
					continue
				}
				if len(k8sResource) == 0 {
					continue
				}
				result := validator.Validate(k8sResource)
				// TODO provide documentIndex in output
				fmt.Printf("\t - %s, %s (%s): %s\n", result.Severity, utils.JoinNotEmptyStrings("/", result.Namespace, result.Name), result.Kind, result.Message)
			}
			return nil

		}, fs.IsYamlFilter)
		if err != nil {
			fmt.Printf("Error while validating %s: %s\n", filename, err)
		}
	}
}

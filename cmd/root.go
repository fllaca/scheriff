package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/fllaca/okay/pkg/fs"
	"github.com/fllaca/okay/pkg/utils"
	"github.com/fllaca/okay/pkg/validate"
	"github.com/spf13/cobra"

	"github.com/gookit/color"
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
		return
	}
	resourceValidator, err := validate.NewOpenApi2Validator(opeanApi2SpecsBytes)
	if err != nil {
		fmt.Printf("Error loading specs from %s: %s\n", openApiSchema, err)
		return
	}
	// TODO implement extraSchemas from CRDs

	fileValidator := validate.NewYamlFileValidator(resourceValidator)

	for _, filename := range filenames {
		err := fs.ApplyToPathWithFilter(filename, recursive, func(file string) error {
			fmt.Printf("Validating manifests in %s\n", file)
			fileBytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file, err)
				// continue processing other files in input
				return nil
			}

			validationResults := fileValidator.Validate(fileBytes)
			outputResult(validationResults)

			return nil

		}, fs.IsYamlFilter)
		if err != nil {
			fmt.Printf("Error while validating %s: %s\n", filename, err)
		}
	}
}

func outputResult(results []validate.ValidationResult) {
	for _, result := range results {
		fmt.Printf("\t - %s, %s (%s): %s\n", colorSeverity(result.Severity), utils.JoinNotEmptyStrings("/", result.Namespace, result.Name), result.Kind, result.Message)
	}
}

func colorSeverity(severity validate.Severity) string {
	red := color.FgRed.Render
	green := color.FgGreen.Render
	yellow := color.FgYellow.Render
	switch severity {
	case validate.SeverityError:
		return red(severity)
	case validate.SeverityWarning:
		return yellow(severity)
	case validate.SeverityOK:
		return green(severity)
	default:
		return (string)(severity)
	}
}

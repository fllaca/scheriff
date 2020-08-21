package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fllaca/scheriff/pkg/fs"
	"github.com/fllaca/scheriff/pkg/kubernetes"
	"github.com/fllaca/scheriff/pkg/utils"
	"github.com/fllaca/scheriff/pkg/validate"
	"github.com/spf13/cobra"

	"github.com/gookit/color"
)

var (
	Version = "development"
	rootCmd = &cobra.Command{
		Use:   "scheriff",
		Short: "Schema Sheriff: A Kubernetes manifests validator tool",
		Long: `Schema Sheriff: A Kubernetes manifests validator tool

Schema Sheriff performs offline validation of Kubernetes configuration manifests by checking them against OpenApi schemas. No connectivity to the Kubernetes cluster is needed`,
		Run: func(cmd *cobra.Command, args []string) {
			exitCode, _ := runValidate(filenames, openApiSchema, crds, recursive)
			os.Exit(exitCode)
		},
	}
	filenames     = make([]string, 0)
	crds          = make([]string, 0)
	openApiSchema = ""
	recursive     = false
	strict        = false
)

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&filenames, "filename", "f", []string{}, "(required) file or directories that contain the configuration to be validated")
	// TODO support OpenApi V3 input
	rootCmd.PersistentFlags().StringVarP(&openApiSchema, "schema", "s", "", "(required) Kubernetes OpenAPI V2 schema to validate against")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "R", false, "process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	rootCmd.PersistentFlags().StringArrayVarP(&crds, "crd", "c", []string{}, "files or directories that contain CustomResourceDefinitions to be used for validation")
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "S", false, "return exit code 1 not only on errors but also when warnings are encountered.")
	rootCmd.MarkPersistentFlagRequired("filename")
	rootCmd.MarkPersistentFlagRequired("schema")
}

// Execute executes the root command.
func Execute(version, date, commit string) error {
	rootCmd.Version = fmt.Sprintf("%s - %s (built %s)", version, date, commit)
	return rootCmd.Execute()
}

func runValidate(filenames []string, schema string, crds []string, recursive bool) (int, []validate.ValidationResult) {
	totalResults := make([]validate.ValidationResult, 0)
	fmt.Printf("Validating config in %s against schema in %s\n", utils.JoinNotEmptyStrings(", ", filenames...), openApiSchema)
	exitCode := 0

	opeanApi2SpecsBytes, err := ioutil.ReadFile(schema)
	if err != nil {
		fmt.Printf("Error loading specs from %s: %s\n", schema, err)
		return 1, totalResults
	}
	resourceValidator, err := validate.NewOpenApi2Validator(opeanApi2SpecsBytes)
	if err != nil {
		fmt.Printf("Error loading specs from %s: %s\n", openApiSchema, err)
		return 1, totalResults
	}

	for _, crd := range crds {
		err := fs.ApplyToPathWithFilter(crd, false, func(file string) error {
			fmt.Printf("Using CustomResourceDefinitions from %s\n", file)
			fileBytes, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			crdResources, err := kubernetes.ParseResourcesFromYaml(fileBytes)
			if err != nil {
				return err
			}
			for _, crdResource := range crdResources {
				err = resourceValidator.AddCrdSchemas(crdResource)
				if err != nil {
					return err
				}
			}
			return nil
		}, fs.IsYamlFilter)
		if err != nil {
			fmt.Printf("Error loading CustomResourceDefinitions from %s: %s\n", crd, err)
			// TODO: log warning instead?
			return 1, totalResults
		}
	}

	fileValidator := validate.NewYamlFileValidator(resourceValidator)

	fmt.Println("Results:")
	for _, filename := range filenames {
		err := fs.ApplyToPathWithFilter(filename, recursive, func(file string) error {
			fmt.Printf("Validating manifests in %s:\n", file)
			fileBytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file, err)
				// continue processing other files in input
				return nil
			}

			validationResults := fileValidator.Validate(fileBytes)
			outputResult(validationResults)
			if containsSeverity(validationResults, strict) {
				exitCode = 1
			}
			totalResults = append(totalResults, validationResults...)
			return nil

		}, fs.IsYamlFilter)
		if err != nil {
			fmt.Printf("Error while validating %s: %s\n", filename, err)
			exitCode = 1
		}
	}
	return exitCode, totalResults
}

func outputResult(results []validate.ValidationResult) {
	for _, result := range results {
		fmt.Printf("\t - %s, %s (%s): %s\n", colorSeverity(result.Severity), utils.JoinNotEmptyStrings("/", result.Namespace, result.Name), result.Kind, result.Message)
	}
	fmt.Println()
}

func containsSeverity(results []validate.ValidationResult, strict bool) bool {
	for _, result := range results {
		switch result.Severity {
		case validate.SeverityError:
			return true
		case validate.SeverityWarning:
			return strict
		}
	}
	return false
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

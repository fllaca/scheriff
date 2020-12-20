package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/fllaca/scheriff/pkg/fs"
	"github.com/fllaca/scheriff/pkg/kubernetes"
	"github.com/fllaca/scheriff/pkg/utils"
	"github.com/fllaca/scheriff/pkg/validate"
	"github.com/spf13/cobra"

	"github.com/gookit/color"
)

type validateOptions struct {
	filenames             []string
	crds                  []string
	openApiSchemaFilename string
	recursive             bool
	strict                bool
	input                 io.Reader
}

var (
	Version = "development"
	rootCmd = &cobra.Command{
		Use:   "scheriff",
		Short: "Schema Sheriff: A Kubernetes manifests validator tool",
		Long: `Schema Sheriff: A Kubernetes manifests validator tool

Schema Sheriff performs offline validation of Kubernetes configuration manifests by checking them against OpenApi schemas. No connectivity to the Kubernetes cluster is needed`,
		Run: func(cmd *cobra.Command, args []string) {
			options.input = cmd.InOrStdin()
			exitCode, _ := runValidate(options)
			os.Exit(exitCode)
		},
	}

	options = validateOptions{}
)

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&options.filenames, "filename", "f", []string{}, "(required) file or directories that contain the configuration to be validated")
	// TODO support OpenApi V3 input
	rootCmd.PersistentFlags().StringVarP(&options.openApiSchemaFilename, "schema", "s", "", "(required) Kubernetes OpenAPI V2 schema to validate against")
	rootCmd.PersistentFlags().BoolVarP(&options.recursive, "recursive", "R", false, "process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	rootCmd.PersistentFlags().StringArrayVarP(&options.crds, "crd", "c", []string{}, "files or directories that contain CustomResourceDefinitions to be used for validation")
	rootCmd.PersistentFlags().BoolVarP(&options.strict, "strict", "S", false, "return exit code 1 not only on errors but also when warnings are encountered.")
	rootCmd.MarkPersistentFlagRequired("filename")
	rootCmd.MarkPersistentFlagRequired("schema")
}

// Execute executes the root command.
func Execute(version, date, commit string) error {
	rootCmd.Version = fmt.Sprintf("%s - %s (built %s)", version, date, commit)
	return rootCmd.Execute()
}

func runValidate(opts validateOptions) (int, []validate.ValidationResult) {
	totalResults := make([]validate.ValidationResult, 0)
	fmt.Printf("Validating config in %s against schema in %s\n", utils.JoinNotEmptyStrings(", ", opts.filenames...), opts.openApiSchemaFilename)
	exitCode := 0

	opeanApi2SpecsBytes, err := ioutil.ReadFile(opts.openApiSchemaFilename)
	if err != nil {
		fmt.Printf("Error loading specs from %s: %s\n", opts.filenames, err)
		return 1, totalResults
	}
	resourceValidator, err := validate.NewOpenApi2Validator(opeanApi2SpecsBytes)
	if err != nil {
		fmt.Printf("Error loading specs from %s: %s\n", opts.openApiSchemaFilename, err)
		return 1, totalResults
	}

	for _, crd := range opts.crds {
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
	for _, filename := range opts.filenames {
		// case stdin:
		if filename == "-" {
			fileBytes, err := ioutil.ReadAll(opts.input)
			if err != nil {
				fmt.Printf("Error reading stdin: %s\n", err)
				return 1, totalResults
			}
			validationResults := fileValidator.Validate(fileBytes)
			outputResult(validationResults)
			totalResults = append(totalResults, validationResults...)
			continue
		}

		err := fs.ApplyToPathWithFilter(filename, opts.recursive, func(file string) error {
			fmt.Printf("Validating manifests in %s:\n", file)
			fileBytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file, err)
				// continue processing other files in input
				return nil
			}

			validationResults := fileValidator.Validate(fileBytes)
			outputResult(validationResults)
			totalResults = append(totalResults, validationResults...)
			return nil

		}, fs.IsYamlFilter)
		if err != nil {
			fmt.Printf("Error while validating %s: %s\n", filename, err)
			exitCode = 1
		}
	}
	if containsSeverity(totalResults, opts.strict) {
		exitCode = 1
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
			if strict {
				return true
			}
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

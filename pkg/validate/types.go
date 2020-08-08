package validate

type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARN"
	SeverityOK      Severity = "OK"
)

type ResourceValidator interface {
	Validate(resource map[string]interface{}) ValidationResult
}

type ValidationResult struct {
	// Message holds a brief description of the validation result
	Message string
	// Severity specifies if the validation is OK/ERROR/WARNING
	Severity Severity
	// Name of the validated resourcce
	Name string
	// Namespace of the validated resourcce
	Namespace string
	// Kind of the validated resourcce
	Kind string
}

type FileValidator interface {
	Validate(data []byte) []ValidationResult
}

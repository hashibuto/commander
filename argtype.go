package commander

type ArgType string

const (
	ArgTypeUnspecified ArgType = ""
	ArgTypeInt         ArgType = "INT"
	ArgTypeFloat       ArgType = "FLOAT"
	ArgTypeString      ArgType = "STRING"
	ArgTypeBool        ArgType = "BOOL"
)

// InferArgType infers the ArgType based on the variable type, or returns ArgTypeUnspecified
// if the variable type does not match one of the known values.
func InferArgType(value any) ArgType {
	switch value.(type) {
	case int:
		return ArgTypeInt
	case float64:
		return ArgTypeFloat
	case string:
		return ArgTypeString
	case bool:
		return ArgTypeBool
	default:
		return ArgTypeUnspecified
	}
}

package commander

import "reflect"

type ArgType string

const (
	ArgTypeUnspecified ArgType = ""
	ArgTypeInt         ArgType = "INT"
	ArgTypeFloat       ArgType = "FLOAT"
	ArgTypeString      ArgType = "STRING"
	ArgTypeBool        ArgType = "BOOL"
)

var (
	intType    reflect.Type = reflect.TypeOf(int(1))
	floatType  reflect.Type = reflect.TypeOf(float64(1))
	stringType reflect.Type = reflect.TypeOf("string")
	boolType   reflect.Type = reflect.TypeOf(false)
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
		baseType := reflect.TypeOf(value)
		if baseType.ConvertibleTo(intType) {
			return ArgTypeInt
		}

		if baseType.ConvertibleTo(floatType) {
			return ArgTypeFloat
		}

		if baseType.ConvertibleTo(stringType) {
			return ArgTypeString
		}

		if baseType.ConvertibleTo(boolType) {
			return ArgTypeBool
		}

		return ArgTypeUnspecified
	}
}

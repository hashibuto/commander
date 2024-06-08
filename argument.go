package commander

import "fmt"

// an argument represents a positional argument.  it is non-defaultable
type Argument struct {
	Name          string
	Description   string
	ArgType       ArgType
	AllowMultiple bool  // if enabled, will be returned as an array of ArgType
	OneOf         []any // if specified, value must belong to collection
	Completer     Completer
}

// Validate returns an error if any part of the argument is invalid
func (a *Argument) Validate() error {
	if len(a.Name) == 0 {
		return fmt.Errorf("argument must have a non-empty name")
	}

	if a.ArgType == ArgTypeUnspecified {
		// Default to string
		a.ArgType = ArgTypeString
	}

	for _, oneOf := range a.OneOf {
		inferredType := InferArgType(oneOf)
		if inferredType != a.ArgType {
			return fmt.Errorf("value in OneOf \"%v\" did not match the argument type \"%s\"", oneOf, a.ArgType)
		}
	}

	return nil
}

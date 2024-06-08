package commander

import "fmt"

type Flag struct {
	Name          string
	ShortName     string
	Description   string
	ArgType       ArgType
	AllowMultiple bool // if enabled, will be returned as an array of ArgType
	DefaultValue  any
	OneOf         []any // if specified, value must belong to collection
	Completer     Completer
	IsRequired    bool
}

// Validate returns an error if any part of the flag is invalid
func (f *Flag) Validate() error {
	if len(f.Name) == 0 && len(f.ShortName) == 0 {
		return fmt.Errorf("flag must specify at least one of a name or short name")
	}

	if len(f.Name) <= 1 {
		return fmt.Errorf("flag name cannot be shorter than 2 characters")
	}

	if len(f.ShortName) > 1 {
		return fmt.Errorf("short name must be exactly 1 character long")
	}

	if f.ArgType == ArgTypeUnspecified {
		// Default to bool
		f.ArgType = ArgTypeBool
	}

	if f.DefaultValue == nil && f.ArgType == ArgTypeBool {
		f.DefaultValue = false
	}

	if f.OneOf != nil && f.ArgType == ArgTypeBool {
		return fmt.Errorf("OneOf is not compatible with boolean flags")
	}

	if f.Completer != nil && f.ArgType == ArgTypeBool {
		return fmt.Errorf("Completer is not compatible with boolean flags")
	}

	if f.AllowMultiple && f.ArgType == ArgTypeBool {
		return fmt.Errorf("AllowMultiple is not compatible with boolean flags")
	}

	for _, oneOf := range f.OneOf {
		inferredType := InferArgType(oneOf)
		if inferredType != f.ArgType {
			return fmt.Errorf("value in OneOf \"%v\" did not match the argument type \"%s\"", oneOf, f.ArgType)
		}
	}

	return nil
}

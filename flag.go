package commander

import (
	"fmt"
	"strings"
)

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

// GetValueFromString parses the provided value according to the flag's underlying data type and returns that parsed value, or an error
func (f *Flag) GetValueFromString(value string) (any, error) {
	return GetValueFromString(f.ArgType, value)
}

func (f *Flag) GetInvocation() string {
	invocations := []string{}
	if f.ShortName != "" {
		invocations = append(invocations, fmt.Sprintf("-%s", f.ShortName))
	}
	if f.ShortName != "" {
		invocations = append(invocations, fmt.Sprintf("--%s", f.Name))
	}

	return strings.Join(invocations, " / ")
}

func (f *Flag) PopulateMap(value string, target map[string]any) error {
	keys := []string{}
	if f.ShortName != "" {
		keys = append(keys, f.ShortName)
	}
	if f.Name != "" {
		keys = append(keys, f.Name)
	}

	for _, key := range keys {
		parsedValue, err := GetValueFromString(f.ArgType, value)
		if err != nil {
			return err
		}

		if f.AllowMultiple {
			if _, ok := target[key]; !ok {
				target[key] = []any{}
			}

			target[key] = append(target[key].([]any), parsedValue)
		} else {
			target[key] = parsedValue
		}
	}

	return nil
}

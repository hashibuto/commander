package commander

import (
	"fmt"
	"strings"

	ns "github.com/hashibuto/nilshell"
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
		return fmt.Errorf("OneOf is not compatible with boolean flags in %s", f.GetInvocation())
	}

	if f.Completer != nil && f.ArgType == ArgTypeBool {
		return fmt.Errorf("Completer is not compatible with boolean flags in %s", f.GetInvocation())
	}

	if f.AllowMultiple && f.ArgType == ArgTypeBool {
		return fmt.Errorf("AllowMultiple is not compatible with boolean flags in %s", f.GetInvocation())
	}

	if f.AllowMultiple && f.DefaultValue != nil {
		if _, ok := f.DefaultValue.([]any); !ok {
			return fmt.Errorf("DefaultValue must be a []any when AllowMultiple is true in %s", f.GetInvocation())
		}
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
	if f.Name != "" {
		invocations = append(invocations, fmt.Sprintf("--%s", f.Name))
	}

	return strings.Join(invocations, " / ")
}

func (f *Flag) GetPaddedInvocation() string {
	invocation := f.GetInvocation()
	if strings.HasPrefix(invocation, "--") {
		return fmt.Sprintf("     %s", invocation)
	}

	return invocation
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

		if f.OneOf != nil {
			if !MatchesOneOf(f.OneOf, parsedValue) {
				return fmt.Errorf("\"%s\" does not belong to the collection defined by the flag", parsedValue)
			}
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

func (f *Flag) PopulateDefault(target map[string]any) error {
	keys := []string{}
	if f.ShortName != "" {
		keys = append(keys, f.ShortName)
	}
	if f.Name != "" {
		keys = append(keys, f.Name)
	}

	for _, key := range keys {
		// Skip anything already populated
		if _, exists := target[key]; exists {
			continue
		}

		if f.DefaultValue == nil && f.IsRequired {
			return fmt.Errorf("flag %s is required", f.GetInvocation())
		}

		target[key] = f.DefaultValue
	}

	return nil
}

func (f *Flag) SuggestValues(prefix string) []*ns.AutoComplete {
	if f.OneOf != nil {
		values := []*ns.AutoComplete{}
		for _, oneOf := range f.OneOf {
			oneOfStr := fmt.Sprintf("%s", oneOf)
			if strings.HasPrefix(oneOfStr, prefix) {
				values = append(values, &ns.AutoComplete{
					Value:   oneOfStr,
					Display: oneOfStr,
				})
			}
		}

		return values
	}

	if f.Completer != nil {
		return f.Completer(prefix)
	}

	return nil
}

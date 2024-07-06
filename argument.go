package commander

import (
	"fmt"
	"strings"

	ns "github.com/hashibuto/nilshell"
)

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

// GetValueFromString parses the provided value according to the argument's underlying data type and returns that parsed value, or an error
func (a *Argument) GetValueFromString(value string) (any, error) {
	return GetValueFromString(a.ArgType, value)
}

func (a *Argument) PopulateMap(value string, target map[string]any) error {
	parsedValue, err := GetValueFromString(a.ArgType, value)
	if err != nil {
		return err
	}

	if a.OneOf != nil {
		if !MatchesOneOf(a.OneOf, parsedValue) {
			return fmt.Errorf("\"%s\" does not belong to the collection defined by the argument", parsedValue)
		}
	}

	if a.AllowMultiple {
		if _, ok := target[a.Name]; !ok {
			target[a.Name] = []any{}
		}

		target[a.Name] = append(target[a.Name].([]any), parsedValue)
	} else {
		target[a.Name] = parsedValue
	}

	return nil
}

func (a *Argument) SuggestValues(prefix string) []*ns.AutoComplete {
	if a.OneOf != nil {
		values := []*ns.AutoComplete{}
		for _, oneOf := range a.OneOf {
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

	if a.Completer != nil {
		return a.Completer(prefix)
	}

	return nil
}

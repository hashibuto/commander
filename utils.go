package commander

import (
	"fmt"
	"strconv"
)

func GetValueFromString(argType ArgType, value string) (any, error) {
	switch argType {
	case ArgTypeInt:
		iVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("value could not be parsed to an integer")
		}

		return iVal, nil
	case ArgTypeFloat:
		fVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("value could not be parsed to a float")
		}
		return fVal, nil
	case ArgTypeBool:
		if value == "true" || value == "t" || value == "1" {
			return true, nil
		}
		if value == "false" || value == "f" || value == "0" {
			return false, nil
		}
		return nil, fmt.Errorf("value could not be parsed into a bool")
	case ArgTypeString:
		return value, nil
	}

	return nil, fmt.Errorf("unknown arg type")
}

func MatchesOneOf(oneOf []any, sample any) bool {
	for _, one := range oneOf {
		if fmt.Sprintf("%s", one) == fmt.Sprintf("%s", sample) {
			return true
		}
	}

	return false
}

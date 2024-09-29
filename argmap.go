package commander

type ArgMap map[string]any

func (m ArgMap) GetString(argName string) string {
	v := m[argName]
	if v == nil {
		return ""
	}

	return v.(string)
}

func (m ArgMap) GetInt(argName string) int {
	v := m[argName]
	if v == nil {
		return 0
	}

	return v.(int)
}

func (m ArgMap) GetFloat(argName string) float64 {
	v := m[argName]
	if v == nil {
		return 0.
	}

	return v.(float64)
}

func (m ArgMap) GetBool(argName string) bool {
	v := m[argName]
	if v == nil {
		return false
	}

	return v.(bool)
}

func (m ArgMap) GetStringArray(argName string) []string {
	v := m[argName]
	if v == nil {
		return []string{}
	}

	anyArr := v.([]any)
	arr := make([]string, len(anyArr))

	for idx, item := range anyArr {
		arr[idx] = item.(string)
	}

	return arr
}

func (m ArgMap) GetIntArray(argName string) []int {
	v := m[argName]
	if v == nil {
		return []int{}
	}

	anyArr := v.([]any)
	arr := make([]int, len(anyArr))

	for idx, item := range anyArr {
		arr[idx] = item.(int)
	}

	return arr
}

func (m ArgMap) GetFloatArray(argName string) []float64 {
	v := m[argName]
	if v == nil {
		return []float64{}
	}

	anyArr := v.([]any)
	arr := make([]float64, len(anyArr))

	for idx, item := range anyArr {
		arr[idx] = item.(float64)
	}

	return arr
}

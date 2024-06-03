package commander

type ArgType string

const (
	ArgTypeUnspecified ArgType = ""
	ArgTypeInt         ArgType = "INT"
	ArgTypeFloat       ArgType = "FLOAT"
	ArgTypeString      ArgType = "STRING"
	ArgTypeBool        ArgType = "BOOL"
)

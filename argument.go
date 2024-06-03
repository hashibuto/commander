package commander

type Argument struct {
	Name          string
	Description   string
	ArgType       ArgType
	AllowMultiple bool // if enabled, will be returned as an array of ArgType
	DefaultValue  any
	OneOf         []any // if specified, value must belong to collection
	Completer     Completer
}

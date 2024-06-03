package commander

// Completer is a type of function which returns a list of strings based on a search string
type Completer func(search string) []string

package commander

import ns "github.com/hashibuto/nilshell"

// Completer is a type of function which returns a list of strings based on a search string
type Completer func(search string) []*ns.AutoComplete

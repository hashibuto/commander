package commander

type Command struct {
	Name        string
	Description string
	Flags       []*Flag
	Arguments   []*Argument
}

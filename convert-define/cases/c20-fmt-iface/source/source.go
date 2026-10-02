package source

type Name string

func (n Name) String() string { return string(n) }

type A struct {
	S Name
	M map[string]int
}

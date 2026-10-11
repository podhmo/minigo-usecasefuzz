package app

type Alpha struct {
	B *Beta `json:"b"`
}

type Beta struct {
	A *Alpha `json:"a"`
}

type Kind int

const (
	KindA Kind = iota
	KindB
)

var _ = 0

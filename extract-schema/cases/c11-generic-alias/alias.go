package app

type Box[T any] struct {
	V T `json:"v"`
}

type IntBox = Box[int]

type Wrap struct {
	B IntBox `json:"b"`
}

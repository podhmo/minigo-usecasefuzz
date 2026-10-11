package app

type Box[T any] struct {
	V T `json:"v"`
}

type Holder struct {
	B Box[int] `json:"b"`
}

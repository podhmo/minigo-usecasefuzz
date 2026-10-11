package app

type Box[T any] struct {
	V T `json:"v"`
}

type Plain struct {
	ID int `json:"id"`
}

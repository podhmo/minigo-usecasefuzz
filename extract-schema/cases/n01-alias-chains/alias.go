package app

type A = B
type B = C
type C = int

type M = map[string]int
type S = []string
type U2 = User

type User struct {
	ID int `json:"id"`
}

type Wrap struct {
	A  A  `json:"a"`
	M  M  `json:"m"`
	S  S  `json:"s"`
	U2 U2 `json:"u2"`
}

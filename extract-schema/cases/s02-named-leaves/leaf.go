package app

type Age int
type Score float64
type Flag bool
type Name string
type Bytes []byte
type Celsius = float64

type Bag struct {
	A Age     `json:"a"`
	S Score   `json:"s"`
	F Flag    `json:"f"`
	N Name    `json:"n"`
	B Bytes   `json:"b"`
	C Celsius `json:"c"`
}

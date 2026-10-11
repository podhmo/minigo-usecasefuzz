package app

type Coll struct {
	Ints    []int                    `json:"ints"`
	Matrix  [][]string               `json:"matrix"`
	ByName  map[string]int           `json:"by_name"`
	ByInt   map[int]string           `json:"by_int"`
	ByNamed map[Key]Val              `json:"by_named"`
	Any     map[string]any           `json:"any"`
	Nested  []map[string][]*Leaf     `json:"nested"`
}

type Key struct {
	Kind string `json:"kind"`
}

type Val struct {
	V int `json:"v"`
}

type Leaf struct {
	ID int `json:"id"`
}

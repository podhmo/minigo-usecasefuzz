package app

type Deep struct {
	L1 struct {
		L2 struct {
			V string `json:"v"`
		} `json:"l2"`
	} `json:"l1"`
	List []struct {
		K string `json:"k"`
	} `json:"list"`
	Table map[string]struct {
		N int `json:"n"`
	} `json:"table"`
	Pair struct {
		A, B int `json:",omitempty"`
	} `json:"pair"`
}

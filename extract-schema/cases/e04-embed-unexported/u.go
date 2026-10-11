package app

type inner struct {
	X int    `json:"x"`
	Y string `json:"y"`
}

type Wrap struct {
	inner
	Own int `json:"own"`
}

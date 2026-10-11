package app

type Base struct {
	ID string `json:"id"`
}

type Outer struct {
	In struct {
		Base
		Extra int `json:"extra"`
	} `json:"in"`
}

package app

type Base struct {
	ID   string `json:"id"`
	Self int    `json:"self"`
}

type Mid struct {
	Base
	M int `json:"m"`
}

type Outer struct {
	Mid
	*PtrBase
	Base `json:"tb"` // tagged embed keeps its own property name
	Own  int `json:"own"`
}

type PtrBase struct {
	PB string `json:"pb"`
}

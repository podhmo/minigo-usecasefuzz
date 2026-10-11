package app

type L1 struct{ V string }
type R1 struct{ V int }

// Clash drops V: both embeds promote V at depth 1.
type Clash struct {
	L1
	R1
}

// Top adds an explicit V at depth 0: the shallow field wins over the
// depth-1 ambiguity.
type Top struct {
	Clash
	V bool `json:"v"`
}

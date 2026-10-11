package app

// shadowing a builtin name is legal inside a package
type string struct {
	V int `json:"v"`
}

// a non-ASCII uppercase name: exported by Go's rule
type Δelta struct {
	ΣX int    `json:"sx"`
	ΔY string `json:"dy"`
}

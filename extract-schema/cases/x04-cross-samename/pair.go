package app

import (
	"example.com/m/a"
	"example.com/m/b"
)

type Pair struct {
	X a.Money `json:"x"`
	Y b.Money `json:"y"`
}

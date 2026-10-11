package app

import (
	"example.com/m/left"
	"example.com/m/right"
)

type Join struct {
	L left.Arm  `json:"l"`
	R right.Arm `json:"r"`
}

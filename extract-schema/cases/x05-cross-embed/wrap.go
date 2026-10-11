package app

import "example.com/m/dep"

type Wrap struct {
	dep.Base
	Own int `json:"own"`
}

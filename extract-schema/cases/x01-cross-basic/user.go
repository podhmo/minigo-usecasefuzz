package app

import "example.com/m/dep"

type User struct {
	Name string   `json:"name"`
	Dep  dep.Meta `json:"dep"`
}

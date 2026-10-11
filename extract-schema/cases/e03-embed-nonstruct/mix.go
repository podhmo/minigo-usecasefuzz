package app

import "time"

type MyInt int

type Mix struct {
	int
	MyInt
	time.Duration
	Name string `json:"name"`
}

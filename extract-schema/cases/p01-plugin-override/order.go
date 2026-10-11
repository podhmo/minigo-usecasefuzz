package app

import "time"

type Money struct {
	Amount int `json:"amount"`
}

type Order struct {
	Price Money     `json:"price"`
	At    time.Time `json:"at"`
}

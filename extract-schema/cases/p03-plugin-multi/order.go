package app

type Money struct {
	Amount int `json:"amount"`
}

type Order struct {
	// Price is the price.
	Price Money `json:"price"`
}

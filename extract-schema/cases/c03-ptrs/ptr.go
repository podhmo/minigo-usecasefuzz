package app

type Wrap struct {
	P   *int     `json:"p"`
	PP  **string `json:"pp"`
	Sub *Sub     `json:"sub"`
	L   *[]int   `json:"l"`
}

type Sub struct {
	V int `json:"v"`
}

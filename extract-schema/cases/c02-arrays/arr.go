package app

type Fixed struct {
	Rgb   [3]int       `json:"rgb"`
	Board [8][8]string `json:"board"`
	Skip  [0]byte      `json:"skip"`
}

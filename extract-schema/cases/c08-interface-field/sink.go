package app

type Sink struct {
	W interface {
		Write([]byte) (int, error)
	} `json:"w"`
}

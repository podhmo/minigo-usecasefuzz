package app

type Bus struct {
	In  chan int    `json:"in"`
	Out chan<- int  `json:"out"`
}

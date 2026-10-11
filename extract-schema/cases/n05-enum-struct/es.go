package app

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

const Origin Point = Point{X: 0, Y: 0}

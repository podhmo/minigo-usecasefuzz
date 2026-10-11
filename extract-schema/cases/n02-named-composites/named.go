package app

type M map[string]int
type S []int
type A [4]int
type P *User
type User struct {
	ID int `json:"id"`
}

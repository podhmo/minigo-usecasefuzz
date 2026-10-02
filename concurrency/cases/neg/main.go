package main

func GoNonFunc() int {
	x := 5
	go x
	return 1
}

func main() {}

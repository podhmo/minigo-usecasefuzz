package main

func main() {
	c := make(chan chan int)
	go func() {
		inner := make(chan int, 1)
		inner <- 7
		c <- inner
	}()
	got := <-c
	println("got", <-got)
}

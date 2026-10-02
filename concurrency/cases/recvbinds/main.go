package main

func RecvAssign() int {
	ch := make(chan int, 1)
	ch <- 4
	var x int
	x = <-ch
	return x // 4
}

func RecvTwoVal() int {
	ch := make(chan int, 1)
	ch <- 9
	v, ok := <-ch
	if !ok {
		return -1
	}
	return v // 9
}

func RecvInCall() int {
	f := func(n int) int { return n * 2 }
	ch := make(chan int, 1)
	ch <- 5
	return f(<-ch) // 10
}

func RecvParen() int {
	ch := make(chan int, 1)
	ch <- 3
	return (<-ch) + 1 // 4
}

func RecvMultipleStmt() int {
	a := make(chan int, 1)
	b := make(chan int, 1)
	a <- 1
	b <- 2
	x, y := <-a, <-b
	return x + y // 3
}

func RecvOfRecv() int {
	inner := make(chan int, 1)
	inner <- 6
	outer := make(chan chan int, 1)
	outer <- inner
	return <-(<-outer) // 6
}

func SendFromRecv() int {
	src := make(chan int, 1)
	dst := make(chan int, 1)
	src <- 7
	go func() { dst <- <-src }()
	return <-dst // 7
}

func main() {}

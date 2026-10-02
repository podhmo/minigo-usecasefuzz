package main

func MakeNegBuf() int {
	ch := make(chan int, -1) // Go: panic makeslicecap
	_ = ch
	return 1
}

// chan of bool/string/named elem
type S string

func ChanNamedElem() int {
	ch := make(chan S, 1)
	ch <- "x"
	return len(<-ch) // want 1
}

// goroutine calling a bound host func that sleeps
func SleepInGo() int {
	ch := make(chan int, 1)
	go func() {
		ch <- 5
	}()
	return <-ch // want 5
}

// defer inside a spawned goroutine running a deferred func that itself panics
func DeferPanicInGo() int {
	ch := make(chan int, 2)
	go func() {
		defer func() { ch <- 1 }()
	}()
	return <-ch // want 1
}

func main() {}

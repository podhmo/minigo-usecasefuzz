package main

func GoNested() int {
	done := make(chan int, 1)
	go func() {
		go func() { done <- 1 }()
	}()
	return <-done // 1
}

func GoNestedPanic() int {
	go func() {
		go func() { panic("deep") }()
		select {}
	}()
	select {}
}

func GoTrapIndex() int {
	go func() { s := []int{1}; _ = s[5] }()
	select {}
}

var initCh = make(chan int, 1)

func init() { go func() { initCh <- 1 }() }

func GoInInit() int { return <-initCh }

type Counter struct {
	n    int
	done chan int
}

func (c *Counter) M() { c.done <- c.n }

func GoBoundMethod() int {
	done := make(chan int, 1)
	c := &Counter{n: 9, done: done}
	go c.M()
	return <-done // 9
}

func Id[T any](x T) T { return x }

func GoGeneric() int {
	go Id(12)
	return 1
}

func GoRecover() int {
	done := make(chan int, 1)
	go func() {
		defer func() {
			if x := recover(); x != nil {
				done <- 7
			} else {
				done <- -1
			}
		}()
		panic("caught")
	}()
	return <-done // 7
}

func GoBuiltin() int {
	go println("spawned")
	return 1
}

// goroutine panic while root is in a blocking select -> run fails
func PanicWhileRootSelects() int {
	ch := make(chan int)
	go func() { panic("during-select") }()
	<-ch
	return 0
}

// goroutine trap (undefined member) while root parked
func TrapWhileRootSelects() int {
	go func() {
		var x struct{}
		_ = x.noSuchMember
	}()
	select {}
}

// loop var capture: Go 1.22 semantics give each iteration its own var
func GoLoopVar() int {
	ch := make(chan int, 3)
	for i := 0; i < 3; i++ {
		go func() { ch <- i }()
	}
	a, b, c := <-ch, <-ch, <-ch
	return a + b + c // 3 regardless of order
}

func GoShorthandMethodValue() int {
	done := make(chan int, 1)
	c := &Counter{n: 2, done: done}
	go func() { c.M() }()
	return <-done // 2
}

func main() {}

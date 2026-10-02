package main

func ChanOfFunc() int {
	ch := make(chan func() int, 1)
	ch <- func() int { return 4 }
	f := <-ch
	return f() // 4
}

type Msg struct{ N int }

func ChanOfStruct() int {
	ch := make(chan Msg, 1)
	ch <- Msg{N: 8}
	m := <-ch
	return m.N // 8
}

func ChanOfPtr() int {
	n := 3
	ch := make(chan *int, 1)
	ch <- &n
	p := <-ch
	return *p // 3
}

func ChanOfAny() int {
	ch := make(chan any, 2)
	ch <- 1
	ch <- "x"
	a := <-ch
	b := <-ch
	i, iok := a.(int)
	s, sok := b.(string)
	if iok && sok && s == "x" {
		return i // 1
	}
	return -1
}

type C chan int

func NamedChanType() int {
	c := make(C)
	go func() { c <- 3 }()
	return <-c // 3
}

func ChanSlice() int {
	chs := []chan int{make(chan int, 1), make(chan int, 1)}
	chs[1] <- 5
	select {
	case v := <-chs[0]:
		return v
	case v := <-chs[1]:
		return v // 5
	}
}

func MapChan() int {
	m := map[string]chan int{"a": make(chan int, 1)}
	m["a"] <- 6
	return <-m["a"] // 6
}

type S struct{ C chan int }

func StructChanField() int {
	s := S{C: make(chan int, 1)}
	s.C <- 2
	return <-s.C // 2
}

func ChanPtrSend() int {
	c := make(chan int)
	go func(ch chan int) { ch <- 9 }(c)
	return <-c // 9
}

func ChanReturn() int {
	mk := func() chan int {
		c := make(chan int, 1)
		c <- 10
		return c
	}
	return <-mk() // 10
}

func ChanOfSlice() int {
	ch := make(chan []int, 1)
	ch <- []int{1, 2, 3}
	s := <-ch
	return s[0] + s[1] + s[2] // 6
}

func ChanOfMap() int {
	ch := make(chan map[string]int, 1)
	ch <- map[string]int{"k": 7}
	m := <-ch
	return m["k"] // 7
}

func main() {}

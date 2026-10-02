package main

type Box struct{ v int }

// select recv into existing field/index/map targets.
func RecvToField() int {
	ch := make(chan int, 1)
	ch <- 7
	b := Box{}
	select {
	case b.v = <-ch:
	}
	return b.v // want 7
}

func RecvToIndex() int {
	ch := make(chan int, 1)
	ch <- 9
	a := []int{0}
	select {
	case a[0] = <-ch:
	}
	return a[0] // want 9
}

func RecvToMap() int {
	ch := make(chan int, 1)
	ch <- 11
	m := map[string]int{}
	select {
	case m["k"] = <-ch:
	}
	return m["k"] // want 11
}

func RecvToBlank() int {
	ch := make(chan int, 1)
	ch <- 13
	select {
	case _ = <-ch:
	}
	return 1 // want 1
}

type MyChan chan int

func LenCapNamed() int {
	ch := make(MyChan, 3)
	ch <- 1
	return len(ch)*10 + cap(ch) // want 13
}

func CloseNamed() int {
	ch := make(MyChan, 1)
	close(ch)
	_, ok := <-ch
	if ok {
		return -1
	}
	return 1 // want 1
}

func main() {}

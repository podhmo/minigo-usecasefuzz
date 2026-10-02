package main

func NestedSelect() int {
	a := make(chan int, 1)
	b := make(chan int, 1)
	a <- 1
	b <- 2
	r := 0
	select {
	case v := <-a:
		select {
		case w := <-b:
			r = v + w
		default:
			r = v
		}
	default:
	}
	return r // want 3
}

func SelectInLoop() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	sum := 0
	done := false
	for !done {
		select {
		case v := <-ch:
			sum += v
		default:
			done = true
		}
	}
	return sum // want 6
}

func DupArms() int {
	ch := make(chan int, 1)
	ch <- 5
	got := 0
	select {
	case v := <-ch:
		got = v
	case v := <-ch:
		got = v + 100
	}
	return got // want 5 or 105 (random ready pick)
}

func ClosedArmWithDefault() int {
	ch := make(chan int, 1)
	close(ch)
	select {
	case v, ok := <-ch:
		if ok {
			return -1
		}
		return v + 10 // ready recv beats default: 10
	default:
		return 9
	}
}

func SelectInGoroutine() int {
	done := make(chan int, 1)
	go func() {
		ch := make(chan int, 1)
		ch <- 3
		select {
		case v := <-ch:
			done <- v
		}
	}()
	return <-done // 3
}

func chOf() chan int {
	c := make(chan int, 1)
	c <- 11
	return c
}

func SelectExprOperand() int {
	select {
	case v := <-chOf():
		return v // 11
	}
}

func SelectInIf() int {
	ch := make(chan int, 1)
	ch <- 4
	x := true
	if x {
		select {
		case v := <-ch:
			return v
		}
	}
	return 0
}

// send arm on a closed channel inside select: Go panics even with a
// default present.
func SelectSendOnClosed() (r int) {
	defer func() {
		if x := recover(); x != nil {
			r = -1
		}
	}()
	ch := make(chan int)
	close(ch)
	select {
	case ch <- 1:
		return 1
	default:
		return 2
	}
}

// select whose recv arm is a send-only-typed expr? (direction erasure)
func SelectTwoValueClosed() int {
	ch := make(chan int, 1)
	close(ch)
	select {
	case v, ok := <-ch:
		if ok {
			return -1
		}
		return v // 0
	default:
		return 9
	}
}

// select result used through switch inside case
func SelectMixedArms() int {
	recv := make(chan int, 1)
	send := make(chan int, 1)
	recv <- 8
	go func() { <-send }()
	out := 0
	select {
	case v := <-recv:
		out = v
	case send <- 9:
		out = 90
	}
	return out // 8 or 90 — random pick; assert membership
}

func main() {}

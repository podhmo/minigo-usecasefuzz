package main

func RangeBreak() int {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	sum := 0
	for v := range ch {
		if v == 3 {
			break
		}
		sum += v
	}
	return sum // 3
}

func RangeTwoVars() int {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	sum := 0
	for k, v := range ch {
		sum += k + v
	}
	return sum
}

func RangeInGoroutine() int {
	ch := make(chan int)
	done := make(chan int, 1)
	go func() {
		sum := 0
		for v := range ch {
			sum += v
		}
		done <- sum
	}()
	ch <- 1
	ch <- 2
	close(ch)
	return <-done // 3
}

func RangeContinue() int {
	ch := make(chan int, 4)
	for i := 0; i < 4; i++ {
		ch <- i
	}
	close(ch)
	sum := 0
	for v := range ch {
		if v == 2 {
			continue
		}
		sum += v
	}
	return sum // 4
}

// range with a labeled break out of nested loops
func RangeNestedBreak() int {
	a := make(chan int, 3)
	a <- 1
	a <- 2
	a <- 3
	close(a)
	out := 0
outer:
	for v := range a {
		for i := 0; i < 3; i++ {
			if v == 2 {
				break outer
			}
		}
		out += v
	}
	return out // 1
}

func main() {}

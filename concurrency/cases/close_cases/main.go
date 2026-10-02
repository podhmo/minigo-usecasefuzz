package main

func CloseTwice() (r int) {
	ch := make(chan int)
	close(ch)
	defer func() {
		if x := recover(); x != nil {
			r = 1
		}
	}()
	close(ch)
	return 0
}

func SendOnClosed() (r int) {
	defer func() {
		if x := recover(); x != nil {
			r = 1
		}
	}()
	ch := make(chan int, 1)
	close(ch)
	ch <- 1
	return 0
}

func CloseNil() (r int) {
	defer func() {
		if x := recover(); x != nil {
			r = 1
		}
	}()
	var ch chan int
	close(ch)
	return 0
}

func CloseWakesReceivers() int {
	ch := make(chan int)
	got := make(chan int, 3)
	for i := 0; i < 3; i++ {
		go func() {
			v, ok := <-ch
			if ok {
				got <- -1
				return
			}
			got <- v
		}()
	}
	close(ch)
	return <-got + <-got + <-got // 0
}

// sender parked on unbuffered chan is panicked by close -> process fails
func SendBlockedThenClosed() int {
	ch := make(chan int)
	go func() { ch <- 1 }()
	close(ch)
	return 0
}

// close named chan type
type NC chan int

func CloseNamed() int {
	ch := make(NC, 1)
	ch <- 1
	close(ch)
	v, ok := <-ch
	if !ok {
		return -1
	}
	return v // 1
}

// closed chan delivers buffered values first
func CloseDrainBuffered() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	a, b, c := <-ch, <-ch, <-ch
	d, ok := <-ch
	if ok {
		return -1
	}
	return a + b + c + d // 6
}

// recv on closed chan in a spawned goroutine
func ClosedRecvInGoroutine() int {
	ch := make(chan int)
	done := make(chan int, 1)
	go func() {
		v, ok := <-ch
		if ok {
			done <- -1
			return
		}
		done <- v
	}()
	close(ch)
	return <-done // 0
}

func main() {}

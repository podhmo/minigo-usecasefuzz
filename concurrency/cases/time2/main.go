package main

import "time"

func AfterDirect() int {
	<-time.After(time.Millisecond)
	return 1
}

func TimerChan() int {
	t := time.NewTimer(time.Millisecond)
	<-t.C
	return 1
}

func TimerStopReset() int {
	t := time.NewTimer(time.Hour)
	t.Stop()
	select {
	case <-t.C:
		return -1
	default:
	}
	t.Reset(time.Millisecond)
	<-t.C
	return 1
}

func TickerRangeBreak() int {
	t := time.NewTicker(time.Millisecond)
	n := 0
	for v := range t.C {
		_ = v
		n++
		if n == 3 {
			break
		}
	}
	t.Stop()
	return n // 3
}

func SleepGoroutine() int {
	done := make(chan int, 1)
	go func() {
		time.Sleep(time.Millisecond)
		done <- 4
	}()
	return <-done // 4
}

func AfterInLoopSelect() int {
	// each iteration creates a fresh timer: fires when nothing else ready
	n := 0
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	for {
		select {
		case v := <-ch:
			n += v
		case <-time.After(time.Millisecond):
			return n // 6
		}
	}
}

func main() {}

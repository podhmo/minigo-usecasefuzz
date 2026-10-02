package main

import "sync"

func RootSelfDeadlock() int {
	ch := make(chan int)
	ch <- 1
	return 0
}

func MutexDoubleLock() int {
	var mu sync.Mutex
	mu.Lock()
	mu.Lock()
	return 0
}

func EmptySelectRoot() int {
	select {}
	return 0
}

func CrossDeadlock() int {
	a := make(chan int)
	b := make(chan int)
	go func() {
		a <- 1
		b <- 2
	}()
	<-b
	<-a
	return 0
}

func main() {}

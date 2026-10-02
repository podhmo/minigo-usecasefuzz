package main

import "sync"

func MutexUnlockUnlocked() (r int) {
	defer func() {
		if x := recover(); x != nil {
			r = 1
		}
	}()
	var mu sync.Mutex
	mu.Unlock()
	return 0
}

func WaitGroupNegative() (r int) {
	defer func() {
		if x := recover(); x != nil {
			r = 1
		}
	}()
	var wg sync.WaitGroup
	wg.Done()
	return 0
}

func WaitGroupWaitInGoroutine() int {
	var wg sync.WaitGroup
	done := make(chan int, 1)
	wg.Add(1)
	go func() {
		wg.Wait()
		done <- 1
	}()
	wg.Done()
	return <-done // 1
}

type Box struct {
	mu sync.Mutex
	n  int
}

func MutexInStruct() int {
	b := &Box{}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.mu.Lock()
			b.n++
			b.mu.Unlock()
		}()
	}
	wg.Wait()
	return b.n // 4
}

func RWMutexUse() int {
	var mu sync.RWMutex
	val := 0
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		mu.Lock()
		val = 5
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		mu.RLock()
		mu.RUnlock()
	}()
	wg.Wait()
	return val // 5
}

func MutexPtrMember() int {
	var mu sync.Mutex
	p := &mu
	p.Lock()
	p.Unlock()
	return 1
}

// wg2 := wg aliases the same host WaitGroup (a Go copy would not share)
func WaitGroupCopyAlias() int {
	var wg sync.WaitGroup
	wg.Add(1)
	wg2 := wg
	go func() { wg2.Done() }()
	wg.Wait()
	return 1
}

func OnceAcrossGoroutines() int {
	var o sync.Once
	n := 0
	f := func() { n++ }
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			o.Do(f)
		}()
	}
	wg.Wait()
	return n // 1
}

func main() {}

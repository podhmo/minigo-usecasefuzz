package main

import (
	"parkprobe"
	"sync"
)

// DetachedWaitProbe: a goroutine parked inside a host call
// (WaitGroup.Wait — no select, no done arm) is not released when the
// process dies — this leaks the host goroutine (unlike channel parking,
// which procExit frees). parkprobe.Parked is bound by the host driver:
// it reports the goroutine is about to park, so the check needs no
// goroutine counting (a net NumGoroutine delta can be masked by any
// unrelated goroutine death).
func DetachedWaitProbe() int {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var inner sync.WaitGroup
		inner.Add(1)
		parkprobe.Parked()
		inner.Wait() // parks inside a host call — survives proc kill
	}()
	return 1
}

// DetachedLeakProbe: a goroutine parked in a channel op watches
// proc.done, so a dead process unwinds it with procExit — its deferred
// parkprobe.Gone fires on the way out, proving the goroutine died with
// its process rather than leaking.
func DetachedLeakProbe() int {
	ch := make(chan int)
	go func() {
		defer parkprobe.Gone()
		<-ch // parks in a channel op — unwinds when the proc dies
	}()
	return 1
}

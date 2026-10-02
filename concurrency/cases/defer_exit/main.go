package main

// root parked in select{} is released by a sibling panic; its defers
// should still run (they show up via println side effects).
func DeferOnProcExit() int {
	defer println("defer-one")
	defer println("defer-two")
	go func() { panic("boom") }()
	select {}
}

// a deferred call that parks during proc exit aborts its own op; the
// remaining defers ideally still run.
func DeferChainOnExit() int {
	defer println("A") // registered first, runs last
	defer println("B")
	defer func() { <-make(chan int) }() // parks -> procExit
	defer println("C")
	go func() { panic("boom") }()
	select {}
}

// recover() inside a deferred call during proc exit must not see it.
func RecoverOnProcExit() int {
	defer func() {
		if x := recover(); x != nil {
			println("saw panic:", x)
		} else {
			println("recover nil")
		}
	}()
	go func() { panic("boom") }()
	select {}
}

func main() {}

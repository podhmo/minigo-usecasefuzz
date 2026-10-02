package main

import "time"

func A_StopThenReset() int {
	println("new")
	t := time.NewTimer(time.Hour)
	println("stop", t.Stop())
	println("reset")
	ok := t.Reset(time.Millisecond)
	println("reset->", ok)
	<-t.C
	println("fired")
	return 1
}

func B_OnlyReset() int {
	t := time.NewTimer(time.Hour)
	t.Reset(time.Millisecond)
	<-t.C
	return 1
}

func C_OnlyStop() int {
	t := time.NewTimer(time.Hour)
	s := t.Stop()
	if s {
		return 1
	}
	return 2
}

func D_DurConst() int {
	d := time.Millisecond
	if d == 0 {
		return -1
	}
	return 1
}

func main() {}

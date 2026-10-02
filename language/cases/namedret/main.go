package main

import "fmt"

func f() (r int) {
	defer func() { r = 99 }() // defer overwrites named return
	r = 1
	return
}

func g() (r int) {
	defer func() { recover(); r = -1 }()
	panic("boom")
}

func h() (r int) {
	defer func() { r += 10 }()
	return 5
}

func p() (r string) {
	defer func() {
		if v := recover(); v != nil {
			r = fmt.Sprintf("caught %v", v)
		}
	}()
	panic("xyzzy")
	return "no"
}

func main() {
	fmt.Println(f())
	fmt.Println(g())
	fmt.Println(h())
	fmt.Println(p())
}

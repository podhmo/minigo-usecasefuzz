package main

import "fmt"

func try(name string, f func() any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s: panic %v\n", name, r)
		}
	}()
	fmt.Printf("%s: %v\n", name, f())
}

func main() {
	try("ifinit", func() any {
		if x := 5; x > 3 {
			return x
		}
		return 0
	})
	try("ifelse", func() any {
		if false {
			return 1
		} else if false {
			return 2
		} else {
			return 3
		}
	})
	try("forCond", func() any {
		i := 0
		for i < 5 {
			i++
		}
		return i
	})
	try("forEver", func() any {
		i := 0
		for {
			i++
			if i > 3 {
				break
			}
		}
		return i
	})
	try("forPost", func() any {
		n := 0
		for i, j := 0, 10; i < j; i, j = i+1, j-1 {
			n++
		}
		return n
	})
	try("switchTag", func() any {
		switch 2 {
		case 1:
			return "one"
		case 2, 3:
			return "two"
		default:
			return "other"
		}
	})
	try("switchNoTag", func() any {
		x := 5
		switch {
		case x < 0:
			return "neg"
		case x < 10:
			return "small"
		default:
			return "big"
		}
	})
	try("switchNoMatch", func() any {
		switch 99 {
		case 1:
			return "one"
		}
		return "none"
	})
	try("switchEmpty", func() any {
		switch {
		}
		return "done"
	})
	try("fallthrough", func() any {
		v := 0
		switch 1 {
		case 1:
			v = 1
			fallthrough
		case 2:
			v += 10
		case 3:
			v += 100
		}
		return v
	})
	try("fallDefault", func() any {
		v := 0
		switch 1 {
		case 1:
			v = 1
			fallthrough
		default:
			v += 100
		}
		return v
	})
	try("switchInit", func() any {
		switch x := 4; {
		case x%2 == 0:
			return "even"
		default:
			return "odd"
		}
	})
	try("labeledBreak", func() any {
		n := 0
	Outer:
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				if j == 2 {
					break Outer
				}
				n++
			}
		}
		return n // 2
	})
	try("labeledContinue", func() any {
		n := 0
	Outer:
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if j == 1 {
					continue Outer
				}
				n++
			}
		}
		return n // 3
	})
	try("gotoFwd", func() any {
		v := 0
		goto L
		v = 99
	L:
		return v
	})
	try("gotoBack", func() any {
		i := 0
	Top:
		i++
		if i < 3 {
			goto Top
		}
		return i
	})
	try("gotoOverDecl", func() any {
		goto E // Go spec: jumping over decl of x is a compile error — check
		return 1
	E:
		return 2
	})
	try("continueSkipsPost", func() any {
		n := 0
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				continue
			}
			n += i
		}
		return n
	})
	try("rangeInt", func() any {
		n := 0
		for i := range 5 {
			n += i
		}
		return n // Go 1.22: 0+1+2+3+4=10
	})
	try("nestedSwitch", func() any {
		switch 1 {
		case 1:
			switch 2 {
			case 2:
				return "inner"
			}
		}
		return "?"
	})
	try("switchScope", func() any {
		x := 1
		switch x := 9; x { // shadow
		case 9:
			return x
		}
		return x
	})
}

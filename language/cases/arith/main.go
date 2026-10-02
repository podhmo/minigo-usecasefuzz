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
	try("intdiv", func() any { return 7 / 2 })
	try("intmod", func() any { return -7 % 3 })
	try("intmod2", func() any { return 7 % -3 })
	try("negdiv", func() any { return -7 / 2 })
	try("overflow8", func() any { var x int8 = 127; x++; return x })
	try("underflow8", func() any { var x int8 = -128; x--; return x })
	try("uintwrap", func() any { var x uint8 = 0; x--; return x })
	try("shiftleft", func() any { return 1 << 10 })
	try("shiftright", func() any { return -16 >> 2 })
	try("ushiftright", func() any { var x int = -16; return x >> 2 })
	try("shiftbyvar", func() any { n := uint(3); return 1 << n })
	try("shiftneg", func() any { n := -1; return 1 << n })
	try("bigshift", func() any { n := uint(70); return 1 << n }) // wraps to 0 in Go
	try("floatdiv", func() any { return 7.0 / 2 })
	try("floatzero", func() any {
		x := 1.0
		return x / 0.0 // +Inf
	})
	try("floatnegzero", func() any {
		x := -1.0
		return x / 0.0
	})
	try("zerodivzero", func() any {
		x := 0.0
		return x / x // NaN
	})
	try("naneq", func() any {
		x := 0.0
		x = x / x
		return x == x // false
	})
	try("intdivzero", func() any { x, y := 7, 0; return x / y })  // runtime panic
	try("intmodzero", func() any { x, y := 7, 0; return x % y })
	try("minintdiv", func() any { x := -9223372036854775808; y := -1; return x / y })
	try("minintmod", func() any { x := -9223372036854775808; y := -1; return x % y })
	try("bitor", func() any { return 0xF0 | 0x0F })
	try("bitand", func() any { return 0xF3 & 0x0F })
	try("bitxor", func() any { return 0xFF ^ 0x0F })
	try("bitandnot", func() any { return 0xFF &^ 0x0F })
	try("bitnot", func() any { return ^5 })
	try("boolshort", func() any { x, y := 1, 0; return true || x/y != 0 })
	try("boolshort2", func() any { x, y := 1, 0; return false && x/y != 0 })
	try("intfloatmix", func() any { return 1 + 0.5 })
	try("runeadd", func() any { return 'a' + 1 })
	try("prec", func() any { return 2 + 3*4 - 10/3 })
	try("unaryF", func() any { return -1.5 })
	try("posU8", func() any { var x uint8 = 5; return -x }) // Go: -5 as int? no — unary - on uint8 gives uint8(251)!
}

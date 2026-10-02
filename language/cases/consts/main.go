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

const (
	A = iota
	B
	C
)
const (
	D = 1 << iota
	E
	F
)
const G = 1.5
const H = "hi"
const big = 1 << 62

const expr = A + B*C
const typedInt int = 7
const typedFloat float64 = 2.5

func main() {
	try("iota", func() any { return fmt.Sprintf("%d%d%d", A, B, C) })
	try("iotaShift", func() any { return fmt.Sprintf("%d%d%d", D, E, F) })
	try("expr", func() any { return expr })
	try("bigConst", func() any { return big })
	/* moved to lim-bigconst (designed: no arbitrary-precision consts) */
	/* moved to lim-bigconst */
	try("untypedOps", func() any { return A + 0.5 })       // untyped const + float → 0.5
	try("typedMix", func() any { return typedInt + 1 })
	try("constStrOps", func() any { return H + "!" })
	try("constDiv", func() any { return 7 / 2 })           // untyped const div → 3 (int)
	try("constFloatDiv", func() any { return 7.0 / 2 })    // 3.5
	try("typedFloatOps", func() any { return typedFloat * 2 })
	try("constShadow", func() any {
		const A = 99 // local shadow of package const
		return A
	})
	try("constInIf", func() any {
		if A == 0 {
			return "zero"
		}
		return "nonzero"
	})
	try("negativeConst", func() any { return -big })       // unary minus on const
	try("constBool", func() any { return true && false })
	try("constCmp", func() any { return A < C })
	try("runeConst", func() any { return 'x' })
	try("constAsFloat64", func() any {
		var f float64 = A // untyped→typed assign
		return f
	})
}

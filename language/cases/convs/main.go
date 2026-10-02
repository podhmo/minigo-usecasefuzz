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

type MyInt int
type MyStr string
type MySlice []int

func main() {
	try("intToFloat", func() any { return float64(7) / 2 })
	try("floatToInt", func() any { f := 3.9; return int(f) })   // truncate → 3
	try("floatToIntNeg", func() any { f := -3.9; return int(f) }) // → -3
	try("floatToIntU8", func() any { f := 300.0; return uint8(f) })
	try("floatNegU8", func() any { f := -1.5; return uint8(f) }) // impl-defined: Go→255
	try("intToU8", func() any { x := 300; return uint8(x) })    // 44
	try("negToU8", func() any { x := -1; return uint8(x) })     // 255
	try("u8ToInt8", func() any { x := 200; return int8(x) })    // -56
	try("i64ToI32", func() any { x := int64(1) << 40; return int32(x) }) // truncate → 0
	try("intToRune", func() any { return rune(97) })
	try("runeToInt", func() any { return int('a') })
	try("intToStr", func() any { x := 65; return string(x) })   // "A"
	try("negIntToStr", func() any { x := -1; return string(x) }) // ""
	try("bigIntToStr", func() any { x := 0x10FFFF + 1; return string(x) }) // invalid rune → RuneError
	try("strToBytes", func() any { return []byte("ab") })
	try("bytesToStr", func() any { return string([]byte{97, 98}) })
	try("strToRunes", func() any { return []rune("aあ") }) // [97 12354]
	try("runesToStr", func() any { return string([]rune{97, 12354}) })
	try("named", func() any { return int(MyInt(5)) })
	try("namedRev", func() any { return MyInt(5) + MyInt(3) })
	try("namedArith", func() any { return MyInt(5) + 3 })
	try("namedSlice", func() any { return len(MySlice{1, 2}) })
	try("namedSliceConv", func() any { return MySlice([]int{1, 2}) })
	try("sliceBack", func() any { return []int(MySlice{1, 2}) })
	try("float32", func() any { return float32(1.5) })
	try("f64ToF32", func() any { return float32(3.14) })
	try("f32Precision", func() any { return float32(0.1) }) // float32 rounding visible via %v? prints 0.1
	try("strToMyStr", func() any { return MyStr("x") + "y" })
	try("boolConv", func() any { return bool(true) })
	try("byteArr", func() any { return len([]byte("hé")) }) // utf8 bytes, len 3
	try("runeSlice", func() any { return fmt.Sprintf("%v", []rune("hé")) })
	try("strBytesAlias", func() any {
		b := []byte("abc")
		b[0] = 'X'
		return string(b) // "Xbc"
	})
	try("sliceToArrPtr", func() any { s := []int{1, 2, 3}; return (*[3]int)(s)[0] }) // Go 1.20
	try("floatToStrIface", func() any { return fmt.Sprintf("%v", any(1.5)) })
}

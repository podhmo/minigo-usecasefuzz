package main

import (
	"encoding/hex"
	"fmt"
)

// FNV-1a implemented in pure script — a "checksum a string" job that
// needs no hash library.
func fnv1a(s string) uint64 {
	h := uint64(14695981039346656037)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

func main() {
	for _, s := range []string{"", "a", "hello world", "日本語"} {
		fmt.Printf("%-12q %016x\n", s, fnv1a(s))
	}
	// hex round trip
	b := []byte{0xde, 0xad, 0xbe, 0xef}
	fmt.Println(hex.EncodeToString(b))
	d, err := hex.DecodeString("ff00")
	fmt.Println(d, err)
	_, err = hex.DecodeString("zz")
	fmt.Println("bad hex err:", err != nil)
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// SHA-256 checksum of a string.
func main() {
	sum := sha256.Sum256([]byte("hello"))
	fmt.Println(hex.EncodeToString(sum[:]))
}

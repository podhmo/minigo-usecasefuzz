package main

import (
	"bytes"
	"fmt"
	"io"
	"probe"
)

var buf []byte
var dst []byte
var src string
var S int

func Setup(n int) {
	buf = make([]byte, n)
	dst = make([]byte, n)
	src = string(buf)
	S = 0
}

func Noop() {}

func MakeBuf(n int) {
	buf = make([]byte, n)
}

// script []byte -> host: bytes.NewReader marshals through goNative
// ([]Value -> []any) then byteSlice ([]any -> []byte) — two passes.
func PNewReader() {
	r := bytes.NewReader(buf)
	_ = r
}

// script []byte -> reflect method arg (toReflectValue element-wise)
// plus the callee's writes mirrored back element-wise.
func PRead() {
	r := probe.NewReader(len(buf))
	n, _ := r.Read(buf)
	S = n
}

// script []byte -> reflect method arg (toReflectValue element-wise),
// read-only direction (no write-back needed but still mirrored).
func PWrite() {
	n, _ := probe.DevNull.Write(buf)
	S = n
}

// io.ReadFull goes through the h.fn path: byteSlice(a[1]) makes a
// FRESH []byte — writes do not propagate back to the script slice.
func PReadFull() {
	n, _ := io.ReadFull(probe.FillReader, buf)
	S = n
}

// PReadCheck reports buf[0] after ReadFull: 171 (0xAB) if writes
// landed on the script slice, 0 if the write-back is missing.
func PReadCheck() int {
	io.ReadFull(probe.FillReader, buf)
	return int(buf[0])
}

// host []byte -> script: scriptVal boxes every byte to a tagged Value.
func PReadAll(n int) {
	b, _ := io.ReadAll(probe.NewReader(n))
	S = len(b)
}

// string(b) — element-wise unwrap + copy.
func PToString() {
	s := string(buf)
	S = len(s)
}

// []byte(s) — element-wise coerce to tagged Values.
func PFromString() {
	b := []byte(src)
	S = len(b)
}

// per-element index read — interpreter-internal baseline.
func PIndex() {
	sum := 0
	for i := 0; i < len(buf); i++ {
		sum += int(buf[i])
	}
	S = sum
}

// per-element range read — interpreter-internal baseline.
func PRange() {
	sum := 0
	for _, x := range buf {
		sum += int(x)
	}
	S = sum
}

// copy(dst, src) — element-wise v.Copy snapshot.
func PCopy() {
	S = copy(dst, buf)
}

// append growth — one element at a time.
func PAppend(n int) {
	var b []byte
	for i := 0; i < n; i++ {
		b = append(b, byte(i))
	}
	S = len(b)
}

// fmt %x — sliceBytes materializes a fresh []byte first.
func PFmtHex() {
	s := fmt.Sprintf("%x", buf)
	S = len(s)
}

// io.Copy — NewReader marshals buf once; the copy itself never
// crosses back into script values (contrast point).
func PIoCopy() {
	r := bytes.NewReader(buf)
	n, _ := io.Copy(io.Discard, r)
	S = int(n)
}

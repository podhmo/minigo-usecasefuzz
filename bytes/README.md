# bytes — []byte boundary cost probes

How expensive is it for `minigo` to move a `[]byte` between script and
host code? This harness measures every crossing a script can trigger —
through the real stdlib bindings and the public `Engine` API — at small
sizes, so the per-byte shape of each path is visible without any load.

Context: podhmo/minigo#361 (comment). `runtime.Slice` stores `Elems
[]Value` — every byte is a tagged heap `*Named{byte, int64}`, and every
boundary marshal walks it element-wise.

## Run

```console
$ go run .                     # default sizes 4KiB,64KiB,1MiB, 5 reps
$ go run . -sizes 4096 -reps 3
$ ./run.sh                     # same, plus `go vet`
```

Requires a `podhmo/minigo` checkout next to this repo (the go.mod
`replace` resolves `../../minigo`, same convention as
`concurrency/host`; `run.sh` clones it if missing). The first line of
output is a **functional check**:
`io.ReadFull` currently never writes into the script's slice
(`byteSlice()` copies), so `buf[0]` reports 0, not 0xAB.

## Measured (linux/amd64, Go 1.26, minigo @ 2026-10-04)

| op | path | ns/byte @1MiB | allocs/byte |
|---|---|---|---|
| `index b[i]` loop | interpreter | ~700 | ~7.0 |
| `append` x1 | interpreter | ~776 | ~9.0 |
| `range` loop | interpreter | ~323 | ~4.0 |
| `r.Read(buf)` | reflect marshal + write-back | ~153 | ~4.0 |
| `w.Write(buf)` | reflect marshal | ~137 | ~4.0 |
| `[]byte(str)` | convert | ~47 | ~1.0 |
| `io.ReadFull` | h.fn marshal (**no write-back**) | ~34 | ~0 |
| `io.Copy(Discard, r)` | NewReader marshal + host copy | ~31 | ~0 |
| `io.ReadAll` | host→script boxing | ~21 | ~1.0 |
| `copy(dst,src)` | interpreter | ~19 | ~1.0 |
| `make([]byte,n)` | interpreter | ~15 | ~1.0 |
| `bytes.NewReader` | h.fn double marshal | ~7 | ~0 |
| `fmt "%x"` | sliceBytes | ~6 | ~0 |
| `string(buf)` | convert | ~3 | ~0 |

(`h.fn` marshal allocates little because byte values ≤255 box via Go's
static small-int cache; it still walks the elements twice.)

## Where the cost is

1. **In-script per-element work is the worst case** — 300-800 ns/byte,
   4-9 allocs/byte. A 1 MiB `b[i]` scan is ~0.7 s. This is generic
   interpreter dispatch, not `[]byte`-specific, and no representation
   change removes it (a read must still produce a `Value`). Scripts that
   scan bytes should hand the buffer to the host instead.
2. **Reflect-call marshal (`Read`/`Write` methods)** — ~150 ns/byte +
   ~4 allocs/byte from the recursive `toReflectValue` per element plus
   the element-wise write-back. This is the boundary the issue comment
   targets, and the cheapest big win: a specialized byte path (tight
   `Unwrap` loop, like `string(buf)` already does at 3 ns/byte) should
   land near ~10-20 ns/byte.
3. **Host→script boxing (`io.ReadAll`, `[]byte(str)`)** — ~20-47
   ns/byte, exactly ~1 alloc/byte: the `*Named` tag per element. That is
   the floor imposed by `Elems []Value`; only a packed representation
   removes it (and the ~24-40 B/element memory it implies).
4. **`io.ReadFull` is a correctness bug, not just slow** — the `h.fn`
   path's `byteSlice()` copies the buffer and writes back nothing; the
   reflect path writes back correctly. The two paths disagree on
   borrow/write-back semantics.

## Layout

- `script/probe.go` — minigo-side probe functions (`Setup(n)` sizes
  `buf`/`dst`/`src`; each `P*` function is one boundary op)
- `main.go` — driver: builds the engine, binds `probe` (a GoValue
  writer, an 0xAB-filling reader, a bounded-reader factory), reports
  ns/op + allocs/op per probe
- `run.sh` — `go vet` + `go run .`

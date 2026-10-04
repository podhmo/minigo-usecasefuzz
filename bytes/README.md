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
output is a **functional check**: `buf[0]` after `io.ReadFull`
should be 0xAB. Before podhmo/minigo#362 landed it reported 0 —
`byteSlice()` copied the buffer and dropped every byte the reader
wrote.

## Measured (linux/amd64, Go 1.26)

Before = minigo @ 2026-10-04 (pre-#362). After = minigo with the
byte fast-path (`toReflectValue` one-pass marshal + `callReflectFunc`
write-back fast-path + `io.ReadFull` borrow/commit).

| op | path | before ns/byte @1MiB | after ns/byte | after allocs/byte |
|---|---|---|---|---|
| `index b[i]` loop | interpreter | ~700 | ~697 | ~7.0 |
| `append` x1 | interpreter | ~776 | ~797 | ~9.0 |
| `range` loop | interpreter | ~323 | ~300 | ~4.0 |
| `r.Read(buf)` | reflect marshal + write-back | ~153 | ~16 | ~1.0 |
| `w.Write(buf)` | reflect marshal | ~137 | ~16 | ~1.0 |
| `[]byte(str)` | convert | ~47 | ~51 | ~1.0 |
| `io.ReadFull` | **h.fn copy → borrow/commit** | ~34 | ~18 | ~1.0 |
| `io.Copy(Discard, r)` | NewReader marshal + host copy | ~31 | ~12 | ~0 |
| `io.ReadAll` | host→script boxing | ~21 | ~18 | ~1.0 |
| `copy(dst,src)` | interpreter | ~19 | ~16 | ~1.0 |
| `make([]byte,n)` | interpreter | ~15 | ~18 | ~1.0 |
| `bytes.NewReader` | h.fn double marshal | ~7 | ~8 | ~0 |
| `fmt "%x"` | sliceBytes | ~6 | ~6 | ~0 |
| `string(buf)` | convert | ~3 | ~4 | ~0 |

Reflect-path marshal went ~9x faster and from ~4 allocs/byte to ~1,
landing at the same floor as `string(buf)`'s tight loop — the
predicted outcome in §2 below. `io.ReadFull` also got faster as a
side effect of borrowing instead of double-copying. In-script
per-element work is unchanged, as expected — it never crosses the
boundary.

## Where the cost is

1. **In-script per-element work is the worst case** — 300-800 ns/byte,
   4-9 allocs/byte. A 1 MiB `b[i]` scan is ~0.7 s. This is generic
   interpreter dispatch, not `[]byte`-specific, and no representation
   change removes it (a read must still produce a `Value`). Scripts that
   scan bytes should hand the buffer to the host instead.
2. **Reflect-call marshal (`Read`/`Write` methods)** — was ~150
   ns/byte + ~4 allocs/byte from the recursive `toReflectValue` per
   element plus the element-wise write-back; the #362 byte fast-path
   brought it to ~16 ns/byte + ~1 alloc/byte, matching the
   `string(buf)` tight loop. The remaining floor is now the same
   `Elems []Value` boxing as everywhere else.
3. **Host→script boxing (`io.ReadAll`, `[]byte(str)`)** — ~20-47
   ns/byte, exactly ~1 alloc/byte: the `*Named` tag per element. That is
   the floor imposed by `Elems []Value`; only a packed representation
   removes it (and the ~24-40 B/element memory it implies).
4. **`io.ReadFull` was a correctness bug, not just slow** — the
   `h.fn` path's `byteSlice()` copied the buffer and wrote back
   nothing, so the script slice stayed zeroed while `r.Read(buf)`
   filled correctly. #362 rebinds it through `borrowBytes` (borrow
   the `*runtime.Slice`, commit its elements back after the call),
   which also made it ~2x faster.

## Layout

- `script/probe.go` — minigo-side probe functions (`Setup(n)` sizes
  `buf`/`dst`/`src`; each `P*` function is one boundary op)
- `main.go` — driver: builds the engine, binds `probe` (a GoValue
  writer, an 0xAB-filling reader, a bounded-reader factory), reports
  ns/op + allocs/op per probe
- `run.sh` — `go vet` + `go run .`

# language fuzz

Differential corpus for `minigo`'s **Go language surface**: spec edge cases
rather than realistic use cases (which live in the top-level `cases/`) or
concurrency semantics (`concurrency/`). Each `cases/<theme>/main.go` is a
`func main()` program run under both `go run` (the oracle) and
`minigo run`, with stdout+stderr diffed by `run.sh`.

Cases share one module (`go.mod` → `module langfuzz`), so intra-corpus
imports look like `langfuzz/cases/pkgs/sub`.

Findings from the round (20+ bugs, all fixed):
`docs/sketch/ja/fuzz-language.md` in the minigo repo (PR #29).

## Verdicts

| verdict | meaning |
|---|---|
| PASS | ran clean, stdout+stderr identical |
| DIFF | both ran, output differs |
| TRAP | Go ran clean; minigo died (panic/trap) |
| PASS-PANIC | both died at runtime (Go `panic:`, minigo error — stack formats aren't compared) |
| NOPANIC | Go panicked but minigo ran clean |
| PASS-REJECT | Go fails to compile; minigo also refuses (e.g. runtime trap) |
| ACCEPT | Go rejects but minigo ran — worst verdict |
| TIMEOUT | either side hit the 15s limit |

`neg-*` dirs hold programs Go rejects at compile time (the expected
verdict is PASS-REJECT — a runtime trap is the nearest loud failure).
`lim-*` dirs probe designed limitations that are *intended* to trap:
`lim-complex` (complex numbers), `lim-threeidx` (`s[a:b:c]`),
`lim-bigconst` (`1<<100` arbitrary-precision consts), `lim-uintptr`.

Latest run (2026-10-02, post-PR-#29): **35 PASS / 4 PASS-REJECT /
0 DIFF / 3 TRAP-by-design**.

## Usage

```sh
./run.sh                 # all cases
./run.sh mapkeys arrays  # selected cases
MINIGO_DIR=/path/to/minigo ./run.sh
```

Outputs land in `out/<case>.{go,mg}.{out,err}` (gitignored).

## Themes

| dir | surface exercised |
|---|---|
| appendalias | append aliasing / cap sharing |
| arith | sized ints, overflow, div/mod, shifts, floats |
| arrays | array↔slice slicing, cap, sharing, `[N]T(s)` conversions |
| assign | compound assign, `(*p) op=`, `x[i]++`, `&^=` |
| builtins | len/cap/copy/append/delete/make/print/println |
| consts | iota, typed/untyped consts, MinInt64 literal |
| control | if/switch/for/range/fallthrough, labeled break/continue |
| convs | numeric/string/slice conversions |
| defers | defer order, recover, panic re-raise |
| destruct | multi-return destructure, comma-ok, blank idents |
| embed | promoted fields/methods, ptr-receiver via addressable, iface satisfaction |
| errors2 | Error(), `%w` wrap, Unwrap, Is/As, Join |
| fmtverbs | `%b %o %x %e %g %q %c` width/precision `%+v %#v %T %[1]d %%` |
| forrangeint | Go 1.22 `range N` |
| funcs | values, closures, variadic, recursion, method exprs |
| generics | `F[T]`, `S[T]{}`, inference, local `type` decls |
| initfuncs | `init()` ordering |
| initorder | package-var init order (dependency DFS) |
| interfaces | method sets, dispatch, `any`, empty-iface |
| iterators | `iter.Seq`/`Seq2` ranges |
| loopvars | Go 1.22 per-iteration capture |
| mapkeys | array/struct/slice elided-literal keys, NaN, canonical identity |
| maps | literal shapes, comma-ok, delete, iteration |
| methods | value/ptr receivers, promoted, method values/exprs |
| namedret | named results + defer/recover mutation |
| neg-arridx | negative/const-out-of-range index — compile reject |
| neg-generics | bad type arg — compile reject |
| neg-mapaddr | `&m[k]` — compile reject |
| neg-strwrite | `s[i]=c` — compile reject |
| nilfieldwrite | nil-ptr field write panic + recover |
| pkgs | cross-package import/vars within the corpus module |
| ptrs | `&`, deref, `**`, pointer arithmetic traps |
| scope | shadowing, decl order, blanks |
| slices | slicing, full-slice exprs, append growth, cap rules |
| strings | index/len/rune iter, concat, conv |
| structs | literals, tags, anon structs, field access |
| typednilcall | method call on typed-nil through interface |
| tyswitch | type switch forms + fallthrough + nil case + comma-ok |

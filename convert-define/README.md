# convert-define/

Compile-oracle corpus for `examples/convert-define` — each `cases/<name>/`
defines a conversion via `define.go` (`//go:build codegen`, interpreted by
minigo), then `run.sh` regenerates `generated.go` and compiles the case as an
independent module (`module example.com/m` + `replace` → the minigo checkout).

```
./run.sh [case ...]    # MINIGO_DIR=../path/to/minigo to override the sibling clone
```

## Verdicts

| verdict | meaning |
|---|---|
| OK | generated converter compiles |
| GEN-FAIL | the tool itself errored (expected for DSL misuse) |
| BUILD-FAIL | generated code does not compile — the loud-failure path for leaf mismatches and unsupported misuse |
| TIDY-FAIL | `go mod tidy` failed — harness/module issue, not a generator verdict |

A `BUILD-FAIL` with `expect: BUILD-FAIL` is a pinned known-limitation, not a
bug: e.g. `int`→`string` leaf mismatch stays a compile error by design
(see `docs/sketch/ja/fuzz-convert-define.md` §3 in podhmo/minigo).

## Cases and what they exercise

| case | surface |
|---|---|
| c01-basic | same-name automatch + `c.Compute` |
| c02-map-convert | `c.Convert` with a `convutil` helper |
| c03-ptr | `*Sub`→`*Sub` dedicated converter call |
| c04-ptrval | `*Sub`→`Sub` nil-guard |
| c05-valptr | `Sub`→`*Sub` |
| c06-deepptr | `**Sub`→`**Sub` generic path |
| c07-slice | `[]Sub` element conversion |
| c08-array | `[N]int` direct assign + `[N]Sub` element loop |
| c09-mapval | `map[string]Sub` value conversion |
| c10-mapkey | `map[int]`→`map[int64]` key cast + struct key pair |
| c11-named | named leaf casts (`SrcStatus`→`DstStatus`, `int`→`int64`) |
| c12-samename | `a.User`→`b.User` and `a.User`→`c.User` name disambiguation |
| c13-embed | promoted field through embedding |
| c14-generic | `[]box.Box[int]` instantiated generic field |
| c15-jsontag | shared `json` tag match |
| c16-unexported | unexported fields skipped |
| c17-selfpkg | types in the generated package itself (self-import in define.go) |
| c18-nested | `c.Map` dotted paths (`dst.Inner.ID`, `src.In.V`) |
| leaf-mismatch | `int64`→`string` — pinned BUILD-FAIL (loud failure by design) |
| neg01-baresrc | `c.Map(dst.V, source.A{}.V)` non-ident src — GEN-FAIL |
| neg02-typo | `c.Map(dst.W, src.W)` unknown src member — GEN-FAIL |
| neg03-rulesig | `define.Rule(badfn)` inside `c.Convert` — BUILD-FAIL |

## Known gaps surfaced by this corpus

- `c.Convert(dst.X, src.Y, define.Rule(fn))` passes the `Rule(...)` call
  through verbatim and only fails at compile time — an early DSL error is a
  candidate improvement.
- Bare (unqualified) type names in `define.Convert` func params hit the
  LoadFile synthetic-path limitation; qualify them via a self-import
  (`import m "example.com/m"`) as c17 does.

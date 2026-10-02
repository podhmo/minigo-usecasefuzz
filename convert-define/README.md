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
| c19-bareself | types in the generated package itself, referenced by bare name (no self-import) |
| leaf-mismatch | `int64`→`string` — pinned BUILD-FAIL (loud failure by design) |
| neg01-baresrc | `c.Map(dst.V, source.A{}.V)` non-ident src — GEN-FAIL |
| neg02-typo | `c.Map(dst.W, src.W)` unknown src member — GEN-FAIL |
| neg03-rulesig | `define.Rule(badfn)` inside `c.Convert` — GEN-FAIL (non-function converter rejected at DSL eval) |
| broken-src | type error in a source-pkg func body — generation succeeds (AST only, no typecheck); BUILD-FAIL is the broken input package itself |
| broken-field | field of unresolved type `Ghost` — generation still succeeds (leaf cast `int64(src.V)` emitted); BUILD-FAIL is the broken input package itself |
| broken-dsl | DSL file that never compiles (unused import, `ghost()`, `var x int = "no"` in the func lit) — OK: minigo runs it, output is clean |
| broken-syntax | syntax error in source — GEN-FAIL: parse failure is the real boundary |
| stale-generated | stale `generated.go` in `package gen` referencing deleted src/dst fields (the whole package does not compile) — OK: `-file` targets the DSL file only, so regeneration overwrites the broken output and the package builds again. The stale state is committed as `generated.go.stale` and copied over the (gitignored) `generated.go` before each run |

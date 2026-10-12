# extract-schema/

Schema-output corpus for `examples/extract-schema` — each `cases/<name>/`
is a standalone module (`module example.com/m` + `replace` → the minigo
checkout) that `run.sh` runs the built tool inside of, so both the
tool's script packages and the fixture's own imports resolve there.

```
./run.sh [case ...]    # MINIGO_DIR=../path/to/minigo to override the sibling clone
```

## Per-case files

| file | meaning |
|---|---|
| `*.go`, `*/` | fixture package(s); target defaults to `.` |
| `target` | target arg (default `.`) |
| `type` | optional type name appended after the target |
| `flags` | extra flags as a shell fragment (`-stdlib walk`, `-plugin ./plug`) |
| `expect` | expected verdict |
| `want` | expected stdout — hand-verified ideal schema |
| `plugin*/` | optional plugin script dirs |

## Verdicts

| verdict | meaning |
|---|---|
| OK | exit 0, valid JSON, stdout matches `want` |
| WARN | exit 0 but stderr non-empty — the tool's own softened-failure channel |
| DIFF | exit 0, valid JSON, but stdout ≠ `want` — pinned semantic divergence |
| BADJSON | exit 0 but stdout is not valid JSON |
| FAIL | non-zero exit via a graceful error |
| TRAP | non-zero exit via `runtime trap` (engine-level fault) |
| PANIC | non-zero exit via `panic:` (script panic propagated) |

Captured output is committed under `out/` (`<name>.stdout`, `<name>.stderr`,
local paths sanitized to `MINIGO`/`CORPUS`/`HOME` tokens). Bug reports
from the findings live in podhmo/minigo `TODO.md`, not here.

## Cases

| case | surface | verdict |
|---|---|---|
| s01-scalars | every builtin field type | OK |
| s02-named-leaves | named builtin types, `[]byte` → base64 string | OK |
| s03-error-field | `error` builtin → loud failure | FAIL |
| s04-no-types | package with no type decls | OK |
| c01-slices-maps | slices, nested maps, struct-keyed map | OK |
| c02-arrays | fixed arrays | OK |
| c03-ptrs | `*T`, `**T`, `*[]int` — transparent + optional | OK |
| c04-func-chan | `func()` field → unsupported shape | FAIL |
| c05-chan-field | `chan` fields → unsupported shape | FAIL |
| c06-generics | `Box[int]` field → type-param leaf | OK |
| c07-generic-decl | generic decl in package mode → type-param leaf | OK |
| c08-interface-field | non-empty interface field | FAIL |
| c09-anon-structs | nested/slice/map anonymous structs | OK |
| c10-anon-embed | anonymous struct embedding a named type | OK |
| c11-generic-alias | alias to `Box[int]` → type-param leaf | OK |
| e01-embed-basic | value/ptr/tagged embeds, multi-level | OK |
| e02-embed-cycle | `A{*B}`,`B{*A}` → drops the cyclic subtree like encoding/json | OK |
| e03-embed-nonstruct | embedded `int`, `MyInt`, `time.Duration` | OK |
| e04-embed-unexported | unexported struct embed promotes fields | OK |
| e05-embed-ambig3 | ambiguity blocks, shallower explicit wins | OK |
| l01-stdlib-leaf | `time`, `net`, `url`, `big`, `json` leaves | OK |
| l02-stdlib-walk | `-stdlib walk` on `net` → `import "C"` files skipped | OK |
| l03-stdlib-walk-time | walk `time.Time` → unexported-only object | OK |
| l04-bound-host | bound `sync.Mutex`/`Once` → stdlib leaf | OK |
| l05-stdlib-walk-dur | walk `time.Duration` → enum scalars + inferred members | OK |
| m01-unknown-type | missing type name → graceful error | FAIL |
| m02-nontype-name | func name as type → graceful error | FAIL |
| m03-multifile | multi-file pkg, `_test.go` excluded | OK |
| m04-pkg-mode | `-pkg example.com/m` | OK |
| m05-pkg-stdlib | `-pkg time` → bound shadow retries source (`inspect.LookupSourceOf`), `Time` → unexported-only object | OK |
| m06-title-desc | `-title`/`-desc` flags | OK |
| m07-loose-false | `-loose=false` | OK |
| m08-bad-target | nonexistent dir → graceful error | FAIL |
| m09-package-mode | whole-package `$defs`, A↔B refs | OK |
| n01-alias-chains | `A = B`, alias to composites/struct | OK |
| n02-named-composites | named map/slice/array/ptr types | OK |
| n03-named-func | named `func` type → unsupported shape | FAIL |
| n04-enum-shapes | string/iota/split/single enums → materialized scalars | OK |
| n05-enum-struct | struct-typed const member → excluded like gc | OK |
| n06-weird-names | builtin-shadow + unicode names → want itself diverges from encoding/json (oracle-side mis-pin) | DIFF |
| p01-plugin-override | override local + stdlib types | OK |
| p02-plugin-desc | DescField/DescDecl filters | OK |
| p03-plugin-multi | two plugins, override last-wins, filters chain | OK |
| p04-plugin-empty | plugin with no configurer → warn | WARN |
| p05-plugin-wrongsig | `func(a,b *C)` fails the arity check → warn | WARN |
| p06-plugin-baddir | missing plugin dir → graceful error | FAIL |
| p07-plugin-trap | configurer panics (nil map write) | PANIC |
| p08-plugin-unexported | unexported entry func → warn | WARN |
| t01-tag-matrix | json/yaml/toml precedence, `-`, `,string` honored | OK |
| t02-required-tag | `required:"bool"` overrides ptr/omitempty | OK |
| t03-override-merge | `jsonschema-override` merges (`required=` applied) | OK |
| t04-override-malformed | malformed override → loud error | FAIL |
| t05-override-desc | override `description` beats doc comment | OK |
| x01-cross-basic | field typed by another package | OK |
| x02-cross-second-order | transitive A→B→C follow | OK |
| x03-cross-diamond | shared dep → shared `$ref` | OK |
| x04-cross-samename | `a.Money`/`b.Money` → `Money_b` rename | OK |
| x05-cross-embed | cross-package embed flattens | OK |

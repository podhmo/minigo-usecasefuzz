# usecasefuzz

Differential corpus for `minigo`: realistic Go use-case programs run under both
`go run` (the oracle) and `minigo run`, with stdout+stderr diffed.

Unlike the earlier fuzz rounds (language spec edge cases, concurrency
semantics), this corpus asks: **do the programs people actually write work?**
Bias is toward text processing / scripting ("LL-ish") usage rather than binary
manipulation or systems code.

## Layout

- `cases/<name>/` — one standalone module per use case (`func main()`)
- `language/` — a separate harness: Go-language-spec edge cases diffed
  against `go run` (see `language/README.md`)
- `concurrency/` — a separate harness: function-level probes + host-side
  checks for `go`/`chan`/`select`/`sync`/`time` semantics (see
  `concurrency/README.md`)
- `run.sh [case ...]` — builds `./out/minigo` from `$MINIGO_DIR` (default
  `~/repos/minigo`) and prints a verdict per case
- `out/` — captured `.want` (go) and `.got` (minigo) outputs

See `cases/README.md` for the verdict table and the case-by-case
index, and `concurrency/README.md` for the concurrency harness.

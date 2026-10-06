# realworld

Sync-tool scripts run against **pinned real-world codebases**. The other
harnesses in this repo ask "does Go code people write run under minigo?";
this one asks minigo's actual pitch: **can a tool that reads a legacy Go
codebase and cross-checks it against another artifact (OpenAPI spec,
plugin manifest, frontend types) run cheaper than a go/packages-based
tool — and keep running when the module graph is incomplete?**

Background and the first round's findings: podhmo/minigo
`docs/sketch/ja/exploration-legacy-sync.md`. The judgement loop (minimize →
pin → TODO) is the `realworld` skill in podhmo/minigo.

## Layout

- `targets.tsv` — `name<TAB>url<TAB>commit`; targets are shallow-fetched at
  that commit into `$SRC_DIR` (default `./.src`) and `go mod download`ed once.
- `tasks/<name>/` — one standalone module per task (`func main()`), reading
  the target checkout from `$TARGET_DIR`; `target` names the targets.tsv row.
  A task with `want.txt` is minigo-only (it imports `minigo.dev/inspect`) and
  is compared against that golden file instead of `go run`.
- `run.sh [task ...]` — verdict + wall times per task. `COLD=1` also times
  the oracle under an empty throwaway `GOCACHE` (the user's cache is never
  touched).

- `prof/` — `runtime/pprof` around the minigo engine (the CLI has no
  profiling flags). `PROFILE=1 ./run.sh [task ...]` builds it against
  `$MINIGO_DIR` through a generated go.mod and writes, per task,
  `out/<task>.cpu.pprof`, `out/<task>.allocs.pprof` and a text summary
  `out/<task>.prof.txt` (flat CPU, cumulative minigo frames, alloc_space).
  A task that traps still gets a profile of the run up to the trap.
- `reports/` — measurement rounds as self-contained HTML (open in a browser).

Verdicts: `PASS`, `DIFF` (silent divergence — a bug), `TRAP`, `REJECT`,
`HANG` (timeout), `ORACLE-FAIL`, `SETUP-FAIL`.

## Tasks

| task | target | reads | checks against |
|---|---|---|---|
| grafana-openapi | grafana | route registrations in `pkg/api/api.go` (func bodies) + `// swagger:route` comments, via `go/parser` | each other |
| grafana-coreplugin | grafana | `coreplugin` plugin-ID consts, via `go/parser` (no package init) | `public/app/plugins/datasource/*/plugin.json` |
| clickhouse-settings | clickhouse-datasource | `Settings` struct json tags, via `inspect` (surface only) | `CHConfig` in `src/types/config.ts` |

## Reports

| report | minigo | summary |
|---|---|---|
| [2026-10-07-measure.html](reports/2026-10-07-measure.html) | podhmo/minigo#534 | time to answer vs `go run` / `go vet` (warm and cold), the missing-modules run, CPU/alloc breakdown of grafana-openapi, GOGC sensitivity, peak RSS, ranked improvement room and usability gaps (Japanese) |

A round re-measures with `./run.sh <task ...>` (3 runs). Native baselines use
`go run .` for stdlib-only tasks and `go vet <pkg>` on the target as a stand-in
for a typed loader, each warm and with a throwaway `GOCACHE`. The incomplete
environment is an empty `GOMODCACHE` with `GOPROXY=off`. Profiles are taken
without `TRACE=1`, because tracing skews the CPU split.

## Adding a task

Pick a sync a real project needs (Go code ↔ another artifact). Prefer
inputs `inspect` can reach (surface) — when a task has to fall back to
`go/parser` because `inspect` can't see something, note *what* it couldn't
see; that gap is a finding. Keep `go run`-able when possible so the oracle
is free.

## Profiling notes

- Open a profile with `go tool pprof -http=: out/prof out/<task>.cpu.pprof`
  (`out/prof` is the binary that produced it).
- **macOS:** heavy goroutine spawning (e.g. a helper goroutine per
  `sync.Mutex.Lock`) can make the CPU profile pin most samples on libc
  (`pthread_cond_signal`/`pthread_cond_wait`) and hide the interpreter.
  `TRACE=1` adds `out/<task>.trace`; `go tool trace` shows whether the main
  goroutine was actually running. If it was, trust the alloc profile and
  re-profile on Linux or with the spawning path disabled.

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
  is compared against that golden file instead of `go run`. A task with an
  executable `task.sh` drives a whole program instead: `task.sh native` is the
  oracle and `task.sh minigo <BIN>` the minigo side; an optional
  `timeout_sec` file overrides `TIMEOUT_SEC`.
- `run.sh [task ...]` — verdict + wall times per task. `COLD=1` also times
  the oracle under an empty throwaway `GOCACHE` (the user's cache is never
  touched).
- `compare.sh OLD_REF NEW_REF [--full] [--profile]` — same-sitting A/B
  performance comparison of two minigo revisions over `probes/` plus light
  tasks (interleaved runs, median deltas, exit 2 on regression). Probe runs
  need no module downloads; `--full` adds grafana-openapi. The
  grafana-swagger-* tasks are pitch/drift checks, not regression probes.

- `PROFILE=1 ./run.sh [task ...]` profiles each task once more and writes,
  per task, `out/<task>.cpu.pprof`, `out/<task>.allocs.pprof`,
  `out/<task>.script.pprof` and a text summary `out/<task>.prof.txt`. The
  summary has the script functions and lines, flat CPU, the cumulative
  minigo frames, and alloc_space. It uses the minigo CLI's
  `$MINIGO_CPUPROFILE` / `$MINIGO_MEMPROFILE` / `$MINIGO_SCRIPTPROFILE`
  (podhmo/minigo#1035, #1036), so `task.sh` tasks are profiled too. A
  `task.sh` that runs minigo several times keeps the last run's profiles. A
  task that traps still gets a profile of the run up to the trap.
- `prof/` — `runtime/pprof` around the minigo engine, built against
  `$MINIGO_DIR` through a generated go.mod. `PROFILE=1` falls back to it
  under `TRACE=1` or for a minigo without the profiling flags (no script
  profile then, and `task.sh` tasks are skipped). `compare.sh --profile`
  also uses it.
- `probes/<name>/` — synthetic `func main()` perf probes for `compare.sh`
  (no target checkout needed; output doubles as a correctness check when
  two builds disagree).
- `reports/` — measurement rounds as self-contained HTML (open in a browser).

Verdicts: `PASS`, `DIFF` (silent divergence — a bug), `TRAP`, `REJECT`,
`HANG` (timeout), `ORACLE-FAIL`, `SETUP-FAIL`.

## Tasks

| task | target | reads | checks against |
|---|---|---|---|
| grafana-openapi | grafana | route registrations inside `registerRoutes`' body (inspect `Body`/`Nodes`) + `// swagger:route` annotations (117 method docs via `Doc`, 4 free via `Comments`), via `inspect` only | each other |
| grafana-coreplugin | grafana | `coreplugin` plugin-ID consts, via `go/parser` (no package init) | `public/app/plugins/datasource/*/plugin.json` |
| clickhouse-settings | clickhouse-datasource | `Settings` struct json tags, via `inspect` (surface only) | `CHConfig` in `src/types/config.ts` |
| oapi-codegen-examples | oapi-codegen | every `go:generate` line in `examples/` (53), running `cmd/oapi-codegen` itself under minigo (`--src text/template,encoding/json,…`) via `task.sh` | the natively built `cmd/oapi-codegen` (rc + sha256 of every written file) |
| grafana-swagger-spec | grafana | `swagger:route`/`response`/`model` doc annotations in `pkg/api` + `pkg/api/dtos`, via `inspect` (surface only, annotation → spec direction) | `public/api-merged.json` |
| gosimilarity-scan | gosimilarity | the `gosimilarity` CLI itself under minigo (fnv-for-maphash patch; real flagx), scanning `input/` — a pinned 5-file slice of podhmo/minigo (~660 functions) — text and `--format json` reports | the natively built gosimilarity on the same patched tree and input |

Run a single task by name: `./run.sh grafana-openapi`, `./run.sh
oapi-codegen-examples`. oapi-codegen-examples is a compatibility/perf
workload, not the pitch: code generation executes nearly the whole program
(kin-openapi, text/template, json/v2, goimports), so minigo's lazy partial
loading does not help and the run shows raw interpreter speed (~5 min for the
53 lines vs ~7 s native). It is not part of `compare.sh`.

gosimilarity-scan is the same kind of perf workload: the interpreted run
executes go/ast normalization plus an O(n²) tree edit distance (run twice:
text and JSON reports), ~90 s per run vs ~0.6 s for the built native binary on the pinned input (~160x; ~110
reported pairs; `oracle_s` also counts rsync + `go build`). The one local
patch it carries (`shim.patch`, fnv for `hash/maphash`) exists only because
`maphash.MakeSeed` hangs when interpreted — see podhmo/minigo TODO.md,
`docs/sketch/ja/experiment-gosimilarity.md` and
`docs/sketch/ja/experiment-gosimilarity-reverify.md`. `PROFILE=1 ./run.sh
gosimilarity-scan` profiles it through the minigo CLI. Because the task
runs minigo twice, the profiles are those of the JSON run;
`out/gosimilarity-scan.prof.txt` leads with the script lines
(`distance.func1` ~86%). Not part of `compare.sh`.

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

- Open a profile with `go tool pprof -http=: out/minigo out/<task>.cpu.pprof`
  (`out/minigo` is the binary that produced it; `out/prof` under the
  `prof/` fallback). The script profile needs no binary:
  `go tool pprof -top -lines out/<task>.script.pprof` names script
  functions and `file:line` — start there to find the hot script code,
  then read the cpu/allocs profiles for what the interpreter spends on it.
- **macOS:** heavy goroutine spawning (e.g. a helper goroutine per
  `sync.Mutex.Lock`) can make the CPU profile pin most samples on libc
  (`pthread_cond_signal`/`pthread_cond_wait`) and hide the interpreter.
  `TRACE=1` adds `out/<task>.trace`; `go tool trace` shows whether the main
  goroutine was actually running. If it was, trust the alloc profile and
  re-profile on Linux or with the spawning path disabled.

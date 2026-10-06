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

Verdicts: `PASS`, `DIFF` (silent divergence — a bug), `TRAP`, `REJECT`,
`HANG` (timeout), `ORACLE-FAIL`, `SETUP-FAIL`.

## Tasks

| task | target | reads | checks against |
|---|---|---|---|
| grafana-openapi | grafana | route registrations in `pkg/api/api.go` (func bodies) + `// swagger:route` comments, via `go/parser` | each other |
| grafana-coreplugin | grafana | `coreplugin` plugin-ID consts, via `go/parser` (no package init) | `public/app/plugins/datasource/*/plugin.json` |
| clickhouse-settings | clickhouse-datasource | `Settings` struct json tags, via `inspect` (surface only) | `CHConfig` in `src/types/config.ts` |

## Adding a task

Pick a sync a real project needs (Go code ↔ another artifact). Prefer
inputs `inspect` can reach (surface) — when a task has to fall back to
`go/parser` because `inspect` can't see something, note *what* it couldn't
see; that gap is a finding. Keep `go run`-able when possible so the oracle
is free.

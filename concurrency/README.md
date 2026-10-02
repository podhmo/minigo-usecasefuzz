# concurrency fuzz

Thought-guided fuzz corpus for `minigo`'s real-concurrency implementation
(host goroutines, `chan Value`, `reflect.Select`, `sync`/`time` intrinsics,
the per-run `proc` model). Unlike the top-level `usecasefuzz` corpus, this
harness drives **exported functions per case** (`func X() int`), not `main`
— each theme dir holds a suite of probes, and panics/deadlocks/timeouts are
first-class expected outcomes, not just diffs.

Findings from the original round (11 bugs, all fixed):
`docs/sketch/ja/fuzz-concurrency.md` in the minigo repo (PR #28).

## Layout

- `cases/<theme>/main.go` — `func ExportedName() int` probes. `// want`
  comments mark expected values; some cases *expect* a fatal
  deadlock panic (host `all goroutines are asleep`) or a recoverable panic.
- `run.sh [dir ...]` — builds `out/minigo` from `$MINIGO_DIR` (default:
  sibling clone `../minigo`, auto-cloned if missing), runs every exported
  func with `timeout 8`, prints `=== dir/Func (exit N)` + output.
- `host/` — host-side Go checks: `runtime.ImportRef.Materialize` panic
  bookkeeping, concurrent `Engine.Run` on one engine (`-race`), goroutine
  leak measurement (`NumGoroutine` before/after).
  `go.mod` uses `replace github.com/podhmo/minigo => ../../../minigo`
  (the same sibling-clone convention as `run.sh`).

## Usage

```sh
./run.sh                 # all themes
./run.sh sync2 time2     # selected themes
MINIGO_DIR=/path/to/minigo ./run.sh   # run against a specific checkout

cd host && go run .      # host-side checks
cd host && go run -race .
```

## Exit-code oracle

| exit | meaning |
|---|---|
| 0 | ran, printed value — compare against the `// want` note |
| 1 | script panic/trap (expected for e.g. `RangeTwoVars`, close-of-closed) |
| 2 | host fatal (deadlock detection, unlock of unlocked mutex — Go-consistent) |
| 124 | timeout — silent hang, always a bug |

## Themes

| dir | surface exercised |
|---|---|
| chan_extra | misc chan shapes (elem types, named chan types, struct-field chans) |
| chan_kind | chan across `Cell`/`Named`/`IfaceNil`/`TypedNil`/`GoValue` value classes |
| chanofchan | `chan chan int` (main-only probe) |
| close_cases | close/close-of-closed/send-on-closed/recv-after-close semantics |
| deadlock | deliberate deadlocks — host `fatal error` expected, not a hang |
| defer_exit | defers draining through procExit teardown |
| gspawn | `go` stmt forms (args, closures, fns as values) |
| loopvar | per-iteration loop vars captured by goroutines (Go 1.22 semantics) |
| misc2 | assorted semantics probes |
| neg | negative-size/buffer edge cases (`make(chan, -1)` etc.) |
| once | `sync.Once` |
| rangex | `for range` over channels incl. the 2-var trap |
| recvbinds | recv into every binding shape (`s.f`, `a[i]`, `m[k]`, `_`, `v,ok`) |
| sel_chan2 | select over unusual chan positions |
| sel_misc | nested/loop/goroutine/expression selects |
| selassign | select cases with `x = <-ch` assign/send arms |
| sync2 | `sync` across goroutines (WaitGroup/Mutex/RWMutex/Once shared) |
| time2 | `time.After`/`NewTimer`/`NewTicker`/`Sleep` incl. select arms |
| timeriso | Timer Stop/Reset semantics on the real clock |

## Expected failures still standing

- `host` T5 / `DetachedWait`: a script goroutine parked inside a
  non-select host call (`WaitGroup.Wait`, `Mutex.Lock`) outlives its
  process — goroutine count grows by 1. Go-consistent limitation; only
  select-based blocking watches `proc.done`.
- Deadlock dirs exit 2 with the host's `fatal error` — that *is* the
  correct Go-like behavior (the CLI dies rather than hanging).

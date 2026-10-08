#!/usr/bin/env bash
# run.sh [task ...] — realworld harness: sync-tool scripts run against pinned
# real-world codebases, `go run` (oracle) vs `minigo run`, with wall times.
#
# Each tasks/<name>/ is a standalone module (func main) reading the target
# checkout from $TARGET_DIR. `target` names a row of targets.tsv. A task with
# want.txt is minigo-only (e.g. imports minigo.dev/inspect) and is compared
# against that golden file instead of `go run`. A task with an executable
# task.sh drives a whole program instead of being one: the oracle is
# `task.sh native`, the minigo side `task.sh minigo <BIN>` (see
# tasks/oapi-codegen-examples). A `timeout_sec` file overrides TIMEOUT_SEC.
#
#   MINIGO_DIR  — podhmo/minigo checkout (default: sibling clone)
#   SRC_DIR     — where targets are cloned (default: ./.src, gitignored)
#   TIMEOUT_SEC — per minigo run (default: 120)
#   COLD=1      — also time the oracle with an empty GOCACHE (slow)
#   PROFILE=1   — also run each task under prof/ (runtime/pprof around the
#                 engine) and write out/<task>.{cpu,allocs}.pprof plus a
#                 text summary out/<task>.prof.txt
#   TRACE=1     — with PROFILE=1, also write out/<task>.trace (on macOS the
#                 CPU profile can misattribute samples; the trace shows
#                 what the main goroutine really did)
#
# GOCACHE is never cleared: cold timings use a throwaway temp dir.
set -u
# Drop inherited GIT_* (e.g. GIT_DIR under `git bisect run` or a hook):
# fetch_target's `git init`/`fetch` would otherwise hit the caller's repo.
unset $(compgen -e GIT_)
ROOT="$(cd "$(dirname "$0")" && pwd)"
MINIGO_DIR="${MINIGO_DIR:-$ROOT/../../minigo}"
SRC_DIR="${SRC_DIR:-$ROOT/.src}"
TIMEOUT_SEC="${TIMEOUT_SEC:-120}"
COLD="${COLD:-0}"
PROFILE="${PROFILE:-0}"
TRACE="${TRACE:-0}"
OUT="$ROOT/out"
BIN="$OUT/minigo"
TIMEOUT=timeout
command -v timeout > /dev/null || TIMEOUT=gtimeout
mkdir -p "$OUT" "$SRC_DIR"

now() { perl -MTime::HiRes=time -e 'printf "%.2f", time'; }
elapsed() { perl -e "printf '%.2f', $2 - $1"; }

if [ ! -d "$MINIGO_DIR" ]; then
	echo "cloning podhmo/minigo into $MINIGO_DIR ..." >&2
	git clone --depth 1 https://github.com/podhmo/minigo "$MINIGO_DIR" || {
		echo "cannot obtain minigo; set MINIGO_DIR to your checkout" >&2
		exit 1
	}
fi
(cd "$MINIGO_DIR" && go build -o "$BIN" ./cmd/minigo) || exit 1

# build_prof — compile prof/ against $MINIGO_DIR through a generated module.
build_prof() {
	local mdir pdir="$OUT/prof-build"
	mdir="$(cd "$MINIGO_DIR" && pwd)"
	rm -rf "$pdir" && mkdir -p "$pdir" && cp "$ROOT/prof/main.go" "$pdir/" &&
		cp "$mdir/go.sum" "$pdir/" &&
		printf 'module realworld/prof\n\ngo 1.26\n\nrequire github.com/podhmo/minigo v0.0.0\n\nreplace github.com/podhmo/minigo => %s\n' "$mdir" > "$pdir/go.mod" &&
		(cd "$pdir" && GOFLAGS=-mod=mod go build -o "$OUT/prof" .)
}

# summarize_prof NAME — top-N text views next to the raw profiles.
summarize_prof() {
	local name="$1" base="$OUT/$1"
	{
		echo "## cpu: flat"
		go tool pprof -top -nodecount=25 "$OUT/prof" "$base.cpu.pprof" 2>/dev/null | sed -n '4,$p'
		echo
		echo "## cpu: cumulative (minigo frames)"
		go tool pprof -top -cum -nodecount=200 "$OUT/prof" "$base.cpu.pprof" 2>/dev/null | grep 'podhmo/minigo' | head -30
		echo
		echo "## allocs: alloc_space"
		go tool pprof -top -nodecount=20 -sample_index=alloc_space "$OUT/prof" "$base.allocs.pprof" 2>/dev/null | sed -n '4,$p'
	} > "$base.prof.txt"
}

[ "$TRACE" = 1 ] && PROFILE=1
if [ "$PROFILE" = 1 ]; then
	build_prof || { echo "cannot build prof/ against $MINIGO_DIR" >&2; exit 1; }
fi

# fetch_target NAME — shallow-fetch the pinned commit and download modules once.
fetch_target() {
	local name="$1" url sha dir
	read -r url sha < <(awk -F'\t' -v n="$name" '$1 == n { print $2, $3 }' "$ROOT/targets.tsv")
	[ -n "${sha:-}" ] || { echo "?? target $name: not in targets.tsv" >&2; return 1; }
	dir="$SRC_DIR/$name"
	if [ "$(git -C "$dir" rev-parse HEAD 2>/dev/null)" != "$sha" ]; then
		echo "fetching $name@${sha:0:7} ..." >&2
		rm -rf "$dir" && mkdir -p "$dir" &&
			git -C "$dir" init -q &&
			git -C "$dir" fetch -q --depth 1 "$url" "$sha" &&
			git -C "$dir" checkout -q FETCH_HEAD || return 1
		rm -f "$dir/.modules-downloaded"
	fi
	if [ ! -f "$dir/.modules-downloaded" ]; then
		echo "go mod download in $name ..." >&2
		(cd "$dir" && go mod download) && touch "$dir/.modules-downloaded" || return 1
	fi
}

if [ $# -gt 0 ]; then
	TASKS="$*"
else
	TASKS="$(cd "$ROOT/tasks" && ls -d */ | tr -d '/')"
fi

printf '%-12s %-24s %9s %9s %9s\n' VERDICT TASK oracle_s minigo_s cold_s
for name in $TASKS; do
	dir="$ROOT/tasks/$name"
	[ -d "$dir" ] || { echo "?? $name: no such task"; continue; }
	target="$(cat "$dir/target")"
	fetch_target "$target" || { printf '%-12s %s\n' SETUP-FAIL "$name"; continue; }
	export TARGET_DIR="$SRC_DIR/$target"

	cold="-"
	tsec="$TIMEOUT_SEC"
	[ -f "$dir/timeout_sec" ] && tsec="$(cat "$dir/timeout_sec")"
	if [ -x "$dir/task.sh" ]; then
		t0=$(now)
		(cd "$dir" && ./task.sh native) > "$OUT/$name.want" 2>&1
		want_rc=$?
		oracle=$(elapsed "$t0" "$(now)")
		if [ "$COLD" = 1 ]; then
			cache="$(mktemp -d)"
			t0=$(now)
			(cd "$dir" && GOCACHE="$cache" ./task.sh native) > /dev/null 2>&1
			cold=$(elapsed "$t0" "$(now)")
			rm -rf "$cache"
		fi
	elif [ -f "$dir/want.txt" ]; then
		cp "$dir/want.txt" "$OUT/$name.want"
		want_rc=0
		oracle="-"
	else
		(cd "$dir" && go build -o "$OUT/$name.native" .) > "$OUT/$name.want" 2>&1
		t0=$(now)
		"$OUT/$name.native" > "$OUT/$name.want" 2>&1
		want_rc=$?
		oracle=$(elapsed "$t0" "$(now)")
		if [ "$COLD" = 1 ]; then
			cache="$(mktemp -d)"
			t0=$(now)
			(cd "$dir" && GOCACHE="$cache" go run .) > /dev/null 2>&1
			cold=$(elapsed "$t0" "$(now)")
			rm -rf "$cache"
		fi
	fi

	t0=$(now)
	if [ -x "$dir/task.sh" ]; then
		(cd "$dir" && $TIMEOUT "$tsec" ./task.sh minigo "$BIN") > "$OUT/$name.got" 2>&1
	else
		(cd "$dir" && $TIMEOUT "$tsec" "$BIN" run .) > "$OUT/$name.got" 2>&1
	fi
	got_rc=$?
	got_s=$(elapsed "$t0" "$(now)")

	if [ "$got_rc" -eq 124 ]; then
		v="HANG"
	elif [ "$want_rc" -ne 0 ]; then
		v="ORACLE-FAIL"
	elif [ "$got_rc" -ne 0 ]; then
		if grep -q 'runtime trap\|OpTrap\|panic:' "$OUT/$name.got"; then
			v="TRAP"
		else
			v="REJECT"
		fi
	elif diff -q "$OUT/$name.want" "$OUT/$name.got" > /dev/null; then
		v="PASS"
	else
		v="DIFF"
	fi
	printf '%-12s %-24s %9s %9s %9s\n' "$v" "$name" "$oracle" "$got_s" "$cold"

	if [ "$PROFILE" = 1 ] && [ -x "$dir/task.sh" ]; then
		printf '%-12s %-24s %s\n' "" "  profile" "skipped (task.sh drives a program, not one main)"
	elif [ "$PROFILE" = 1 ]; then
		targs=""
		[ "$TRACE" = 1 ] && targs="-trace"
		res=$(cd "$dir" && $TIMEOUT "$TIMEOUT_SEC" "$OUT/prof" -dir . -out "$OUT/$name" $targs 2> "$OUT/$name.prof.err" | tail -1)
		summarize_prof "$name"
		printf '%-12s %-24s %s -> out/%s.prof.txt\n' "" "  profile" "$res" "$name"
	fi
done

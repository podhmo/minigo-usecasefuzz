#!/usr/bin/env bash
# run.sh [case ...] — convert-define fuzz: generate converters per case, then
# `go build` the output inside the case's own module (the oracle).
#
# Each cases/<name>/ holds: define.go (//go:build codegen), source/,
# destination/ (+ any extra pkgs), expect (OK|GEN-FAIL|BUILD-FAIL).
# go.mod is regenerated per run so the replace points at GEN_DIR.
#
#   MINIGO_DIR — podhmo/minigo checkout (default: sibling clone — if
#                missing, it is cloned shallowly next to this repo)
set -u
ROOT="$(cd "$(dirname "$0")" && pwd)"
MINIGO_DIR="${MINIGO_DIR:-$ROOT/../../minigo}"
REPO_DIR="$MINIGO_DIR"
GEN_DIR="$MINIGO_DIR/examples/convert-define"

if [ ! -d "$MINIGO_DIR" ]; then
	echo "cloning podhmo/minigo into $MINIGO_DIR ..." >&2
	git clone --depth 1 https://github.com/podhmo/minigo "$MINIGO_DIR" || {
		echo "cannot obtain minigo; set MINIGO_DIR to your checkout" >&2
		exit 1
	}
fi
OUT="$ROOT/out"
mkdir -p "$OUT"

if [ $# -gt 0 ]; then
	CASES="$*"
else
	CASES="$(cd "$ROOT/cases" && ls -d */ 2>/dev/null | tr -d '/')"
fi

pass=0; fail=0
for name in $CASES; do
	dir="$ROOT/cases/$name"
	[ -d "$dir" ] || { echo "?? $name: no such case"; continue; }
	expect="$(cat "$dir/expect" 2>/dev/null || echo OK)"

	cat > "$dir/go.mod" <<MOD
module example.com/m

go 1.24

require github.com/podhmo/minigo/examples/convert-define v0.0.0
require github.com/podhmo/minigo v0.0.0

replace github.com/podhmo/minigo => $REPO_DIR
replace github.com/podhmo/minigo/examples/convert-define => $GEN_DIR
MOD

	rm -f "$dir/generated.go" "$dir/go.sum"
	gen_log="$OUT/$name.gen.log"; build_log="$OUT/$name.build.log"
	verdict=""
	(cd "$dir" && go mod tidy) >"$gen_log" 2>&1 || verdict="TIDY-FAIL"
	if [ -z "$verdict" ]; then
		(cd "$GEN_DIR" && go run . -file "$dir/define.go" -output "$dir/generated.go") >>"$gen_log" 2>&1 \
			|| verdict="GEN-FAIL"
	fi
	if [ -z "$verdict" ] && [ ! -s "$dir/generated.go" ]; then
		verdict="GEN-FAIL" # generator "succeeded" but wrote nothing
	fi
	# second tidy: generated.go adds imports (model, convutil, ...) the first
	# pass could not see
	if [ -z "$verdict" ]; then
		(cd "$dir" && go mod tidy) >>"$gen_log" 2>&1 || verdict="TIDY-FAIL"
	fi
	if [ -z "$verdict" ]; then
		(cd "$dir" && go build ./...) >"$build_log" 2>&1 || verdict="BUILD-FAIL"
	fi
	[ -z "$verdict" ] && verdict="OK"

	if [ "$verdict" = "$expect" ]; then
		echo "PASS $name ($verdict)"
		pass=$((pass + 1))
	else
		echo "FAIL $name: got $verdict, want $expect (logs: out/$name.*.log)"
		fail=$((fail + 1))
	fi
done
echo "== $pass expected / $fail unexpected =="
[ $fail -eq 0 ]

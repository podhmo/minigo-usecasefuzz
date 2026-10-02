#!/usr/bin/env bash
# run.sh [case ...] — differential test: `go run` (oracle) vs `minigo run`.
# Each cases/<name>/ dir is a standalone module so both run from the dir itself.
set -u
ROOT="$(cd "$(dirname "$0")" && pwd)"
MINIGO_DIR="${MINIGO_DIR:-$HOME/repos/minigo}"
OUT="$ROOT/out"
BIN="$OUT/minigo"
mkdir -p "$OUT"

if [ ! -x "$BIN" ] || [ "$MINIGO_DIR/intrinsics.go" -nt "$BIN" ]; then
	(cd "$MINIGO_DIR" && go build -o "$BIN" ./cmd/minigo) || exit 1
fi

if [ $# -gt 0 ]; then
	CASES="$*"
else
	CASES="$(cd "$ROOT/cases" && ls -d */ | tr -d '/')"
fi

for name in $CASES; do
	dir="$ROOT/cases/$name"
	[ -d "$dir" ] || { echo "?? $name: no such case"; continue; }

	(cd "$dir" && go run .) > "$OUT/$name.want" 2>&1
	want_rc=$?
	(cd "$dir" && timeout 15 "$BIN" run .) > "$OUT/$name.got" 2>&1
	got_rc=$?

	if [ "$want_rc" -ne 0 ]; then
		if [ "$got_rc" -ne 0 ]; then
			v="PASS-REJECT" # both refuse
		else
			v="ACCEPT"      # minigo ran what Go rejects -- worst
		fi
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
	printf '%-14s %s (go=%d minigo=%d)\n' "$v" "$name" "$want_rc" "$got_rc"
done

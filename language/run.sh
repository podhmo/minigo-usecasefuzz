#!/usr/bin/env bash
# run.sh [case ...] — differential harness for Go-language edge cases:
# `go run` (the oracle) vs `minigo run`, stdout+stderr diffed per case.
# Cases live in the shared `langfuzz` module (./go.mod covers cases/*/).
#
# Point MINIGO_DIR at a podhmo/minigo checkout (default: sibling clone —
# if missing, it is cloned shallowly next to this repo).
set -u
ROOT="$(cd "$(dirname "$0")" && pwd)"
MINIGO_DIR="${MINIGO_DIR:-$ROOT/../../minigo}"
OUT="$ROOT/out"
case "$(uname -s)" in
MINGW* | MSYS* | CYGWIN*)
	BIN="$OUT/minigo.exe"
	TIMEOUT=/usr/bin/timeout # System32\timeout is a different command
	;;
*)
	BIN="$OUT/minigo"
	TIMEOUT=timeout
	;;
esac
mkdir -p "$OUT"

if [ ! -d "$MINIGO_DIR" ]; then
	echo "cloning podhmo/minigo into $MINIGO_DIR ..." >&2
	git clone --depth 1 https://github.com/podhmo/minigo "$MINIGO_DIR" || {
		echo "cannot obtain minigo; set MINIGO_DIR to your checkout" >&2
		exit 1
	}
fi

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

	go_out="$OUT/$name.go.out"; go_err="$OUT/$name.go.err"
	mg_out="$OUT/$name.mg.out"; mg_err="$OUT/$name.mg.err"

	(cd "$dir" && $TIMEOUT 15 go run .) >"$go_out" 2>"$go_err"
	gec=$?
	(cd "$dir" && $TIMEOUT 15 "$BIN" run .) >"$mg_out" 2>"$mg_err"
	mec=$?

	# classify each side: ok / panic (died with a runtime panic) / reject / timeout
	gokind=ok
	if [ $gec -eq 124 ]; then gokind=timeout
	elif [ $gec -ne 0 ]; then
		if grep -qE 'panic:|fatal error' "$go_err"; then gokind=panic
		else gokind=reject; fi
	fi
	mkind=ok
	if [ $mec -eq 124 ]; then mkind=timeout
	elif [ $mec -ne 0 ]; then mkind=die; fi

	verdict="?"
	if [ "$gokind" = ok ] && [ "$mkind" = ok ]; then
		if cmp -s "$go_out" "$mg_out" && cmp -s "$go_err" "$mg_err"; then verdict=PASS
		else verdict=DIFF; fi
	elif [ "$gokind" = ok ] && [ "$mkind" != ok ]; then verdict=TRAP
	elif [ "$gokind" = panic ] && [ "$mkind" = die ]; then verdict=PASS-PANIC
	elif [ "$gokind" = panic ] && [ "$mkind" = ok ]; then verdict=NOPANIC
	elif [ "$gokind" = reject ] && [ "$mkind" = die ]; then verdict=PASS-REJECT
	elif [ "$gokind" = reject ] && [ "$mkind" = ok ]; then verdict=ACCEPT
	elif [ "$gokind" = timeout ] || [ "$mkind" = timeout ]; then verdict=TIMEOUT
	else verdict=OTHER; fi
	printf '%-12s %-16s go=%-7s minigo=%-6s\n' "$verdict" "$name" "$gokind" "$mkind"
done

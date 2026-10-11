#!/usr/bin/env bash
# run.sh [case ...] — extract-schema corpus: build the tool once, then run it
# inside each case's own module so both the script packages and the fixture's
# imports resolve through the case go.mod's `replace => $MINIGO_DIR`.
#
# Each cases/<name>/ is a module root (go.mod is regenerated per run):
#   *.go          fixture package (target "." by default)
#   sub/...       sub-packages the fixture imports (example.com/m/sub)
#   plugin*/      optional plugin script dirs (-plugin ./plugin*)
#   target        optional target arg, default "."
#   type          optional type name appended after the target
#   flags         optional extra flags, one line (e.g. "-stdlib walk")
#   expect        expected verdict: OK|WARN|DIFF|BADJSON|FAIL|TRAP|PANIC
#
#   TRAP = 'runtime trap' (engine-level fault); PANIC = 'panic:' (script
#   panic propagated); FAIL = a graceful error/usage exit; WARN = exit 0
#   with stderr output (the tool's own softened-failure channel).
#   want          expected stdout — the hand-verified schema (OK/DIFF cases)
#
# Captured output lands in out/ (committed): <name>.stdout / <name>.stderr,
# sanitized of local paths so runs on other machines stay diff-stable.
#
#   MINIGO_DIR — podhmo/minigo checkout (default: sibling clone — if
#                missing, it is cloned shallowly next to this repo)
set -u
ROOT="$(cd "$(dirname "$0")" && pwd)"
MINIGO_DIR="${MINIGO_DIR:-$ROOT/../minigo}"
GEN_DIR="$MINIGO_DIR/examples/extract-schema"
OUT="$ROOT/out"
BIN="$OUT/extract-schema"
case "$(uname -s)" in
MINGW* | MSYS* | CYGWIN*)
	BIN="$BIN.exe"
	TIMEOUT=/usr/bin/timeout
	;;
*)
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

# rebuild when any host file is newer than the binary (the tool is small —
# rebuilding unconditionally keeps the check simple and stays correct when
# script/ or inspect internals change)
(cd "$GEN_DIR" && go build -o "$BIN" .) || exit 1

if [ $# -gt 0 ]; then
	CASES="$*"
else
	CASES="$(cd "$ROOT/cases" && ls -d */ 2>/dev/null | tr -d '/')"
fi

# sanitize writes stable output: local paths become tokens so committed
# out/ files do not thrash between machines.
sanitize() {
	sed -e "s|$MINIGO_DIR|MINIGO|g" -e "s|$ROOT|CORPUS|g" -e "s|$HOME|HOME|g"
}

pass=0; fail=0
for name in $CASES; do
	dir="$ROOT/cases/$name"
	[ -d "$dir" ] || { echo "?? $name: no such case"; continue; }
	expect="$(cat "$dir/expect" 2>/dev/null || echo OK)"

	cat > "$dir/go.mod" <<MOD
module example.com/m

go 1.24

require github.com/podhmo/minigo v0.0.0

replace github.com/podhmo/minigo => $MINIGO_DIR
MOD

	ln -sfn "$GEN_DIR/script" "$dir/script"

	target="$(cat "$dir/target" 2>/dev/null || echo .)"
	type="$(cat "$dir/type" 2>/dev/null || true)"
	flags="$(cat "$dir/flags" 2>/dev/null || true)"

	stdout="$OUT/$name.stdout"; stderr="$OUT/$name.stderr"
	verdict=""
	# flags is a shell fragment (quoted multi-word values allowed) — corpus-
	# authored input, eval'd inside the case dir
	if [ -n "$flags" ]; then
		(cd "$dir" && eval "set -- $flags" && $TIMEOUT 60 "$BIN" "$@" "$target" $type) \
			> "$stdout.raw" 2> "$stderr.raw"
	else
		# shellcheck disable=SC2086 # type is intentionally word-split
		(cd "$dir" && $TIMEOUT 60 "$BIN" "$target" $type) \
			> "$stdout.raw" 2> "$stderr.raw"
	fi
	rc=$?
	sanitize < "$stdout.raw" > "$stdout"; rm -f "$stdout.raw"
	sanitize < "$stderr.raw" > "$stderr"; rm -f "$stderr.raw"

	if [ "$rc" -ne 0 ]; then
		if grep -q '^runtime trap' "$stdout" "$stderr"; then
			verdict="TRAP"
		elif grep -q '^panic' "$stdout" "$stderr"; then
			verdict="PANIC"
		else
			verdict="FAIL"
		fi
	elif ! jq -e . "$stdout" >/dev/null 2>&1; then
		verdict="BADJSON"
	elif [ -f "$dir/want" ] && ! diff -q "$dir/want" "$stdout" >/dev/null; then
		verdict="DIFF"
	elif [ -s "$stderr" ]; then
		verdict="WARN"
	else
		verdict="OK"
	fi

	if [ "$verdict" = "$expect" ]; then
		echo "PASS $name ($verdict)"
		pass=$((pass + 1))
	else
		echo "FAIL $name: got $verdict, want $expect (out/$name.*)"
		fail=$((fail + 1))
	fi
done
echo "== $pass expected / $fail unexpected =="
[ $fail -eq 0 ]

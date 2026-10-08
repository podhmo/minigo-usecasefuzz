#!/usr/bin/env bash
# task.sh native|minigo [BIN] — run every `go:generate go run .../cmd/oapi-codegen`
# line under $TARGET_DIR/examples, in a scratch copy of the checkout.
#
# native: builds cmd/oapi-codegen inside the examples module (what `go run`
# in a go:generate line executes) and runs it per line.
# minigo: `BIN run --src ... <pkg> -- ARGS` per line.
#
# Prints one block per line: rc and the sha256 of every file the run wrote,
# so the two sides compare as text. Exits 1 when any line failed.
#
# oapi-codegen runs nearly the whole program (kin-openapi, text/template,
# json/v2, goimports): a compatibility/perf workload, not minigo's lazy
# partial-read use case.
set -u
MODE="${1:?usage: task.sh native|minigo [BIN]}"
BIN="${2:-}"
: "${TARGET_DIR:?TARGET_DIR is not set}"
PKG=github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
# The bound text/template, encoding/json, … are too partial for
# oapi-codegen; these run from source (podhmo/minigo TODO.md).
SRC_FLAGS="--src text/template,encoding/json,encoding/hex,encoding/base64,encoding/base32,slices,maps"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
rsync -a --exclude .git "$TARGET_DIR/" "$WORK/src/" || exit 1
cd "$WORK/src/examples" || exit 1

case "$MODE" in
native)
	go build -o "$WORK/oapi-codegen" "$PKG" || exit 1
	run() { "$WORK/oapi-codegen" "$@"; } ;;
minigo)
	[ -x "$BIN" ] || { echo "task.sh: minigo binary required" >&2; exit 1; }
	run() { "$BIN" run $SRC_FLAGS "$PKG" -- "$@"; } ;;
*)
	echo "task.sh: unknown mode $MODE" >&2; exit 1 ;;
esac

failed=0 n=0
while IFS=: read -r file line rest; do
	dir="$(dirname "${file#./}")"
	args="${rest#*cmd/oapi-codegen }"
	n=$((n + 1))
	marker="$WORK/marker"
	: > "$marker"
	(cd "$dir" && eval run "$args") > "$WORK/log" 2>&1
	rc=$?
	echo "== $dir: $args rc=$rc"
	if [ "$rc" -ne 0 ]; then
		failed=1
		grep -m1 -v '^[[:space:]]' "$WORK/log" | sed 's/^/  error: /'
	fi
	(cd "$dir" && find . -type f -newer "$marker" | sort | while read -r f; do
		echo "  wrote ${f#./} $(shasum -a 256 "$f" | cut -c1-16)"
	done)
done < <(grep -rn --include='*.go' "go:generate go run $PKG " . | sort -t: -k1,1 -k2,2n)

[ "$n" -gt 0 ] || { echo "task.sh: no go:generate lines found" >&2; exit 1; }
exit "$failed"

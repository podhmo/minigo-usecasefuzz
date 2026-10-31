#!/usr/bin/env bash
# task.sh native|minigo [BIN] — run gosimilarity (structural AST similarity,
# https://github.com/Songmu/gosimilarity) over input/, a pinned 5-file slice
# of podhmo/minigo source (~660 functions).
#
# Both sides run the SAME patched copy of the pinned gosimilarity checkout.
# shim.patch swaps hash/maphash for hash/fnv: maphash.MakeSeed loops
# forever when interpreted (its bodiless `//go:linkname runtime.rand` hits
# minigo's zero-return shim and keeps answering 0). The real
# github.com/Songmu/flagx runs as-is since runtime.AddCleanup is bound
# (podhmo/minigo#1023); its weak.Make silently yields an empty pointer under
# minigo, harmless here because flagx only uses it as a map key for the one
# FlagSet. Both gaps are tracked in podhmo/minigo TODO.md.
#
# Output is the default text report: `--format json` still diverges
# (minigo's json.Marshal emits unexported fields — TODO.md).
#
# Perf workload, not the lazy-read pitch: interpretation covers nearly the
# whole tool (go/ast normalization + O(n^2) tree edit distance). Reference
# numbers (docs/sketch/ja/experiment-gosimilarity-reverify.md): ~0.6s for
# the built native binary vs ~90s interpreted (~160x; oracle_s also counts
# rsync + go build), identical 110 pairs.
set -u
MODE="${1:?usage: task.sh native|minigo [BIN]}"
BIN="${2:-}"
: "${TARGET_DIR:?TARGET_DIR is not set}"
TASK_DIR="$(cd "$(dirname "$0")" && pwd)"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
rsync -a --exclude .git "$TARGET_DIR/" "$WORK/tool/" || exit 1
mkdir -p "$WORK/tool/input" && cp "$TASK_DIR"/input/*.go "$WORK/tool/input/"
(cd "$WORK/tool" && git apply "$TASK_DIR/shim.patch") || exit 1
cd "$WORK/tool" || exit 1

case "$MODE" in
native)
	go build -o "$WORK/gosimilarity" ./cmd/gosimilarity || exit 1
	run() { "$WORK/gosimilarity" "$@"; } ;;
minigo)
	[ -x "$BIN" ] || { echo "task.sh: minigo binary required" >&2; exit 1; }
	run() { "$BIN" run ./cmd/gosimilarity -- "$@"; } ;;
*)
	echo "task.sh: unknown mode $MODE" >&2; exit 1 ;;
esac

run input

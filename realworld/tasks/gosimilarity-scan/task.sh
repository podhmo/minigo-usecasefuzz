#!/usr/bin/env bash
# task.sh native|minigo [BIN] — run gosimilarity (structural AST similarity,
# https://github.com/Songmu/gosimilarity) over input/, a pinned 5-file slice
# of podhmo/minigo source (~660 functions).
#
# Both sides run the SAME patched copy of the pinned gosimilarity checkout:
# flagshim.go replaces github.com/Songmu/flagx (flagx calls
# runtime.AddCleanup, which minigo does not bind) and shim.patch swaps
# hash/maphash for hash/fnv (maphash.MakeSeed hangs when interpreted; both
# gaps are tracked in podhmo/minigo TODO.md).
#
# The scan input must stay outside any .git: gosimilarity's git-root
# discovery calls exec.CommandContext, which minigo does not bind.
#
# Perf workload, not the lazy-read pitch: interpretation covers nearly the
# whole tool (go/ast normalization + O(n^2) tree edit distance). Reference
# numbers (report): ~1.0s native vs ~160s interpreted, identical 110 pairs.
set -u
MODE="${1:?usage: task.sh native|minigo [BIN]}"
BIN="${2:-}"
: "${TARGET_DIR:?TARGET_DIR is not set}"
TASK_DIR="$(cd "$(dirname "$0")" && pwd)"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
rsync -a --exclude .git "$TARGET_DIR/" "$WORK/tool/" || exit 1
cp "$TASK_DIR/flagshim.go" "$WORK/tool/flagshim.go"
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

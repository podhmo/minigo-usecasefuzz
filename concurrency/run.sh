#!/usr/bin/env bash
# run.sh [dir ...] — run every exported Func in each case dir through minigo.
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

cd "$ROOT"
dirs=("$@")
if [ ${#dirs[@]} -eq 0 ]; then
	dirs=($(ls -d cases/*/ | xargs -n1 basename))
fi
for d in "${dirs[@]}"; do
	fns=$(grep -oE '^func [A-Z][A-Za-z0-9_]*\(' "cases/$d/main.go" | sed 's/^func //; s/($//; s/(//')
	for fn in $fns; do
		out=$($TIMEOUT 8 "$BIN" "./cases/$d" "$fn" 2>&1)
		code=$?
		printf '=== %s/%s (exit %s)\n%s\n' "$d" "$fn" "$code" "$out"
	done
done

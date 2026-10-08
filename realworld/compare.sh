#!/usr/bin/env bash
# compare.sh OLD_REF NEW_REF [options] — honest A/B performance comparison of
# two minigo revisions, measured in the same sitting (absolute numbers drift
# between machines and days; only same-sitting deltas mean anything).
#
# Each side is checked out into a scratch worktree under $MINIGO_DIR and built.
# Probes run interleaved (A B A B …, RUNS rounds each) and medians are compared.
# Exit code: 0 = no regression, 2 = at least one probe regressed, 1 = setup error.
#
#   MINIGO_DIR    — podhmo/minigo checkout with both refs resolvable
#                   (default: sibling clone). git fetch first if needed.
#   SRC_DIR       — probe-target checkouts (default: ./.src, gitignored)
#   RUNS          — interleave rounds per probe (default: 3)
#   THRESHOLD_PCT — |delta|% that counts as a verdict (default: 10)
#   TIMEOUT_SEC   — per probe run (default: 120)
#   --full        — include the grafana-openapi task (fetches the grafana
#                   checkout; no module download needed for the probe)
#   --profile     — after the table, build prof/ per side and write
#                   out/cmp/<probe>.{A,B}.{cpu,allocs}.pprof plus
#                   out/cmp/<probe>.diff.txt (pprof -diff_base) for every
#                   regressed probe
#
# Probes: synthetic dirs under probes/ (no deps) plus light realworld tasks.
# The grafana-swagger-spec / grafana-* tasks are pitch/drift checks, not
# regression probes — only --full pulls grafana-openapi in.
set -u
unset $(compgen -e GIT_) 2>/dev/null || true
ROOT="$(cd "$(dirname "$0")" && pwd)"
MINIGO_DIR="${MINIGO_DIR:-$ROOT/../../minigo}"
SRC_DIR="${SRC_DIR:-$ROOT/.src}"
RUNS="${RUNS:-3}"
THRESHOLD_PCT="${THRESHOLD_PCT:-10}"
TIMEOUT_SEC="${TIMEOUT_SEC:-120}"
OUT="$ROOT/out"
CMP="$OUT/cmp"
FULL=0
PROFILE=0
TIMEOUT=timeout
command -v timeout > /dev/null || TIMEOUT=gtimeout

usage() { sed -n '2,26p' "$0"; exit 1; }
[ $# -ge 2 ] || usage
OLD_REF="$1"; NEW_REF="$2"; shift 2
while [ $# -gt 0 ]; do
	case "$1" in
		--full) FULL=1 ;;
		--profile) PROFILE=1 ;;
		*) usage ;;
	esac
	shift
done
mkdir -p "$OUT" "$SRC_DIR" "$CMP"

now() { perl -MTime::HiRes=time -e 'printf "%.3f", time'; }
elapsed() { perl -e "printf '%.3f', $2 - $1"; }
median() { # median of space-separated numbers on stdin
	perl -e '@x = sort { $a <=> $b } @ARGV; printf "%.3f", $x[int($#x/2)]' "$@"
}

# --- build both sides ---------------------------------------------------
build_side() { # side ref -> prints worktree path; binary at $wt/minigo.cmp
	local side="$1" ref="$2" wt="$CMP/wt-$1"
	rm -rf "$wt"
	git -C "$MINIGO_DIR" worktree add --detach -f "$wt" "$ref" > /dev/null 2>&1 || {
		echo "cannot checkout $ref (git fetch in $MINIGO_DIR first?)" >&2; exit 1; }
	mkdir -p "$wt/bin"
	(cd "$wt" && go build -o "$wt/bin/minigo" ./cmd/minigo) || {
		echo "build failed for $ref" >&2; exit 1; }
	echo "$wt"
}
cleanup() {
	for s in A B; do
		git -C "$MINIGO_DIR" worktree remove --force "$CMP/wt-$s" > /dev/null 2>&1
	done
}
trap cleanup EXIT

echo "building A=$OLD_REF and B=$NEW_REF under $MINIGO_DIR" >&2
WTA=$(build_side A "$OLD_REF"); BIN_A="$WTA/bin/minigo"
WTB=$(build_side B "$NEW_REF"); BIN_B="$WTB/bin/minigo"

# --- light target fetch (git only; no `go mod download`) ----------------
light_fetch() { # target name -> exports nothing; ensures $SRC_DIR/<name>
	local name="$1" url sha dir
	read -r url sha < <(awk -F'\t' -v n="$name" '$1 == n { print $2, $3 }' "$ROOT/targets.tsv")
	[ -n "${sha:-}" ] || { echo "?? target $name: not in targets.tsv" >&2; return 1; }
	dir="$SRC_DIR/$name"
	if [ "$(git -C "$dir" rev-parse HEAD 2>/dev/null)" != "$sha" ]; then
		echo "fetching $name@${sha:0:7} (checkout only) ..." >&2
		rm -rf "$dir" && mkdir -p "$dir" &&
			git -C "$dir" init -q &&
			git -C "$dir" fetch -q --depth 1 "$url" "$sha" &&
			git -C "$dir" checkout -q FETCH_HEAD || return 1
	fi
}

# --- probe runner --------------------------------------------------------
# A probe is "dir|env" lines: dir = minigo-run target, env = optional VAR=val
run_probe() { # bindir probe_dir -> wall seconds (stdout to $CMP/probe.out)
	local bin="$1" dir="$2" t0
	t0=$(now)
	(cd "$dir" && $TIMEOUT "$TIMEOUT_SEC" "$bin" run .) > "$CMP/probe.out" 2>&1
	rc=$?
	elapsed "$t0" "$(now)"
}

PROBES="micro"
[ -d "$ROOT/tasks/clickhouse-settings" ] && PROBES="$PROBES clickhouse-settings"
[ "$FULL" = 1 ] && PROBES="$PROBES grafana-openapi"

# per-probe results live in files: macOS ships bash 3.2, which has no
# associative arrays (declare -A).
RES="$CMP/res"; rm -rf "$RES"; mkdir -p "$RES"
put() { printf '%s' "$3" > "$RES/$1.$2"; }          # put KEY PROBE VALUE
get() { cat "$RES/$1.$2" 2>/dev/null || true; }      # get KEY PROBE
for probe in $PROBES; do
	case "$probe" in
		clickhouse-settings) light_fetch clickhouse-datasource || { echo "skip $probe: target fetch failed" >&2; continue; }
			export TARGET_DIR="$SRC_DIR/clickhouse-datasource"; dir="$ROOT/tasks/clickhouse-settings" ;;
		grafana-openapi) light_fetch grafana || { echo "skip $probe: target fetch failed" >&2; continue; }
			export TARGET_DIR="$SRC_DIR/grafana"; dir="$ROOT/tasks/grafana-openapi" ;;
		*) unset TARGET_DIR || true; dir="$ROOT/probes/$probe" ;;
	esac
	[ -d "$dir" ] || { echo "?? no probe dir $dir" >&2; continue; }

	ta=""; tb=""; out_a=""; out_b=""
	for r in $(seq 1 "$RUNS"); do
		s_a=$(run_probe "$BIN_A" "$dir"); out_a="$out_a$(cat "$CMP/probe.out")"
		s_b=$(run_probe "$BIN_B" "$dir"); out_b="$out_b$(cat "$CMP/probe.out")"
		ta="$ta $s_a"; tb="$tb $s_b"
	done
	med_a=$(median $ta); med_b=$(median $tb)
	delta=$(perl -e "printf '%+.1f', ($med_b - $med_a) * 100 / ($med_a || 1)")
	put med_a "$probe" "$med_a"; put med_b "$probe" "$med_b"; put delta "$probe" "$delta"
	if [ "$out_a" != "$out_b" ]; then
		put verdict "$probe" "OUTPUT-DIFF"
	else
		put verdict "$probe" "$(perl -e 'my ($d, $t) = @ARGV; print $d >= $t ? "REGRESSED" : $d <= -$t ? "IMPROVED" : "same"' -- "$delta" "$THRESHOLD_PCT")"
	fi
done

# --- report -------------------------------------------------------------
echo
printf '%-22s %9s %9s %8s  %s\n' probe "A($OLD_REF)" "B($NEW_REF)" 'delta%' verdict
worst=0
for probe in $PROBES; do
	[ -n "$(get med_a "$probe")" ] || continue
	v=$(get verdict "$probe")
	printf '%-22s %9s %9s %8s  %s\n' "$probe" "$(get med_a "$probe")" "$(get med_b "$probe")" "$(get delta "$probe")" "$v"
	[ "$v" = "REGRESSED" ] && worst=1
	[ "$v" = "OUTPUT-DIFF" ] && worst=1
done

# --- optional profile diff on regressed probes --------------------------
build_prof() { # side worktree -> $OUT/cmp/prof-<side>
	local side="$1" mdir="$2"
	local pdir="$CMP/prof-build-$side"
	rm -rf "$pdir" && mkdir -p "$pdir" && cp "$ROOT/prof/main.go" "$pdir/" &&
		cp "$mdir/go.sum" "$pdir/" &&
		printf 'module realworld/prof\n\ngo 1.26\n\nrequire github.com/podhmo/minigo v0.0.0\n\nreplace github.com/podhmo/minigo => %s\n' "$mdir" > "$pdir/go.mod" &&
		(cd "$pdir" && GOFLAGS=-mod=mod go build -o "$CMP/prof-$side" .)
}
if [ "$PROFILE" = 1 ]; then
	build_prof A "$WTA" && build_prof B "$WTB" || echo "prof build failed; skipping" >&2
	for probe in $PROBES; do
		[ "$(get verdict "$probe")" = "REGRESSED" ] || continue
		case "$probe" in
			clickhouse-settings|grafana-openapi) dir="$ROOT/tasks/$probe" ;;
			*) dir="$ROOT/probes/$probe" ;;
		esac
		(cd "$dir" && "$CMP/prof-A" -dir . -out "$CMP/$probe.A" > /dev/null 2>&1)
		(cd "$dir" && "$CMP/prof-B" -dir . -out "$CMP/$probe.B" > /dev/null 2>&1)
		{
			echo "## cpu diff (A base -> B):"
			go tool pprof -top -nodecount=25 -diff_base "$CMP/$probe.A.cpu.pprof" "$CMP/prof-B" "$CMP/$probe.B.cpu.pprof" 2>/dev/null | sed -n '4,$p'
			echo
			echo "## alloc_space diff:"
			go tool pprof -top -nodecount=25 -sample_index=alloc_space -diff_base "$CMP/$probe.A.allocs.pprof" "$CMP/prof-B" "$CMP/$probe.B.allocs.pprof" 2>/dev/null | sed -n '4,$p'
		} > "$CMP/$probe.diff.txt"
		echo "profile diff -> out/cmp/$probe.diff.txt" >&2
	done
fi

echo
[ "$worst" = 1 ] && { echo "RESULT: regression detected (or output diverged)"; exit 2; }
echo "RESULT: no regression beyond ${THRESHOLD_PCT}%"
exit 0

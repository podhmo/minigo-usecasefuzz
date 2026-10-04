#!/usr/bin/env bash
# run.sh — vet + run the []byte boundary probes.
#
# The go.mod replace resolves podhmo/minigo at ../../minigo (sibling
# clone convention, same as concurrency/host).
set -u
ROOT="$(cd "$(dirname "$0")" && pwd)"

if [ ! -d "$ROOT/../../minigo" ]; then
	echo "cloning podhmo/minigo into $ROOT/../../minigo ..." >&2
	git clone --depth 1 https://github.com/podhmo/minigo "$ROOT/../../minigo" || {
		echo "cannot obtain minigo; clone it next to this repo" >&2
		exit 1
	}
fi

cd "$ROOT"
go vet . && go run . "$@"

#!/usr/bin/env sh
# ShellCheck over every scripts/*.sh, or, given files relative to the root, just
# those (as flai test's shellcheck tier names the ones it selected), so a change
# to the scripts is linted before review (I-0199). It follows what each script
# sources (-x) from beside it (-P SCRIPTDIR), such as env.sh. Findings print in
# ShellCheck's default format, which names the file and line.
# Runs the pinned ShellCheck that scripts/install-tools.sh puts in the main
# checkout's bin/, which env.sh puts first on PATH.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
sc="$(command -v shellcheck || true)"
[ -n "$sc" ] || { echo "shellcheck: not found; run scripts/install-tools.sh (make install-tools) in the main checkout" >&2; exit 1; }
echo "shellcheck $("$sc" --version | sed -n 's/^version: //p') ($sc)"
if [ "$#" -eq 0 ]; then
  set -- scripts/*.sh
fi
exec "$sc" -x -P SCRIPTDIR "$@"

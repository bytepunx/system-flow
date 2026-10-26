#!/usr/bin/env sh
# Reproduces I-0177 in a scratch copy of this tree's env.sh and
# flaiover-install.sh: a flaiover install made from another store, as a bare
# pnpm add makes it from the default store, then flaiover-install.sh
# --if-needed with no TTY, which must pass and leave the install on env.sh's
# store. The packages are local folders, so it needs no network.
set -eu
. "$(dirname "$0")/env.sh"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/scripts" "$tmp/flaiover/dep" "$tmp/dep2"
cp "$ROOT/scripts/env.sh" "$ROOT/scripts/flaiover-install.sh" "$tmp/scripts/"
printf '{"name":"dep","version":"1.0.0"}\n' >"$tmp/flaiover/dep/package.json"
printf '{"name":"dep2","version":"1.0.0"}\n' >"$tmp/dep2/package.json"
printf '{"name":"fixture","private":true,"dependencies":{"dep":"file:./dep"}}\n' >"$tmp/flaiover/package.json"
# CI=true makes pnpm purge without asking, which would hide the abort
unset CI
cd "$tmp/flaiover"
if ! out="$(npm_config_store_dir="$tmp/other-store" pnpm add file:../dep2 </dev/null 2>&1)"; then
  printf '%s\n' "$out"
  echo "flaiover-install-test: the bare pnpm add from another store failed"
  exit 1
fi
# pnpm add can leave the lockfile no newer than .modules.yaml: make it newer, as I-0177 found it
sleep 1
touch pnpm-lock.yaml
if ! out="$("$tmp/scripts/flaiover-install.sh" --if-needed </dev/null 2>&1)"; then
  printf '%s\n' "$out"
  echo "flaiover-install-test: flaiover-install.sh --if-needed failed on an install from another store (I-0177)"
  exit 1
fi
if ! grep -q "\"storeDir\": \"$tmp/.flai-cache/pnpm-store/" node_modules/.modules.yaml; then
  sed -n '/storeDir/p' node_modules/.modules.yaml
  echo "flaiover-install-test: the install is not on env.sh's store $tmp/.flai-cache/pnpm-store"
  exit 1
fi
if [ ! -d node_modules/dep ] || [ ! -d node_modules/dep2 ]; then
  ls node_modules
  echo "flaiover-install-test: the install lacks the lockfile's packages"
  exit 1
fi
echo "flaiover-install-test: an install from another store is made again from env.sh's store"

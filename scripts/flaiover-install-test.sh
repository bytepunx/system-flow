#!/usr/bin/env sh
# Reproduces I-0177 in a scratch copy of this tree's env.sh and
# flaiover-install.sh: a flaiover install made from another store, as a bare
# pnpm add makes it from the default store, then flaiover-install.sh
# --if-needed with no TTY, which must pass and leave the install on env.sh's
# store. Then reproduces I-0259: one install from flaiover's pnpm-workspace.yaml
# must hold both a glibc and a musl optional package, and --if-needed must
# install again once the install configuration is newer than the last install.
# The packages are local folders, so it needs no network.
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

# I-0259: a worktree's node_modules is shared by the host (glibc) and the
# Alpine agent container (musl), so one install must hold both libcs' native
# bindings, and an install made before the setting must be made again. Two
# local optional packages, one per libc, stand in for rolldown's bindings.
lib="$tmp/libc"
mkdir -p "$lib/scripts" "$lib/flaiover" "$lib/gnu" "$lib/musl"
cp "$ROOT/scripts/env.sh" "$ROOT/scripts/flaiover-install.sh" "$lib/scripts/"
printf '{"name":"bind-gnu","version":"1.0.0","libc":["glibc"]}\n' >"$lib/gnu/package.json"
printf '{"name":"bind-musl","version":"1.0.0","libc":["musl"]}\n' >"$lib/musl/package.json"
printf '{"name":"fixture","private":true,"optionalDependencies":{"bind-gnu":"file:../gnu","bind-musl":"file:../musl"}}\n' >"$lib/flaiover/package.json"
cd "$lib/flaiover"
# pnpm's resolution skips an optional package for another libc, as it skips
# rolldown's; a frozen install links local folders whatever their libc, so
# each check resolves afresh
resolve() {
  rm -rf node_modules pnpm-lock.yaml
  if ! out="$(pnpm install --no-frozen-lockfile </dev/null 2>&1)"; then
    printf '%s\n' "$out"
    echo "flaiover-install-test: the libc fixture's install failed"
    exit 1
  fi
}
resolve
if [ -d node_modules/bind-gnu ] && [ -d node_modules/bind-musl ]; then
  echo "flaiover-install-test: both libc variants landed without supportedArchitectures, so the fixture shows nothing"
  exit 1
fi
cp "$ROOT/flaiover/pnpm-workspace.yaml" pnpm-workspace.yaml
resolve
if [ ! -d node_modules/bind-gnu ] || [ ! -d node_modules/bind-musl ]; then
  ls node_modules
  echo "flaiover-install-test: flaiover/pnpm-workspace.yaml does not install both the glibc and the musl variant (I-0259)"
  exit 1
fi
# --if-needed installs again once the lockfile or the install configuration is
# newer than its stamp, and skips otherwise
installs() {
  if ! out="$("$lib/scripts/flaiover-install.sh" --if-needed </dev/null 2>&1)"; then
    printf '%s\n' "$out"
    echo "flaiover-install-test: flaiover-install.sh --if-needed failed in the libc fixture"
    exit 1
  fi
  case "$out" in *"installing flaiover's dependencies"*) return 0 ;; esac
  return 1
}
sleep 1
touch node_modules/.modules.yaml
if installs; then
  echo "flaiover-install-test: --if-needed installed again with nothing newer than the last install"
  exit 1
fi
for f in pnpm-workspace.yaml package.json; do
  sleep 1
  touch "$f"
  installs || { echo "flaiover-install-test: --if-needed skipped the install after $f changed (I-0259)"; exit 1; }
  if installs; then
    echo "flaiover-install-test: --if-needed installed again right after installing for $f"
    exit 1
  fi
done
echo "flaiover-install-test: one install holds both libc variants, and a newer install configuration installs again"

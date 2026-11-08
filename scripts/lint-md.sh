#!/usr/bin/env sh
# Lint every markdown file with the same globs and config as system-flow-check.yml,
# or, given files relative to the root, just those (as flai test's markdown tier
# names the ones it selected).
# In a close-out, which exports the story as CLOSE_OUT_STORY, the run with no
# files leaves out wip/, which the branch holds as the main branch last committed
# it, and lints instead the markdown under wip/ that the story changes against
# the main branch (committed since the merge base, uncommitted, or untracked), so
# a bad line another agent commits to main's wip fails no other story (S-0345,
# I-0117). In a merge queue batch, which exports its stories comma-separated as
# MERGE_QUEUE_STORIES, it does the same for what the batch's checkout changes,
# so a bad line on main's wip fails no batch either (I-0174). The main branch
# is the one the main checkout has, as close-out.sh takes it.
# Runs the markdownlint-cli2 on PATH when it is the pinned version, as
# scripts/install-tools.sh installs it into the main checkout's bin/, which
# env.sh puts on PATH; that saves npx's start, about 0.7 to 0.9 s a run (S-0609,
# cause 19 in design/system/server-performance.md). Otherwise it runs the pinned
# version through npx (Node only; the npm cache lives in .flai-cache via env.sh),
# so it works on a fresh clone and in CI without installing anything.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
# The one pinned version: scripts/install-tools.sh installs this one.
MDLINT_VERSION="${MDLINT_VERSION:-0.20.0}"
# mdlint_installed says whether the markdownlint-cli2 on PATH is MDLINT_VERSION,
# from the package.json beside the script its link resolves to, without
# starting node.
mdlint_installed() {
  mdlint_bin="$(command -v markdownlint-cli2 2>/dev/null)" || return 1
  mdlint_real="$(readlink -f "$mdlint_bin" 2>/dev/null)" || return 1
  mdlint_pkg="$(dirname "$mdlint_real")/package.json"
  [ -f "$mdlint_pkg" ] && grep -q "\"version\": \"$MDLINT_VERSION\"" "$mdlint_pkg"
}
if mdlint_installed; then
  echo "lint-md: markdownlint-cli2 $MDLINT_VERSION, installed"
  mdlint() { markdownlint-cli2 "$@"; }
else
  command -v npx >/dev/null 2>&1 || { echo "lint-md: npx not found; install Node, or markdownlint-cli2 $MDLINT_VERSION with scripts/install-tools.sh" >&2; exit 1; }
  echo "lint-md: markdownlint-cli2 $MDLINT_VERSION, through npx"
  mdlint() { npx --yes "markdownlint-cli2@$MDLINT_VERSION" "$@"; }
fi
if [ "$#" -gt 0 ]; then
  mdlint "$@"
  exit
fi
if [ -z "${CLOSE_OUT_STORY:-}" ] && [ -z "${MERGE_QUEUE_STORIES:-}" ]; then
  mdlint \
    "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**" "!flaiover/build/**" "!**/.svelte-kit/**" "!.flai-cache/**"
  exit
fi

main="$(git -C "$CACHE_ROOT" rev-parse --abbrev-ref HEAD)"
base="$(git merge-base "$main" HEAD)" || { echo "lint-md: no merge base between $main and HEAD" >&2; exit 1; }
# Against the merge base, git diff compares the worktree, so it holds what the
# branch committed and what is not committed yet; deleted files are left out.
changed="$(git diff --name-only --diff-filter=d "$base" -- 'wip/*.md' && git ls-files --others --exclude-standard -- 'wip/*.md')"
if [ -n "${CLOSE_OUT_STORY:-}" ]; then
  echo "lint-md: close-out of $CLOSE_OUT_STORY: wip/ only where the story changes it against $main"
else
  echo "lint-md: merge queue batch of $MERGE_QUEUE_STORIES: wip/ only where the batch changes it against $main"
fi
rc=0
mdlint \
  "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**" "!flaiover/build/**" "!**/.svelte-kit/**" "!.flai-cache/**" "!wip/**" || rc=$?
# markdownlint-cli2's negated globs apply to named files too, so the changed
# wip files have a run of their own.
while IFS= read -r f; do
  if [ -n "$f" ] && [ -f "$f" ]; then
    set -- "$@" "$f"
  fi
done <<EOF
$changed
EOF
if [ "$#" -gt 0 ]; then
  mdlint "$@" || rc=$?
else
  echo "lint-md: the checkout changes no markdown under wip/"
fi
exit "$rc"

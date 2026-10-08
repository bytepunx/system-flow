#!/usr/bin/env sh
# Lint every markdown file with the same globs and config as system-flow-check.yml,
# or, given files relative to the root, just those (as flai test's markdown tier
# names the ones it selected).
# In a close-out, which exports the story as CLOSE_OUT_STORY, the run with no
# files leaves out wip/, which the branch holds as the main branch last committed
# it, and lints instead the markdown under wip/ that the story changes against
# the main branch (committed since the merge base, uncommitted, or untracked), so
# a bad line another agent commits to main's wip fails no other story (S-0345,
# I-0117). The main branch is the one the main checkout has, as close-out.sh
# takes it.
# Runs markdownlint-cli2 through npx (Node only; the npm cache lives in .flai-cache
# via env.sh), so it works on a fresh clone and in CI without pnpm.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
command -v npx >/dev/null 2>&1 || { echo "lint-md: npx not found; install Node" >&2; exit 1; }
MDLINT_VERSION="${MDLINT_VERSION:-0.20.0}"
echo "lint-md: markdownlint-cli2 $MDLINT_VERSION"
if [ "$#" -gt 0 ]; then
  exec npx --yes "markdownlint-cli2@$MDLINT_VERSION" "$@"
fi
if [ -z "${CLOSE_OUT_STORY:-}" ]; then
  exec npx --yes "markdownlint-cli2@$MDLINT_VERSION" \
    "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**" "!flaiover/build/**" "!**/.svelte-kit/**" "!.flai-cache/**"
fi

main="$(git -C "$CACHE_ROOT" rev-parse --abbrev-ref HEAD)"
base="$(git merge-base "$main" HEAD)" || { echo "lint-md: no merge base between $main and HEAD" >&2; exit 1; }
# Against the merge base, git diff compares the worktree, so it holds what the
# branch committed and what is not committed yet; deleted files are left out.
changed="$(git diff --name-only --diff-filter=d "$base" -- 'wip/*.md' && git ls-files --others --exclude-standard -- 'wip/*.md')"
echo "lint-md: close-out of $CLOSE_OUT_STORY: wip/ only where the story changes it against $main"
rc=0
npx --yes "markdownlint-cli2@$MDLINT_VERSION" \
  "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**" "!flaiover/build/**" "!**/.svelte-kit/**" "!.flai-cache/**" "!wip/**" || rc=$?
# markdownlint-cli2's negated globs apply to named files too, so the story's
# wip files have a run of their own.
while IFS= read -r f; do
  if [ -n "$f" ] && [ -f "$f" ]; then
    set -- "$@" "$f"
  fi
done <<EOF
$changed
EOF
if [ "$#" -gt 0 ]; then
  npx --yes "markdownlint-cli2@$MDLINT_VERSION" "$@" || rc=$?
else
  echo "lint-md: the story changes no markdown under wip/"
fi
exit "$rc"

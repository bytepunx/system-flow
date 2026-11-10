#!/usr/bin/env sh
# Smoke tests: end-to-end. Render the template and check it; check this repository
# and its CI pins (scripts/ci-pins.sh, S-0630).
# Under flai verify, in a close-out or a merge queue batch, the repository check
# and the markdown lint run here: TestMonorepoIsClean, in the integration tier,
# skips there, and so does TestRepositoryLintsClean unless the stories change
# mdlint or the markdownlint configuration (S-0548). The gateway checks are not
# here: scripts/gateway-smoke.sh is the gateway tier, which runs only by name,
# flai test gateway, or on the host's schedule, as it spends gateway money.
#
# Under flai verify a step is skipped when the same run already passed it, as
# FLAI_VERIFY_PASSED says: space-separated words, <tier> for a tier that passed
# for the paths, <tier>:all for one that ran its whole command, and step:check
# for the close-out's check step (S-0610, ADR-0266). Unset or empty, as under
# make smoke and flai test smoke, every step runs.
#   - The template test is skipped when the template tier passed.
#   - The repository check is skipped when the close-out's check step passed and
#     the branch changes nothing under flai/. That step runs the running flai's
#     check in process, while check.sh runs the tree's flai (scripts/flai.sh),
#     which, when the branch changes flai/, may check what the running one does
#     not; check.sh is then the only check by the tree's rules under flai verify
#     (ADR-0232).
#   - The markdown lint is skipped when the markdown tier passed whole, or passed
#     for the paths and the branch changes neither a markdownlint configuration
#     file nor lint-md.sh: the markdown it did not lint is the main branch's,
#     unchanged, linted when it last changed. TestRepositoryLintsClean follows
#     the same rule.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
passed=" ${FLAI_VERIFY_PASSED:-} "

# has says whether FLAI_VERIFY_PASSED holds the word $1.
has() {
  case "$passed" in
    *" $1 "*) return 0 ;;
  esac
  return 1
}

# changes says whether the branch changes, against the main branch, a path that
# one of the case patterns it is given matches, committed or not; yes when it
# cannot tell.
changes() {
  if [ -z "${changed_read:-}" ]; then
    changed_read=1
    changed_ok=
    if main="$(git -C "$CACHE_ROOT" rev-parse --abbrev-ref HEAD 2>/dev/null)" &&
      base="$(git merge-base "$main" HEAD 2>/dev/null)" &&
      diff="$(git diff --name-only "$base")" &&
      untracked="$(git ls-files --others --exclude-standard)"; then
      changed="$diff
$untracked"
      changed_ok=1
    fi
  fi
  [ -n "$changed_ok" ] || return 0
  while IFS= read -r f; do
    for pattern in "$@"; do
      # shellcheck disable=SC2254 # each is a case pattern on purpose
      case "$f" in
        $pattern) return 0 ;;
      esac
    done
  done <<EOF
$changed
EOF
  return 1
}

if has template || has template:all; then
  echo "smoke: template render and check: skipped, the template tier passed in this run"
else
  echo "smoke: template render and check"
  "$ROOT/scripts/template-test.sh"
fi

if has step:check && ! changes 'flai/*'; then
  echo "smoke: repository check: skipped, the close-out's check step passed in this run and the branch changes nothing under flai/"
else
  echo "smoke: repository check"
  "$ROOT/scripts/check.sh"
fi

echo "smoke: CI pins: actions by commit SHA, read-only top-level permissions, exact images, one lockfile"
"$ROOT/scripts/ci-pins.sh"

if has markdown:all; then
  echo "smoke: markdown lint: skipped, the markdown tier linted every file in this run"
elif has markdown && ! changes .markdownlint-cli2.jsonc .markdownlint-cli2.yaml .markdownlint.jsonc .markdownlint.json .markdownlint.yaml .markdownlint.yml scripts/lint-md.sh; then
  echo "smoke: markdown lint: skipped, the markdown tier passed in this run and the branch changes neither the markdownlint configuration nor lint-md.sh"
else
  echo "smoke: markdown lint"
  "$ROOT/scripts/lint-md.sh"
fi

echo "smoke: installer and self-upgrade, from a release built from the tree and served on 127.0.0.1"
"$ROOT/scripts/install-test.sh"

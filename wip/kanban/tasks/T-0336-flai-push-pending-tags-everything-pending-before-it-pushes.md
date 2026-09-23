---
id: T-0336
type: task
nature: feature
title: flai push --pending tags everything pending before it pushes
status: done
parent: S-0094
owner: alex
created: 2026-09-23T00:55:34Z
updated: 2026-09-23T00:57:23Z
transitions:
  - to: ready
    at: 2026-09-23T00:55:59Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T00:56:05Z
    by: system-flow
  - to: done
    at: 2026-09-23T00:57:23Z
    by: system-flow
stream: S-0094
tags: []
---

# T-0336 flai push --pending tags everything pending before it pushes

## Work
Since S-0087, `flai accept` computes no release, and nothing else ever called `release.Pending` except `flai release --pending` itself — a fully separate, easy-to-forget command. `flai push --pending` now calls a new shared helper, `computeApplyAndTagPending` (extracted from `flai release --pending`'s own `publishPending`), before it decides what to push: it applies and tags everything `release.Pending` finds accumulated since each component's last tag, then proceeds with the existing push logic, which now also sees the new tags (and the version-bump commit) as part of what is ahead. `--dry-run` only previews the plan (`release.Pending`, read-only); nothing is applied, tagged, or pushed. `flai release --pending` is otherwise unchanged, still available for seeing or forcing a release ahead of a push.

## Done when
All three acceptance criteria: no tag at accept (already true, S-0087, reconfirmed); a real push creates the tag(s) and pushes them with the branch in one operation; a dry run previews without changing anything. Go tests: a dedicated end-to-end test (`TestPushPendingTagsWhatAcceptLeftUnreleased`) proving no tag at accept, no tag on a dry run, a real push creating and pushing the tag, and a resumable no-op on a second run; existing push/release tests updated where they asserted the old "push tags nothing" behavior. `go test ./...`, `go vet`, `gofmt -l .`, `golangci-lint run ./...`, `flai check --strict`, `scripts/lint-md.sh` clean. Verified live, twice, with a locally built `bin/flai` against real scratch git repositories (no remote, then a real bare-repo remote): accept created no tag, a dry run created no tag, a real `flai push --pending --json` created `cli/v1.1.0` locally and pushed it to the remote in the same command, and a second run retagged nothing.

## Notes
Docs updated to match: `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, `docs/operators/index.md`, and the project addition in `design/conventions/git.md`. Found and fixed, in the same pass since they directly contradicted what I was documenting: `docs/users/flai.md` and `docs/operators/index.md` both still described the pre-S-0087 behavior (accept itself computing, tagging, and pushing a release, and flags like `--no-push`/`--no-release` that no longer exist) — stale since S-0087 shipped, not something this story introduced.

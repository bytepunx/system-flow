---
id: T-3441
type: task
nature: feature
title: flai release --pending runs the suites verify_at.publish lists at the main branch's head before it applies or tags, and exits with a code of its own naming the failed tier
status: in-progress
parent: S-0811
owner: alex
created: 2026-11-16T23:48:01Z
updated: 2026-11-18T08:16:43Z
transitions:
  - to: ready
    at: 2026-11-18T08:16:43Z
    by: agent-S-0811
  - to: in-progress
    at: 2026-11-18T08:16:43Z
    by: agent-S-0811
stream: S-0811
tags: [flai]
touches: [flai/cmd/release.go, flai/cmd/release_publish_suites.go, flai/cmd/release_publish_suites_test.go, docs/users/flai-reference.md]
---
# T-3441 flai release --pending runs the suites verify_at.publish lists at the main branch's head before it applies or tags, and exits with a code of its own naming the failed tier

## Work

ADR-0312 §6. In `flai/cmd/release.go`, `publishPending` gains a step before `computeApplyAndTagPending`: once the `push` gate has passed, the remote is in step, and something is pending, it reads `verify_at.publish` from the workflow file (S-0806's parse) and, when it names suites, runs them in the main checkout at its head. Put the reading, the selection, and the run in a new file, `flai/cmd/release_publish_suites.go`.

- Select every tier of those suites, or of every suite for `all`, as `flai test --suite <name> --all` does (S-0805): not by changed paths, and never a `manual` tier.
- Run them with `verify.Run` from the repository's root. Keep the whole output in a log under `.flai-cache/`, and name its path in the answer.
- When one fails, write no version file, changelog, or tag, push nothing, and exit with a new constant beside `exitReleaseNotMet` and `exitReleasePushOff`: the next unused code, 5 unless another is taken. Its message names the suite, the tier, and the log. `--json` gives the suites, each tier's outcome, and the log path.
- With nothing pending, run nothing: the answer stays as it is today. With `--dry-run`, name the suites that would run and run none.
- With no `publish` in the file, as on the baseline, run no tier, and the answer stays as it is today.
- Say all of this in the command's help, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

This task waits for no other: the formula, the dashboard's mapping, and the documents read the exit it defines.

## Done when

- Tests in `flai/cmd/release_publish_suites_test.go`, on a fixture repository with something pending, show each of these:
  - a `verify_at.publish` suite whose tier fails exits with the new code, names the suite and the tier, keeps the output at the path it names, and leaves no new tag, version file, or changelog entry and no push;
  - a passing suite publishes as before;
  - a file with no `publish` runs no tier;
  - `--dry-run` runs no tier.
- `flai test` passes on the changed paths, `TestReferenceIsCurrent` included.

## Notes

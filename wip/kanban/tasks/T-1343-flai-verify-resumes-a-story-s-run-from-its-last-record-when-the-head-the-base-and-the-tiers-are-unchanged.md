---
id: T-1343
type: task
nature: improvement
title: flai verify resumes a story's run from its last record when the head, the base, and the tiers are unchanged
status: in-progress
parent: S-0341
owner: alex
created: 2026-10-08T08:05:06Z
updated: 2026-10-08T08:38:35Z
transitions:
  - to: ready
    at: 2026-10-08T08:38:35Z
    by: agent-S-0341
  - to: in-progress
    at: 2026-10-08T08:38:35Z
    by: agent-S-0341
stream: S-0341
tags: [cli]
touches: [flai/internal/verify/verify.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go]
---
# T-1343 flai verify resumes a story's run from its last record when the head, the base, and the tiers are unchanged

## Work

The core of the story, in `flai/internal/verify`. It waits for nothing: the CLI, the dashboard, the close-out, and the docs all build on the record and the state it defines.

- Add the state `Reused` (`"reused"`) beside `Passed`, `Failed`, and `NotReached` in `verify.go`.
- Give the report a fingerprint of what selects and runs the tiers: the manifest's `tests` and each selected tier's name, command, and directory, hashed. Write it in every record.
- Give a reused step the time of the run it comes from (`reused_from`, a UTC timestamp). A step the last record had as reused keeps that record's `reused_from`, so a chain of resumed runs names the run that actually ran the tier.
- In the story run (`story.go`), read the story's last record with `LastReport`. Resume when its `commit` is the branch head, its `base` is the main commit, and its fingerprint matches; otherwise run in full. A record with no fingerprint, written by an older flai, runs in full.
- When resuming, run `rebase`, `sync`, `narrative`, and `check` as now. Then mark each tier before the first one the record shows as `failed` or `not-reached` as `Reused`, and run the tiers from that one on (`run.go` takes where to start, or the results to keep). A record that stopped before the tiers runs every tier. A record that passed every tier gives a run that reuses every tier.
- Add an option, `Fresh`, that forces a full run.
- Write the resumed run's record whole: every step, reused ones with their state and `reused_from`, and `passed` and `stopped_at` over the whole result.
- `flai test` (`verify.Test`) does not resume; only a story's verify does.

## Done when

- Tests in `story_test.go` and `run_test.go`, with the existing fake git and fake process, cover a resumed run that reruns only from the failed tier and reports the earlier tiers as `reused` with `reused_from`; a full run on each of a changed head, a changed base, a changed tier selection, a changed manifest `tests`, a changed tier command, and a record without a fingerprint; `Fresh`; and the record a resumed run writes holding every tier.
- `flai test flai/internal/verify` passes.

## Notes

Test helpers: `fakeProc` and `step` in `run_test.go`, `scriptGit` and `newVerifyFixture` in `story_test.go`.

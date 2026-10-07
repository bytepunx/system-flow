---
id: T-1145
type: task
nature: improvement
title: A close-out records no wip.overlap in an issue, and its scoped check lists as notes only the overlaps that name the story
status: done
parent: S-0279
owner: alex
created: 2026-10-07T01:20:27Z
updated: 2026-10-07T09:41:12Z
transitions:
  - to: ready
    at: 2026-10-07T09:35:16Z
    by: agent-S-0279
  - to: in-progress
    at: 2026-10-07T09:35:16Z
    by: agent-S-0279
  - to: done
    at: 2026-10-07T09:41:12Z
    by: agent-S-0279
stream: S-0279
tags: [flai]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go]
after: [T-1143]
usage:
  source: log
  seconds: 356
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 37
      output: 10649
      cache_read: 1695432
      cache_write: 64495
      cost: 0.9686
---
# T-1145 A close-out records no wip.overlap in an issue, and its scoped check lists as notes only the overlaps that name the story

## Work

Change how `flai check --story S-nnnn --record-issues` treats `wip.overlap`, as T-1141's ADR decides. This is the change that stops I-0076 from recurring. It waits for T-1143, which carries the two story IDs on each `wip.overlap` finding. This task scopes the findings by those IDs, and T-1143 waits for T-1141's ADR.

The recommended rule:

- In `ScopeToStory` (`flai/internal/check/scope.go`), a `wip.overlap` that does not name the story is dropped from the scoped result. One that names it stays a note outside the story, as today, so the agent still sees it. Tell the two apart by the story IDs T-1143 puts on the finding, never by parsing its message.
- `recordOutside` in `flai/cmd/check.go` records no `wip.overlap`. The pull hold, the `overlapped` change, and the trial merge at `flai stream sync` already report an overlap, so an issue adds nothing.
- Update the comment on `ScopeToStory`, and update `flai check`'s help (`Long`), which says every `wip.overlap` is outside the story.

## Done when

- A test in `flai/cmd/check_test.go` reproduces I-0076 and passes. Two other stories in progress overlap, and a scoped check with `--record-issues` for a third story opens or bumps no issue. Its output carries no `wip.overlap`.
- A test passes where the closing story is in the overlapping pair. The finding is printed as a note outside the story, and no issue is recorded.
- `TestScopedCheckFailsOnlyOnTheStorysOwnFindings` and `TestCheckRecordIssuesOpensThenBumpsOncePerStory` pass, updated where the rule changes them. Another rule's findings outside the story are still recorded.
- `flai test flai/internal/check flai/cmd` passes.

## Notes

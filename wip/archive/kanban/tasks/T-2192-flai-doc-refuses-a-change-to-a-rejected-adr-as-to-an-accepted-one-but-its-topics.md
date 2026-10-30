---
id: T-2192
type: task
nature: remediation
title: flai doc refuses a change to a rejected ADR as to an accepted one, but its topics
status: done
parent: S-0512
owner: alex
created: 2026-10-28T08:28:01Z
updated: 2026-10-30T14:15:23Z
transitions:
  - to: ready
    at: 2026-10-30T14:13:43Z
    by: agent-S-0512
  - to: in-progress
    at: 2026-10-30T14:13:44Z
    by: agent-S-0512
  - to: done
    at: 2026-10-30T14:15:23Z
    by: agent-S-0512
stream: S-0512
tags: [cli]
touches: [flai/internal/docedit/docedit.go, flai/internal/docedit/docedit_test.go]
after: [T-2190]
usage:
  source: log
  seconds: 99
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 4298
      cache_read: 964544
      cache_write: 42375
      cost: 0.5651
---
# T-2192 flai doc refuses a change to a rejected ADR as to an accepted one, but its topics

## Work

In `flai/internal/docedit/docedit.go`, `acceptedStatus` matches `accepted`, `superseded`, and `deprecated`, and `acceptedADR` makes such an ADR read-only for `flai doc save`, `Show`, and the dashboard's editor, but for its `topics` (`onlyTopics`). Make a rejected ADR immutable the same way: build the match from the one list T-2190 adds, every status but `proposed`, rather than adding a fourth word to the regexp, so that the two cannot drift again. Update the comment on `acceptedADR`: a proposed ADR is a draft; once accepted, superseded, deprecated, or rejected, only `superseded_by` and `topics` may change.

In `docedit_test.go`, add a test that a save changing the body of an ADR with status `rejected` is refused as an accepted one's is, and that a save changing only its `topics` goes through.

This task waits for T-2190, which adds the list it reads.

## Done when

- `flai doc save` refuses a body change to an ADR of status `rejected` with the same refusal as for an accepted one, and accepts a change to its `topics` alone.
- The match is built from the shared list, not a second copy of the statuses.
- The Go tests of `flai/internal/docedit` pass, and `flai test` on the changed paths passes.

## Notes

Drafted by planner-S-0512. `docedit` imports `workitem` and `topics` today; importing the package T-2190 chooses for the list must not make a cycle, which is why that task is asked to put it where `docedit` can reach it.

---
id: S-0290
type: story
nature: remediation
title: TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
status: done
owner: alex
created: 2026-10-06T09:56:52Z
updated: 2026-10-08T09:08:56Z
transitions:
  - to: ready
    at: 2026-10-08T08:40:35Z
    by: alex
  - to: in-progress
    at: 2026-10-08T09:00:23Z
    by: agent-S-0290
  - to: review
    at: 2026-10-08T09:08:12Z
    by: agent-S-0290
  - to: done
    at: 2026-10-08T09:08:56Z
    by: orchestrator
tags: [flai, tests]
touches: [flai/internal/workitem/workitem_test.go, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 479
  turns:
    - day: 2026-10-08
      ceremony: 2
      test_runs: 2
      hand_edits: 2
      work: 35
  models:
    - model: claude-opus-5-5
      input: 84
      output: 13500
      cache_read: 5303853
      cache_write: 166472
      cost: 2.6629
  strategic:
    - kind: orchestrator
      seconds: 235
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 788
          cache_read: 15126415
          cache_write: 14626
          cost: 3.7305
        - model: claude-sonnet-5-5
          input: 4
          output: 36
          cache_read: 29103
          cache_write: 17519
          cost: 0.0402
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-06T09:56:52Z
  value: 268
  by: planner-S-0290
  at: 2026-10-08T08:42:37Z
forecast:
  duration: 30m
  delivery: 2026-10-08T09:31:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 5, behind S-0232, S-0288, S-0321, S-0322 and S-0341."
  by: flai
  at: 2026-10-08T08:59:49Z
finalized:
  by: alex
  at: 2026-10-07T02:19:37Z
---
# S-0290 TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run

## Goal

This story remediates [I-0079](../../../design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md), "TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0079 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0079 is closed with `flai issue close I-0079 --reason` saying what fixed it

## Tasks
- T-1367 TestRoundTripRepositoryItems round-trips each item from one read of its file and skips a file gone since the list
- T-1368 Close I-0079 with the one-read round trip as its fix

## Notes

Cost of delay inputs set by flai from I-0079. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-05T05:55:58Z, 1.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the nine instances. In a worktree, `Open` puts `wip/` in the main checkout (`Repo.WipDir`). The test parses each item through `List`, then reads the file a second time with `os.ReadFile`. Every instance is an item rewritten between those two reads. `Repo.Save` writes atomically, so one read is always a whole file. The fix parses and compares from that one read, and skips a file gone since the list (T-1367).

Touches, all files, none a folder:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/workitem_test.go` | layout: the test I-0079 names is defined here |
| `design/issues/I-0079-…md` | design: criterion 2 closes the issue |
| `design/issues/summary.md` | co-change: `flai issue close` rewrites it |

`flai touches suggest` from the test file listed `work-hierarchy.md`, `docs/users/flai.md`, and `item.go` as co-changes. They are left out because the fix changes only a test, not the item format or a command. No design document names the test. Both `design/issues` paths are in `claims.shared`. S-0321, in progress, also touches the I-0079 file.

Forecast: `flai forecast` gave 14m (size 5). It is raised to 30m because the close-out runs the whole integration tier, which is where the race showed, and the story adds a reproduction test. flai placed the delivery from the 30m and the pull order.

Cost of delay: `flai cod` gives 15 USD a week from the input of 6m a cycle. That input counted one occurrence when the story was made. The value is set to 268 USD a week, from I-0079 as it stands:

```text
9 occurrences x 5m = 45m between 2026-10-05T05:55:58Z and 2026-10-08T04:24:20Z (70.5h)
45m x 168h / 70.5h = 107m a week, about 1.79h x 150 USD = 268 USD a week
```

The input is left as it is because it is the operator's. TH-0377 recommends raising it to 1h47m, which makes `flai cod` agree.

### Accepted by the orchestrator

- Verified: 75462ec068919a4d4c2acf158c4496b1bdebcba9
- At: 2026-10-08T09:08:56Z

Verdict: accept. flai verify passed every step at the branch head 75462ec0, and the verifier matched both criteria to the diff. A minor note: the rewritten-file test's comment describes the old behaviour loosely.
- 1: flai/internal/workitem/workitem_test.go
- 2: design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/summary.md

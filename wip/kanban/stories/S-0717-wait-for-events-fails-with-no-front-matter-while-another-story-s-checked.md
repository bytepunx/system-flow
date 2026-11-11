---
id: S-0717
type: story
nature: remediation
title: wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file
status: ready
owner: alex
created: 2026-11-11T01:43:38Z
updated: 2026-11-11T22:31:13Z
transitions:
  - to: ready
    at: 2026-11-11T22:31:08Z
    by: orchestrator
tags: [flai]
touches: [flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder_test.go, design/issues/I-0345-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md, design/issues/summary.md]
after: [S-0716]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 33
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 3250
          cache_read: 1510737
          cache_write: 18132
          cost: 0.5565
    - kind: orchestrator
      seconds: 38
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 117
          cache_read: 424497
          cache_write: 93339
          cost: 0.1449
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-11T01:43:38Z
  value: 12.5
  by: planner-S-0717
  at: 2026-11-11T09:08:16Z
forecast:
  duration: 1h
  delivery: 2026-11-11T23:54:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, behind S-0720."
  by: flai
  at: 2026-11-11T22:31:13Z
finalized:
  by: alex
  at: 2026-11-11T09:01:50Z
---
# S-0717 wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file

## Goal

This story remediates [I-0345](../../../design/issues/I-0345-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md), "wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0345 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0345 is closed with `flai issue close I-0345 --reason` saying what fixed it

## Tasks
- T-3120 Find what wrote PLACEHOLDER into T-3079's file and correct I-0345's account of it
- T-3121 wait_for_events, in one project and across a folder of projects, answers past a task file that holds only PLACEHOLDER, with tests
- T-3122 Close I-0345 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0345. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T11:10:23Z, 0.6 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from the instance. I-0345 and I-0344 record one event: from about 11:04:00Z to 11:04:25Z on 2026-11-10, T-3079's file held only `PLACEHOLDER`. Every `Repo.List` read of the project failed with "no front matter" while it did. The issue blames a checked creation (ADR-0272), but the code does not bear that out. `itemnew.Create` writes a whole item through `atomicfile.WriteFile`, and no Go code under `flai/` holds the string. T-3079 was created at 11:03:52Z, before the placeholder appeared. So the file held a whole item, then the placeholder, then the item again, and the writer is not yet known.

The failure is the loader's: `Repo.List` returns the first file that does not parse as an error for the whole call. S-0716 fixes that for I-0344: `List` skips such a file and keeps its path and error, and `flai check` reports it as one finding. This story builds on that fix and adds what S-0716 leaves out:

- T-3120 finds the writer and corrects the issue's account of it.
- T-3121 pins `wait_for_events` in one project and across a folder of projects, whose `catchUp` fails every project on one project's error.
- T-3122 closes I-0345.

`after` S-0716: this story reads S-0716's change to `Repo.List` in `flai/internal/workitem/store.go` and `repo.go`. T-3121's tests pass only once it is on main, and T-3122's reason cites it.

Touches:

- Layout, from the explorer's read of the code: `flai/internal/mcpserver/cursor.go` (`server.catchUp`) and `folder.go` (`folderWaits.catchUp`, before the hold and inside its callback). T-3121 changes them only if a test still fails once S-0716 is on main.
- Co-change, from `flai touches suggest` on those two files: `flai/internal/mcpserver/server_test.go` (49%) and `folder_test.go` (35%), where `wait_for_events` is tested today. I left out `server.go` (68%). Nothing in the plan changes it, and `flai task done` widens the touches if a task does.
- The criteria: `design/issues/I-0345-…md` and `design/issues/summary.md`, which T-3120 and `flai issue close` write.

No folder touch is kept. `inbox_test.go` and `stream_test.go` are S-0716's (T-3109) and are not repeated here.

Forecast. With the six touches in place, flai gives 18m: a median 135 s per unit over 59 done remediation stories in the medium band, times a size of 8. I raised it to 1h for three reasons. T-3120 searches run logs and the serve journal for a writer that no code names. T-3121 needs a fixture for a folder of two projects. The close-out runs the integration and smoke tiers, because the story changes `flai/`. Delivery is flai's 23:25Z plus the 42 minutes added: 2026-11-12T00:07:00Z.

Cost of delay. The value of 12.50 USD a week stands, as flai worked it out from the operator's input: 5m lost per 168h cycle at 150 USD an hour. The issue gives no reason for more. The instance cost 5m once, and S-0716 removes the failure itself.

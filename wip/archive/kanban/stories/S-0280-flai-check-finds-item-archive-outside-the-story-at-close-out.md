---
id: S-0280
type: story
nature: improvement
title: "flai check finds `item.archive` outside the story at close-out"
status: done
owner: alex
created: 2026-10-05T07:09:05Z
updated: 2026-10-07T23:45:02Z
transitions:
  - to: ready
    at: 2026-10-07T23:21:00Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T23:21:22Z
    by: agent-S-0280
  - to: review
    at: 2026-10-07T23:43:12Z
    by: agent-S-0280
  - to: done
    at: 2026-10-07T23:45:02Z
    by: orchestrator
tags: [flai, template]
topics: [cli, conventions, template]
touches: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1334
  turns:
    - day: 2026-10-07
      test_runs: 2
      hand_edits: 3
      work: 44
  models:
    - model: claude-opus-5-5
      input: 100
      output: 26138
      cache_read: 7400592
      cache_write: 221333
      cost: 3.7739
  strategic:
    - kind: orchestrator
      seconds: 208
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 403
          cache_read: 5378417
          cache_write: 124902
          cost: 1.3565
        - model: claude-sonnet-5-5
          input: 10
          output: 59
          cache_read: 103786
          cache_write: 44423
          cost: 0.1278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 45m
    by: orchestrator
    at: 2026-10-07T22:36:31Z
  value: 112.5
  by: planner-S-0280
  at: 2026-10-07T22:37:10Z
forecast:
  duration: 25m
  delivery: 2026-10-08T04:20:00Z
  basis: "Its own forecast of 25m; 9th in the pull order with an in-progress limit of 3, behind S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239 and S-0241."
  by: flai
  at: 2026-10-07T23:20:33Z
finalized:
  by: orchestrator
  at: 2026-10-07T22:37:17Z
---
# S-0280 flai check finds `item.archive` outside the story at close-out

## Goal

This story remediates [I-0078](../../../design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md), "flai check finds `item.archive` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0078 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0078 is closed with `flai issue close I-0078 --reason` saying what fixed it

## Tasks
- T-1164 An ADR refining ADR-0085 and ADR-0115 records the remedy for I-0078, proposed from its instances
- T-1165 A close-out records no item.archive in an issue, and its scoped check leaves out one that does not name the story
- T-1166 I-0078 is closed with flai issue close, saying that a close-out records no item.archive
- T-1167 The convention, the design, and the users' guide say that a close-out records no item.archive

## Notes

### Planning

Planned by planner-S-0280 on 2026-10-07. The plan thread on S-0280 lists the tasks, their layers, and the assumptions.

The proposed remedy comes from I-0078's 38 instances. T-1164 records it in an ADR before any code is built.

- **One cause.** All 38 instances name S-0250. The operator cancelled it from backlog on 2026-10-05, and it stayed in `wip/kanban/stories` until commit 1f7b7296 archived it by hand on 2026-10-07. No instance came after that.
- **Never the story's.** `item.archive` warns on a done or cancelled epic or story that is not archived. The closing story is in progress, so the finding never names it. No story branch makes it, and no story's agent can fix it in its worktree.
- **Recorded every time.** `ScopeToStory` marks it outside, and `recordOutside` records every outside rule but `wip.overlap`.
- **Remedy.** As ADR-0115 did for `wip.overlap`: a close-out records no `item.archive` and leaves out one that does not name the story. `flai check` in the main checkout still warns.
- **Not taken.** Archiving on cancel. It would reverse ADR-0055's move back from cancelled and ADR-0028's narrative left in place.

Layers:

1. T-1164, the ADR.
2. T-1165, the scoped check and the recording, with tests that reproduce I-0078.
3. T-1166, closing I-0078; T-1167, the convention, the design, and the guides. They share no path and run together.

Touches:

- **Declared:** none. The story declared no touches.
- **Layout:**
  - `flai/internal/check/scope.go` and `scope_test.go`: `ScopeToStory`.
  - `flai/cmd/check.go` and `check_test.go`: `recordOutside` and the help.
- **Co-change:** `docs/users/flai-reference.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` change with `flai/cmd/check.go` in 31 to 38% of its commits. The reference is generated from the changed help.
- **Design:**
  - The two copies of `work-management.md`: the close-out rule names `wip.overlap`.
  - `template/CHANGELOG.md`: the template convention changes.
  - `design/system/continuous-improvement.md`: it describes what a close-out records (ADR-0085).
- **Criterion 2:** the I-0078 file and `design/issues/summary.md`, which `flai issue close` writes.
- **Folder touch kept:** `design/adrs`. T-1164 adds an ADR whose number and slug `flai adr new` picks, so no task can name the file yet. It lies inside `claims.shared`, so it holds no story.
- **Not taken from `flai touches suggest`:** `flai/internal/check/check.go` and `check_test.go`, since the rule itself does not change. Also `design/system/workflow.md`, `docs/operators/settings.md`, and the other co-changes the goal names none of. T-1165 widens its touches if `check.go` must change.
- **Overlap with S-0333**, in progress: it claims both copies of `work-management.md`, which T-1167 changes. S-0280 is held while S-0333 is in progress. Whichever goes second rebases onto the other.

Forecast 25m.

- `flai forecast` gave 21m: 78 s per unit of size over 43 done large-band improvement stories on claude-opus-5-5, times size 16 (2 criteria, 14 touches).
- Raised by 4m. The story has three serial layers, each with its own commit, sync, and test cycle. S-0279, the same remedy for I-0076 with one more layer, took 26m of agent time.
- Delivery is flai's, replanned from the pull order as stories are accepted.

Cost of delay 112.50 USD a week, as `flai cod` gives it.

- The input is `time_lost_per_cycle: 45m`, set by the orchestrator on the cost thread on S-0280, as recommended there.
- I-0078 has 38 instances in about two days while S-0250 lingered. At about 1m of agent time and one issue-bump commit each, averaged over weeks with and without a lingering cancelled item, that is 45m per 168h cycle, at 150 USD an hour.
- The figure stands as computed.

### Accepted by the orchestrator

- Verified: 04ddae77a3de6dca391e8f91e7bde2f00cc21dd8
- At: 2026-10-07T23:45:02Z

Verdict: meets all criteria (verifier at 04ddae77a3de6dca391e8f91e7bde2f00cc21dd8; flai verify passed every step at that commit). ADR-0122 records the remedy; the I-0079 bump on the branch is the close-out's own record, per continuous-improvement.md.

- 1: flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/adrs/0122-a-close-out-records-no-item-archive-and-a-check-scoped-to-a-story-leaves-out.md
- 2: design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/summary.md

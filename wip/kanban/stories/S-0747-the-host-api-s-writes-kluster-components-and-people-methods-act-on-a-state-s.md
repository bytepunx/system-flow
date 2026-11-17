---
id: S-0747
type: story
nature: feature
title: The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name
status: review
parent: E-0034
owner: alex
created: 2026-11-12T16:00:29Z
updated: 2026-11-15T02:43:17Z
transitions:
  - to: ready
    at: 2026-11-15T02:00:13Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-15T02:00:20Z
    by: agent-S-0747
  - to: review
    at: 2026-11-15T02:43:17Z
    by: agent-S-0747
tags: [cli]
topics: [cli, dashboard]
touches: [flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/writes_describe.go, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/people.go, flai/internal/hostapi/people_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/flow.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/hostapi_test.go, design/system/flaiover-dashboard.md]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2590
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 3
      hand_edits: 1
      work: 56
  models:
    - model: claude-opus-5-5
      input: 344
      output: 96135
      cache_read: 22331046
      cache_write: 668873
      cost: 10.662
  strategic:
    - kind: orchestrator
      seconds: 229
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 63
          output: 409
          cache_read: 2416794
          cache_write: 248611
          cost: 0.8039
        - model: claude-sonnet-5-5
          input: 14
          output: 99
          cache_read: 143024
          cache_write: 77032
          cost: 0.173
cost_of_delay:
  value: 67.7
  by: system-flow
  at: 2026-11-14T17:44:08Z
forecast:
  duration: 38m
  delivery: 2026-11-15T03:09:00Z
  basis: "Its own forecast of 38m; 4th in the pull order with an in-progress limit of 3, behind S-0742, S-0743, S-0741, S-0789 and S-0241."
  by: flai
  at: 2026-11-15T01:56:42Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:55Z
---
# S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name

## Goal

What the dashboard asks of flai on the host follows the workflow. A move from the board is refused or allowed by who may make its transition, an acceptance from the board makes the acceptance transition, and the descriptions of each write, the Deployments page's namespaces, the components' worktrees, and the people's open work read a state's kind, so the dashboard's writes work on a workflow they have never seen.

## Acceptance criteria

- [x] The host API's files in `flai/internal/hostapi` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [x] `item.move` from the dashboard refuses a transition the operator may not make there, naming the transition and who may make it, and makes the acceptance transition through acceptance as today
- [x] `writes_describe.go` describes each move in the workflow's names and the reason it asks for, from the transition, not from a list of states
- [x] Tests run the move, the acceptance, the descriptions, and the kluster and people reads against a workflow with a state renamed and one added between `in-progress` and `review`
- [x] With the baseline, every existing `hostapi` test and flaiover's tests of the built flai pass unchanged

## Tasks

- T-3228 The dashboard's item.move refuses a transition the operator may not make there, and an acceptance from the board makes the acceptance transition
- T-3229 The host API's kluster, components, people, import, board, and settings reads tell a story's place by kind
- T-3230 writes_describe.go describes each move from its transition, in the workflow's names and with the reason it asks for

## Notes

- `hostapi` holds about 73 references in 12 files, 31 of them in `writes_describe.go`. The method that gives the dashboard the workflow's columns is the board's story's; this one changes no answer's shape.

### Planning

- Touches, all files:
  - Layout: the twelve files of `flai/internal/hostapi` that name a state, with the tests of the writes and of the kluster and people reads.
  - Co-change: `flai/internal/hostapi/writes.go` (62% in `flai touches suggest`), where the write methods are registered.
  - `hostapi.go` is touched by S-0753 too, which adds the columns to `board.get`; the trial merge judges that.
  - Left out from `flai touches suggest`: `flaiover/src/lib/server/agent.ts` (37%) and `cmd/serve_actions.go` (28%). No answer's shape changes, so flaiover needs no change here; the `flaiover-flai` tier runs flaiover's tests of the built flai since `hostapi` changes.
- Tag `dashboard` dropped from the first draft: every path is flai's. The topic stays, since the story is about what the dashboard reaches.
- `after`: S-0737. The writes read who may make each transition, through S-0737's mover check, and the kinds from the workflow S-0737 gives; the tests use its test hook.
- Tasks: T-3228 (moves and acceptance) and T-3229 (the reads) together; T-3230 (descriptions) after T-3228, since both change `writes_test.go` and the descriptions name the moves T-3228 allows.
- Forecast 38m, `flai forecast`'s rate for size 21, up from its 37m for the first draft's 20.
- Cost of delay: 67.70 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 38m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.

### As built

- Criterion 1: no file of `flai/internal/hostapi` compares a status with a state's name. Three words that look like one stay, because no answer's shape changes: the designer inbox's entry kind and count key `review`, and the Deployments answer's JSON field `review`.
- Two checks outside `hostapi` still read names, so a workflow that renames states is not followed end to end yet: `flai move` asks a reason by the names cancelled and review to in-progress, not by the `reason-given` rule (`flai/internal/workitem/rules.go`, S-0738's), and `kluster.DeployBy` deploys only a story in `in-progress` or `review` (`flai/internal/kluster/deploy.go`, S-0743's). The host API's own checks read the kind.

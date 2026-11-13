---
id: S-0747
type: story
nature: feature
title: The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:29Z
updated: 2026-11-13T17:24:09Z
transitions: []
tags: [cli]
topics: [cli, dashboard]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_run.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/writes_describe.go, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_forward.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/people.go, flai/internal/hostapi/people_test.go, flai/internal/hostapi/import.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/settings.go]
after: [S-0737]
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
    - kind: orchestrator
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 108313
          cache_write: 18018
          cost: 0.0361
draft: true
forecast:
  duration: 38m
  delivery: 2026-11-13T21:27:00Z
  basis: "Its own forecast of 38m; 13th in the pull order with an in-progress limit of 3, behind S-0759, S-0760, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744 and S-0746."
  by: flai
  at: 2026-11-13T17:24:09Z
---
# S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name

## Goal

What the dashboard asks of flai on the host follows the workflow. A move from the board is refused or allowed by who may make its transition, an acceptance from the board makes the acceptance transition, and the descriptions of each write, the Deployments page's namespaces, the components' worktrees, and the people's open work read a state's kind, so the dashboard's writes work on a workflow they have never seen.

## Acceptance criteria

- [ ] The host API's files in `flai/internal/hostapi` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [ ] `item.move` from the dashboard refuses a transition the operator may not make there, naming the transition and who may make it, and makes the acceptance transition through acceptance as today
- [ ] `writes_describe.go` describes each move in the workflow's names and the reason it asks for, from the transition, not from a list of states
- [ ] Tests run the move, the acceptance, the descriptions, and the kluster and people reads against a workflow with a state renamed and one added between `in-progress` and `review`
- [ ] With the baseline, every existing `hostapi` test and flaiover's tests of the built flai pass unchanged

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
- Cost of delay: not set yet, pending TH-1097, as on S-0736.

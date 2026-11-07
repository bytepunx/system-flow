---
id: S-0582
type: story
nature: improvement
title: "flai check finds `agent.image` outside the story at close-out"
status: done
owner: alex
created: 2026-11-06T23:25:47Z
updated: 2026-11-07T00:49:02Z
transitions:
  - to: ready
    at: 2026-11-07T00:33:10Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T00:33:15Z
    by: agent-S-0582
  - to: review
    at: 2026-11-07T00:48:12Z
    by: agent-S-0582
  - to: done
    at: 2026-11-07T00:49:02Z
    by: alex
tags: [flai, check, container]
topics: [check, close-out, issues, container, cli, conventions, template]
touches: [flai/internal/check/image.go, flai/internal/check/image_test.go, flai/internal/check/check.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/adrs, design/adrs/README.md, docs/users/flai.md, design/system/flai-cli.md, design/system/workflow.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0264-flai-check-finds-agent-image-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 908
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 5
      hand_edits: 1
      work: 32
  models:
    - model: claude-opus-5-5
      input: 138
      output: 33464
      cache_read: 15837054
      cache_write: 316497
      cost: 6.2438
  strategic:
    - kind: planner
      seconds: 249
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 70
          output: 20123
          cache_read: 3055772
          cache_write: 101398
          cost: 1.8251
    - kind: orchestrator
      seconds: 80
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 217
          cache_read: 1021840
          cache_write: 146325
          cost: 0.2814
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: orchestrator
    at: 2026-11-07T00:29:18Z
  value: 25
  by: planner-S-0582
  at: 2026-11-07T00:31:32Z
forecast:
  duration: 45m
  delivery: 2026-11-07T08:20:00Z
  basis: "flai forecast's 31m (92 s per unit of size times size 20, 119 done large improvement stories), raised to 45m because S-0446, the same leave-out for components.examples with an ADR and the same files, took 47m from in-progress to review; delivery moved by the 14m added."
  by: planner-S-0582
  at: 2026-11-07T00:31:32Z
finalized:
  by: alex
  at: 2026-11-07T00:27:54Z
---
# S-0582 flai check finds `agent.image` outside the story at close-out

## Goal

This story remediates [I-0264](../../../design/issues/I-0264-flai-check-finds-agent-image-outside-the-story-at-close-out.md), "flai check finds `agent.image` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0264 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0264 is closed with `flai issue close I-0264 --reason` saying what fixed it

## Tasks
- T-2493 An ADR decides that a story's scoped check leaves out the agent.image notes only the host clears
- T-2494 A story's scoped check leaves out the agent.image notes only the host clears, unless the story changes the manifest or the layer
- T-2495 The docs and the close-out rule say a story's scoped check leaves out the agent.image notes only the host clears
- T-2496 I-0264 is closed with the reason that the scoped check now leaves the agent.image notes out

## Notes

### Planning

Proposed fix, for the story's agent to settle in an ADR before building it. `check.AgentImage` (`flai/internal/check/image.go`) gives three advisory `agent.image` warnings that only the host clears: docker not on `PATH`, docker could not say, and the layer's image not built. Each sits on `system-flow.yaml`, so a close-out of a story that does not change the manifest finds it outside the story. `--record-issues` then records it, and inside an agent container, where docker is never on `PATH`, every close-out stops on the uncommitted issue file (ADR-0253). The scoped check should leave these three out, as ADR-0171 does for the `components.examples` note, unless the story changes `system-flow.yaml` or the layer's Dockerfile. The warning for a layer that is not a file stays as it is, since a story can fix it. Precedent: S-0446 (ADR-0171, I-0161) did the same for `components.examples`, with the same files.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/check/image.go`, `image_test.go` | design: I-0264's finding is written here |
| `flai/internal/check/check.go` | layout: `Finding` and its `unchecked` mark (ADR-0171) |
| `flai/internal/check/scope.go`, `scope_test.go` | layout: `ScopeToStories` leaves the notes out |
| `flai/cmd/check.go`, `check_test.go` | layout: the `--story` help and `recordOutside` name what is left out |
| `docs/users/flai-reference.md` | co-change: regenerated from the `--story` help, as in S-0446 |
| `design/adrs`, `design/adrs/README.md` | design: the new ADR and its index |
| `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/workflow.md` | co-change: S-0446 updated the same scope lists |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | co-change: the close-out rule lists what the scoped check leaves out |
| `design/issues/I-0264-…md`, `design/issues/summary.md` | declared by criterion 2 |

`flai touches suggest` listed nothing, since the story declared no touches. Run again with the touches above, it listed only co-changes at 10% or less, none of which the fix reaches. The folder touch `design/adrs` stays because the ADR's number and slug are not known until it is written; T-2493's touch names the predicted file.

Figures:

- Forecast: 45m, delivery 2026-11-07T08:20:00Z. `flai forecast` gave 31m (92 s per unit of size times size 20, over 119 done large improvement stories). Raised because S-0446, the same leave-out for `components.examples` with an ADR and the same files, took 47m from in-progress to review.
- Cost of delay: 25 USD a week, as `flai cod` works it out from the orchestrator's input on TH-0833, `time_lost_per_cycle: 10m`, at 150 USD an hour and one cycle a week. It stands: inside an agent container, every close-out loses about that much to the stopped commit until this lands.

Layers: T-2493 first; T-2494 and T-2495 together, after T-2493, with no file in common; T-2496 last, after both. No story `after`: S-0582 reads nothing another open story writes.

---
id: S-0699
type: story
nature: improvement
title: A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context
status: done
owner: alex
created: 2026-11-10T08:55:29Z
updated: 2026-11-11T18:16:27Z
transitions:
  - to: ready
    at: 2026-11-11T08:59:53Z
    by: alex
  - to: in-progress
    at: 2026-11-11T15:03:49Z
    by: agent-S-0699
  - to: review
    at: 2026-11-11T15:31:51Z
    by: agent-S-0699
  - to: done
    at: 2026-11-11T18:16:27Z
    by: alex
tags: [flai, cli, go]
touches: [flai/internal/topics/story.go, flai/internal/topics/story_test.go, flai/internal/context/context_test.go, flai/internal/context/items.go, flai/internal/context/context.go, flai/internal/context/design_test.go, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md, design/issues/summary.md, flai/internal/context/parts.go, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1698
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 1
      work: 53
  models:
    - model: claude-opus-5-5
      input: 210
      output: 45780
      cache_read: 9506369
      cache_write: 403966
      cost: 5.2947
  strategic:
    - kind: planner
      seconds: 24
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 3661
          cache_read: 1151724
          cache_write: 37652
          cost: 0.5206
    - kind: orchestrator
      seconds: 88
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 197
          cache_read: 947408
          cache_write: 91574
          cost: 0.2913
cost_of_delay:
  inputs:
    time_lost_per_cycle: 18m
    by: flai
    at: 2026-11-10T08:55:29Z
  value: 45
  by: planner-S-0699
  at: 2026-11-10T10:31:38Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T17:26:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 3, behind S-0689, S-0716 and S-0719."
  by: flai
  at: 2026-11-11T14:52:10Z
finalized:
  by: alex
  at: 2026-11-11T08:59:51Z
---
# S-0699 A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context

## Goal

This story remediates [I-0261](../../../design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md), "A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0261 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0261 is closed with `flai issue close I-0261 --reason` saying what fixed it

## Tasks
- T-2992 A story's topic carries one claim source per item and sub-project, with the count of the paths that reach it
- T-2996 Measure S-0425's and S-0586's packs again, record the parts in agent-context.md, and close I-0261
- T-3111 A brief's JSON carries the count of its other reasons in place of the also list

## Notes

Cost of delay inputs set by flai from I-0261. time_lost_per_cycle 18m: 10m per occurrence × 2 occurrences ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:44:22Z, 7.5 days before this story).

### Planning

The planner (planner-S-0699) measured `flai prime --story --json` on main on 2026-11-10, with the 80 KB budget:

| Story | Pack | Topics | Conventions | Named | Briefs | ADR briefs one link step away | `also` lists |
|-------|------|--------|-------------|-------|--------|-------------------------------|--------------|
| S-0425 | 840 KB | 12 KB, 136 sources | 106 KB | 362 KB | 342 KB, 288 briefs | 167 KB, 155 briefs | 134 KB |
| S-0586 | 552 KB | 13 KB, 130 sources | 106 KB | 74 KB | 342 KB, 295 briefs | 164 KB, 160 briefs | 132 KB |

Later parts do not repeat the topic sources: since T-0987, every part after the first carries an empty `topics` list (`flai/internal/context/parts.go`, `open`). The sources cost part 1 alone. The packs are large because every brief loads on a pack over its budget (TH-0032, ADR-0241), because each brief's JSON lists all of its other reasons, and because a topic holds one source per claimed path.

The operator settled TH-1023 on 2026-11-11 by confirming the orchestrator's recommendation:

- Q1, link-step ADR briefs to the catalog: no. Every brief still loads, ADR-0241 stands, and the story writes no ADR.
- Q2, a count of a brief's other reasons in place of `also`: yes. T-3111 builds it. A loaded item keeps its list.

Tasks and layers:

1. T-2992 (topic sources folded per item and sub-project) and T-3111 (a brief's `also` as a count). They share no file and run together.
2. T-2996 measures both packs again, records them, and closes I-0261. It waits for T-2992 and T-3111, since it measures what they saved, and T-3111 also edits `design/system/agent-context.md`.

Touches:

- `flai/internal/topics/story.go`, `flai/internal/topics/story_test.go`: layout. `topics.OfStory` writes the sources, and T-2992 folds them.
- `flai/internal/context/context_test.go`, `flai/internal/mcpserver/prime_test.go`: layout and co-change. Each is where T-2992's pack-size bound can be tested.
- `flai/internal/context/items.go`, `flai/internal/context/context.go`, `flai/internal/context/design_test.go`: layout. `Selection.Items` sets a brief's `Also`, `Item.briefReason` counts it for the text pack, and `TestPrintedBriefsNameTheirFirstReasonAndCountTheRest` covers it. T-3111 changes them.
- `design/system/agent-context.md`: design. Its `## Fitting the pack to a budget` says "the JSON keeps every reason", which T-3111 changes, and its `### Delivered in parts` records the parts, which T-2996 measures.
- `design/system/flai-cli.md`, `docs/users/flai.md`: co-change, from `flai touches suggest`. The first lists the pack's JSON `items` fields, `also` among them; the second says how a brief counts its later reasons. T-3111 updates both.
- `design/issues/I-0261-…md`, `design/issues/summary.md`: declared by the second criterion, which closes the issue.
- `flai touches suggest` also lists `flai/cmd/prime.go`, `flai/internal/mcpserver/folder.go`, and `docs/users/flai-reference.md`. None names the `also` field or "briefs are never cut" in a way Q2 changes, and with Q1 no nothing else does, so they stay out.

Figures:

- Forecast: `flai forecast` gave 17m, which is 107 s per unit of size over 35 done improvement stories, times size 9. I keep it raised to 1h30m. The story has two Go tasks and a measuring task. Its close-out runs the integration and smoke tiers, as every story that touches `flai/` does, and those tiers alone took S-0672 8m and S-0669 20m. The delivery, 13:20, is where flai places 1h30m in the pull order.
- Cost of delay: 45 USD a week, as `flai cod` worked it out from flai's input of 18m lost per 168h cycle at 150 USD an hour. It stands. The input is the operator's, and it counts only the two recorded instances.

The story names no `after`: it reads nothing another open story writes.

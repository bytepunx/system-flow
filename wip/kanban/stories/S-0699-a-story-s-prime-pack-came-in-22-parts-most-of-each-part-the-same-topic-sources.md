---
id: S-0699
type: story
nature: improvement
title: A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context
status: backlog
owner: alex
created: 2026-11-10T08:55:29Z
updated: 2026-11-11T01:40:45Z
transitions: []
tags: [flai, cli, go]
touches: [flai/internal/topics/story.go, flai/internal/topics/story_test.go, flai/internal/context/context_test.go, flai/internal/mcpserver/prime_test.go, design/system/agent-context.md, design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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
  delivery: 2026-11-11T15:36:00Z
  basis: "Its own forecast of 1h30m; 24th in the pull order with an in-progress limit of 3, behind S-0676, S-0678, S-0674, S-0704, S-0681, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0686, S-0689, S-0690, S-0691, S-0692, S-0693, S-0694, S-0695, S-0696, S-0697 and S-0698."
  by: flai
  at: 2026-11-11T01:40:45Z
---
# S-0699 A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context

## Goal

This story remediates [I-0261](../../../design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md), "A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0261 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0261 is closed with `flai issue close I-0261 --reason` saying what fixed it

## Tasks
- T-2992 A story's topic carries one claim source per item and sub-project, with the count of the paths that reach it
- T-2996 Measure S-0425's and S-0586's packs again, record the parts in agent-context.md, and close I-0261

## Notes

Cost of delay inputs set by flai from I-0261. time_lost_per_cycle 18m: 10m per occurrence × 2 occurrences ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:44:22Z, 7.5 days before this story).

### Planning

The planner (planner-S-0699) measured `flai prime --story --json` on main on 2026-11-10, with the 80 KB budget:

| Story | Pack | Topics | Conventions | Named | Briefs | ADR briefs one link step away | `also` lists |
|-------|------|--------|-------------|-------|--------|-------------------------------|--------------|
| S-0425 | 840 KB | 12 KB, 136 sources | 106 KB | 362 KB | 342 KB, 288 briefs | 167 KB, 155 briefs | 134 KB |
| S-0586 | 552 KB | 13 KB, 130 sources | 106 KB | 74 KB | 342 KB, 295 briefs | 164 KB, 160 briefs | 132 KB |

Later parts do not repeat the topic sources: since T-0987, every part after the first carries an empty `topics` list (`flai/internal/context/parts.go`, `open`). The sources cost part 1 alone. The packs are large because every brief loads on a pack over its budget (TH-0032, ADR-0241), because each brief's JSON lists all of its other reasons, and because a topic holds one source per claimed path. The proposal and its two questions are on TH-1023. T-2992 needs no answer. The briefs tasks wait for it.

Touches:

- `flai/internal/topics/story.go`, `flai/internal/topics/story_test.go`: layout. `topics.OfStory` writes the sources, and T-2992 folds them.
- `flai/internal/context/context_test.go`, `flai/internal/mcpserver/prime_test.go`: layout and co-change. Each is where a pack-size bound can be tested.
- `design/system/agent-context.md`: design. Its `### Delivered in parts` section records the parts, and T-2996 measures there.
- `design/issues/I-0261-…md`, `design/issues/summary.md`: declared by the second criterion, which closes the issue.
- `flai touches suggest` also lists `docs/users/flai.md`, `design/system/flai-cli.md`, `flai/cmd/prime.go`, `flai/internal/mcpserver/folder.go`, and `docs/users/flai-reference.md`. They describe "briefs are never cut" and the `also` reasons. They join the touches only with the briefs tasks, if TH-1023 keeps them.

Figures:

- Forecast: `flai forecast` gave 17m, which is 107 s per unit of size over 30 done improvement stories, times size 9. I raised it to 1h30m. The story has three to four Go tasks, possibly an ADR, and a regenerated reference. Its close-out runs the integration and smoke tiers, as every story that touches `flai/` does, and those tiers alone took S-0672 8m and S-0669 20m today. The delivery is flai's 16:56 plus that difference.
- Cost of delay: 45 USD a week, as `flai cod` worked it out from flai's input of 18m lost per 168h cycle at 150 USD an hour. It stands. The input is the operator's, and it counts only the two recorded instances.

The story names no `after`: it reads nothing another open story writes.

---
id: S-0426
type: story
nature: research
title: Agent memory system
status: done
owner: alex
created: 2026-10-16T05:47:45Z
updated: 2026-10-16T20:56:07Z
transitions:
  - to: ready
    at: 2026-10-16T05:47:46Z
    by: alex
  - to: in-progress
    at: 2026-10-16T06:01:33Z
    by: agent-S-0426
  - to: review
    at: 2026-10-16T07:52:00Z
    by: agent-S-0426
  - to: done
    at: 2026-10-16T20:56:07Z
    by: alex
tags: []
topics: [cli]
touches: [design/system/agent-memory.md, design/system/agent-context.md, design/system/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6639
  priced_by: harness
  turns:
    - day: 2026-10-16
      ceremony: 2
      hand_edits: 2
      work: 45
  models:
    - model: claude-haiku-5-5
      input: 495868
      output: 42897
      cache_read: 0
      cache_write: 0
      cost: 0.131
    - model: claude-opus-5-5
      input: 570
      output: 240842
      cache_read: 28858440
      cache_write: 2048271
      cost: 21.688
  strategic:
    - kind: orchestrator
      seconds: 48690
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 274
          output: 4763
          cache_read: 61757805
          cache_write: 94456
          cost: 15.5062
        - model: claude-sonnet-5-5
          input: 10
          output: 64
          cache_read: 108674
          cache_write: 45754
          cost: 0.1264
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: planner-S-0426
    at: 2026-10-16T19:31:54Z
  value: 150
  by: planner-S-0426
  at: 2026-10-16T19:32:08Z
forecast:
  duration: 1h30m
  delivery: 2026-10-16T09:05:00Z
  basis: "flai forecast's 8m (232 s per size unit over 6 research stories, size 2) raised to 1h30m, since S-0413, a comparable research story, was forecast at 50m and this one adds a trial of local retrieval in a scratch directory; delivery is the start at 06:01Z plus 1h30m times flai's cycle factor of 2.03."
  by: planner-S-0426
  at: 2026-10-16T06:01:49Z
---
# S-0426 Agent memory system

## Goal

Research open source solutions for agent memory that improve the recall and search time for agents over a corpus of textual data.

We are looking to adapt designs and technology to solve this problem for flai or to adopt an existing open source solution via an optional, configured adapter.

## Acceptance criteria
- [x] Present evaluations of various open source solutions and/or technologies and architectures to solve for this
- [x] Consider state-of-the-art approaches that may not be directly available for flai, but could serve as inspiration and direction for a solution that would improve the speed and reduce the overhead that exists in flai today

## Tasks
- T-1795 Record how flai recalls and searches for agents today, what it costs, and the headings every candidate is judged under
- T-1796 Evaluate open source agent memory and retrieval solutions under the document's headings
- T-1797 Survey state-of-the-art memory and retrieval approaches flai cannot adopt directly but could learn from
- T-1798 Try the most promising candidates against this repository's corpus in a scratch directory and record recall and search time
- T-1799 Recommend adapting or adopting, name the follow-on stories, and link the finding from agent-context.md

## Notes

### Planning

Touches, by where each came from:

| Touch | From | Why |
|-------|------|-----|
| `design/system/agent-memory.md` | design | The finding, a new design document in the shape of `agent-context.md` (S-0125) and `agent-prose.md` (S-0413) |
| `design/system/README.md` | layout | The index of `design/system/` gets the finding's row |
| `design/system/agent-context.md` | design, co-change seed | The pack's design gets one line linking the finding |

The story declared no touches. `flai touches suggest` seeded with `agent-context.md` and `design/tech/search.md` lists mostly the CLI's documentation (`flai-cli.md`, `docs/users/flai.md`), which implementation stories changed; a research story changes none of them, so none was taken. No folder touch was kept. `design/tech/search.md` is left out: the research adds no dependency.

`design/system/README.md` held S-0426 on overlap with S-0413; flai asked S-0413's agent about it in MS-0127, and S-0426 started at 06:01Z.

Topics: `cli` added, since every task reads flai's pack and search code (`flai/internal/context/`, `flai/internal/search/`) and its design.

Forecast: `flai forecast` gave 8m from the median 232 s per size unit over 6 research stories, at size 2. Raised to 1h30m: S-0413, a comparable research story, was forecast at 50m, and this one adds a trial in a scratch directory. Delivery 2026-10-16T09:05Z is the start at 06:01Z plus 1h30m times flai's cycle factor of 2.03; it was 10:00Z while the hold was expected to last until S-0413 reached review.

Cost of delay: 150 USD a week, as `flai cod` works it out and unadjusted. The input, `time_lost_per_cycle: 1h`, is the operator's: they confirmed it on TH-0475 at 2026-10-16T19:31Z, and the planner wrote it as confirmed. That is 1h lost per 168h cycle at `planning.hour_rate` 150.

Tasks run in five layers, one each, because all five write `design/system/agent-memory.md` in order: T-1795, T-1796, T-1797, T-1798, T-1799.

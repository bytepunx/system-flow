---
id: S-0725
type: story
nature: remediation
title: The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container
status: done
owner: alex
created: 2026-11-11T18:16:51Z
updated: 2026-11-12T16:04:52Z
transitions:
  - to: ready
    at: 2026-11-12T15:53:27Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-12T15:53:41Z
    by: agent-S-0725
  - to: review
    at: 2026-11-12T15:57:39Z
    by: agent-S-0725
  - to: done
    at: 2026-11-12T16:04:52Z
    by: orchestrator
tags: [docs]
topics: [harness, docker, orchestration]
touches: [design/system/agent-adapters.md, docs/operators/index.md, design/issues/I-0350-the-container-harness-s-design-and-operator-docs-still-say-the-close-out-and.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 250
  priced_by: harness
  turns:
    - day: 2026-11-12
      ceremony: 4
      hand_edits: 1
      work: 33
  models:
    - model: claude-opus-5-5
      input: 84
      output: 11637
      cache_read: 7699567
      cache_write: 224467
      cost: 3.5687
  strategic:
    - kind: planner
      seconds: 209
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 82
          output: 331
          cache_read: 2655697
          cache_write: 194594
          cost: 0.1948
        - model: claude-opus-5-5
          input: 212
          output: 12825
          cache_read: 9422165
          cache_write: 837977
          cost: 3.6507
    - kind: orchestrator
      seconds: 122
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 258
          cache_read: 1584065
          cache_write: 167185
          cost: 0.4982
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-11T18:16:51Z
  value: 12.5
  by: planner-S-0725
  at: 2026-11-12T15:52:23Z
forecast:
  duration: 17m
  delivery: 2026-11-12T23:03:00Z
  basis: "flai forecast: median 160 s per unit of size over 23 done remediation stories on claude-opus-5-5 in the small band, times size 6 (2 criteria, 4 touches), 8th in the pull order."
  by: planner-S-0725
  at: 2026-11-12T15:52:23Z
finalized:
  by: orchestrator
  at: 2026-11-12T15:53:12Z
---
# S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container

## Goal

This story remediates [I-0350](../../../design/issues/I-0350-the-container-harness-s-design-and-operator-docs-still-say-the-close-out-and.md), "The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0350 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0350 is closed with `flai issue close I-0350 --reason` saying what fixed it

## Tasks
- T-3138 Say in agent-adapters.md what the container harness has built since S-0618 and S-0628
- T-3139 Say in the operators guide that the orchestrator in a container verifies and accepts through MCP on the host
- T-3140 Close I-0350 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0350. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T16:04:40Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed fix: correct the two stale passages I-0350 names, each against what the story that changed it built, and close the issue. It is a documentation fix: no code changes.

Touches, each a file:

- `design/system/agent-adapters.md`: design, the issue's first instance. Its § A container harness of flai's own, paragraph "What is built.", ends "Still not done: the close-out in the image (S-0510, S-0520), and the `permission_prompt` check …". S-0628 (ADR-0274) built the first and S-0618 the second (`design/system/flai-cli.md` § The check inside an agent image (S-0618)).
- `docs/operators/index.md`: design, the issue's second instance. § Strategic agents in a container, **Where a story's code runs** and **What it cannot do yet**, say the orchestrator's shell `flai verify` and `flai accept` run in its container; ADR-0278 makes them refuse there in favour of the MCP tools `verify` and `accept`.
- `design/issues/I-0350-…md` and `design/issues/summary.md`: layout, written by `flai issue close` for the second criterion.

`flai touches suggest S-0725` from the two documents lists `design/system/flai-cli.md` and `docs/users/flai.md` (69% co-change) and `design/system/strategic-agents.md` (17%). They are left out: their passages on the container (`flai-cli.md` § The orchestrator's acceptance on the host, `docs/users/flai.md` on a role's container, `strategic-agents.md` § In a container (S-0509)) already say what ADR-0278 decided. No folder touch. No `after`: nothing this story reads is unbuilt.

Forecast 17m, as `flai forecast` gives it with the four touches written: two short prose edits and an issue close, the size of the remediation stories it is measured over. The first run, before the touches, gave 6m on size 2; the touches are the work, so the later figure stands.

Cost of delay value 12.50 USD a week, as `flai cod` gives it from flai's input of 5m lost per 168h cycle at 150 USD an hour. It stands: one occurrence, a reader misled, not a failure.

No test fits the first criterion: the cause is prose that no test checks against the code. T-3140 records that in the Notes.

### Criterion 1 without a test

The cause was two statements in prose that later stories made untrue: `design/system/agent-adapters.md` § A container harness of flai's own and `docs/operators/index.md` § Strategic agents in a container. No test in this repository checks a document's claims against the code, so none reproduces it. The criterion is met by the corrected passages, each checked against the code and ADRs it describes: ADR-0274 and `flai-cli.md` for the close-out in the image, `claudecheck.go` and `flai-cli.md` for the image's check, ADR-0278 and `strategic-agents.md` for the orchestrator's verify and accept.

### Accepted by the orchestrator

- Verified: 63b0b28dea70e990091a865b435788e5c6fd5830
- At: 2026-11-12T16:04:52Z

Verdict: pass
- 1: design/system/agent-adapters.md, docs/operators/index.md
- 2: design/issues/I-0350-the-container-harness-s-design-and-operator-docs-still-say-the-close-out-and.md, design/issues/summary.md

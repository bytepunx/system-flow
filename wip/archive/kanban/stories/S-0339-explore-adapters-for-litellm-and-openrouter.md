---
id: S-0339
type: story
nature: research
title: Explore Adapters for LiteLLM and OpenRouter
status: done
owner: alex
created: 2026-10-08T07:17:16Z
updated: 2026-10-08T08:37:05Z
transitions:
  - to: ready
    at: 2026-10-08T07:17:18Z
    by: alex
  - to: in-progress
    at: 2026-10-08T07:39:13Z
    by: agent-S-0339
  - to: review
    at: 2026-10-08T08:34:45Z
    by: agent-S-0339
  - to: done
    at: 2026-10-08T08:37:05Z
    by: alex
tags: []
topics: [cli, template]
touches: [design/system/agent-adapters.md, design/system/README.md, design/adrs, design/adrs/README.md, design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: medium
usage:
  source: log
  seconds: 2102
  turns:
    - day: 2026-10-08
      ceremony: 1
      hand_edits: 4
      work: 39
  models:
    - model: claude-fable-5-1
      input: 2307
      output: 126861
      cache_read: 12226315
      cache_write: 687734
      cost: 20.1623
    - model: claude-haiku-4-5-20251001
      input: 853006
      output: 55817
      cache_read: 3654023
      cache_write: 260277
      cost: 2.0628
  strategic:
    - kind: planner
      seconds: 276
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 122
          output: 6214
          cache_read: 816973
          cache_write: 90009
          cost: 0.2254
        - model: claude-opus-5-5
          input: 60
          output: 18971
          cache_read: 2831261
          cache_write: 123558
          cost: 1.9344
    - kind: orchestrator
      seconds: 468
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 66
          output: 954
          cache_read: 11648548
          cache_write: 33677
          cost: 2.8784
        - model: claude-sonnet-5-5
          input: 10
          output: 74
          cache_read: 136633
          cache_write: 47329
          cost: 0.1586
cost_of_delay:
  inputs:
    penalty_per_week: 500
    by: alex
    at: 2026-10-08T07:17:16Z
  value: 500
  by: planner-S-0339
  at: 2026-10-08T07:21:19Z
forecast:
  duration: 1h
  delivery: 2026-10-08T09:06:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0287 and S-0326."
  by: flai
  at: 2026-10-08T07:21:53Z
---
# S-0339 Explore Adapters for LiteLLM and OpenRouter

## Goal

Perform the research necessary to write a new epic to implement agent adapters for both LiteLLM and OpenRouter while exploring the necessary changes and abstractions to capture intent free from vendor specifics (Anthropic SDK API vs. OpenAI's).

## Acceptance criteria
- [x] A new epic to create agent adapters for LiteLLM and OpenRouter
- [x] New ADRs where necessary to capture changes necessary to move flai away from being Anthropic specific
- [x] Documentation capturing the abstractions and changes needed to expand support for other vendor's agents that can inform the planner that will write the epic's stories

## Tasks
- T-1333 Record where flai is tied to Anthropic and Claude Code today
- T-1334 Set out what LiteLLM and OpenRouter offer and which agent harnesses can drive them
- T-1335 Set out the vendor-neutral abstractions flai needs, with options and a recommendation
- T-1336 Report the options to the operator and record the decision in ADRs
- T-1337 Create the epic for the LiteLLM and OpenRouter adapters and link it from the finding

## Notes

### Planning

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `design/system/agent-adapters.md` | design | The finding document, new, in the shape of the research findings `release-signing.md` (S-0193) and `agent-coordination.md` (S-0124). It is the third criterion's documentation. |
| `design/system/README.md` | co-change, design | The index lists every finding. It is co-changed with the seed paths in 2 of 26 commits. |
| `design/adrs` | design | Folder touch, kept: the second criterion adds ADRs whose numbers and slugs `flai adr new` gives only once the operator decides, so no task can name them yet. It is in the manifest's `claims.shared`, so it holds no ready story (ADR-0096). |
| `design/adrs/README.md` | co-change | `flai adr new` adds each ADR's index row. It is co-changed in 4 of 26 commits. |

`flai touches suggest S-0339` found no declared touches. Seeded with `design/system/release-signing.md`, `design/system/agent-coordination.md`, and `flai/internal/harness/adapters.go`, it listed mostly code (`flai/internal/harness/`, `flai/internal/serve/`, `flai/cmd/serve_actions.go`) and `design/system/flai-cli.md`. Those are left out: a research story changes no code, and the living design changes when the epic builds the adapters. The epic is a work item that flai writes in `wip/` on the main branch, so it is no touch.

Figures:

- **Forecast: 1h, delivery 2026-10-08T09:45:00Z.** `flai forecast` gave 39m: 332 s per unit of size, the median over 4 done research stories, times size 7. That is raised to 1h. S-0193, a research story on one topic with a decision thread, took 35m. This story covers two vendors, several harnesses, more than one ADR, and an epic. The delivery adds the extra 21m to flai's 09:06Z and allows for the operator's answer on T-1336's thread.
- **Cost of delay: 500 USD a week.** `flai cod` gives 500 from the operator's penalty of 500 USD a week alone. It stands.

Tasks, in layers, one per layer, since each edits `design/system/agent-adapters.md` after the one before:

1. T-1333 records today's ties to Anthropic and Claude Code.
2. T-1334 sets out LiteLLM, OpenRouter, and the harnesses that can drive them.
3. T-1335 sets out the vendor-neutral abstractions, options, and a recommendation.
4. T-1336 reports to the operator and records the decision in ADRs.
5. T-1337 creates the epic, with no stories, for the planner to plan.

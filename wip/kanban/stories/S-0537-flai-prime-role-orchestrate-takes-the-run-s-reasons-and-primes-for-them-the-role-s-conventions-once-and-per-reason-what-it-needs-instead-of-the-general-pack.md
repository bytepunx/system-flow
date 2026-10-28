---
id: S-0537
type: story
nature: improvement
title: flai prime --role orchestrate takes the run's reasons and primes for them, the role's conventions once and per reason what it needs, instead of the general pack
status: backlog
parent: E-0030
owner: alex
created: 2026-10-28T21:02:39Z
updated: 2026-10-28T21:02:39Z
transitions: []
tags: [flai, prime, orchestrator]
topics: [orchestration, cli]
touches: [flai/internal/context, flai/internal/mcpserver, flai/cmd, design/system/flai-cli.md, design/system/strategic-agents.md, design/adrs, docs]
after: [S-0536]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0537 flai prime --role orchestrate takes the run's reasons and primes for them, the role's conventions once and per reason what it needs, instead of the general pack

## Goal

ADR-0209 §3. The orchestrator's general pack is three to four parts of some 35k characters each, loaded at every start and again after every compaction, whatever the run is for. With a run started for named reasons, the pack can be primed for them. `flai prime --role orchestrate --reason <reason>...` and MCP `prime` with `reasons` load the role's conventions once and, per reason, what it needs: for a thread, the thread, its item, and the conventions on threads; for a story in review, the story, its last `verify` result, and the acceptance rules; for a plan, the planner's thread and the items it wrote; for a release, `release_evaluate`'s figures; for a draft, the draft and the finalizing rules; for the ready column, the board and `promote_candidates`. The general pack is not loaded. flai serve's start prompt names the reasons in the form `prime` takes, and the digest's answers ranked against the reasons' threads (ADR-0194) stay within the budget's eighth.

## Acceptance criteria
- [ ] `flai prime --role orchestrate --reason` and MCP `prime` with `reasons` accept each reason the trigger names, load what the goal lists for it, and refuse an unknown reason naming the known ones, with a test per reason
- [ ] A pack primed for one thread is under a quarter of the general pack's size, measured in the test
- [ ] The orchestrator's start prompt names its reasons in the form `prime` takes, and the role's prompt tells it to prime with them
- [ ] `flai-cli.md`, `strategic-agents.md`, and the agent docs describe the reason-primed pack
- [ ] The orchestrator in this repository runs a day on reason-primed packs, and the narrative records the context tokens per action from `flai stats` beside the previous story's figure

## Tasks

## Notes

A reason the trigger gains later needs its pack rule here; the refusal naming the known reasons makes a missing one visible at once.

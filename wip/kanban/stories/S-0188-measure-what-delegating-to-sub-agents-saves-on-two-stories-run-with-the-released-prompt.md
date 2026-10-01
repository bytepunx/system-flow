---
id: S-0188
type: story
nature: research
title: Measure what delegating to sub-agents saves on two stories run with the released prompt
status: ready
owner: arobson
created: 2026-10-01T08:13:37Z
updated: 2026-10-01T08:13:53Z
transitions:
  - to: ready
    at: 2026-10-01T08:13:53Z
    by: alex
tags: [flai]
topics: [conventions]
touches: [design/system/agent-context.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0188 Measure what delegating to sub-agents saves on two stories run with the released prompt

## Goal

S-0175 told the agents `flai serve` starts with `claude-code` to hand search, test runs, and the check before review to an explorer and a verifier sub-agent (ADR-0059, ADR-0060). It could measure only its own run, which delegated by hand, and S-0118, which delegated once on its own, because no story runs with the new prompt until a flai with it is released and installed on the host (TH-0042). Measure two stories that run with the released prompt, against comparable runs that did not delegate, so the design says what delegation costs and saves.

## Acceptance criteria
- [ ] Two stories run by `flai serve` with a flai that includes S-0175, each delegating to the explorer or the verifier, are measured from their logs as `design/system/agent-context.md § Sub-agents` measures S-0175 and S-0118: cost, cache reads, model calls, and turns for the story's agent and its sub-agents, and wall time
- [ ] Each is set against comparable runs that did not delegate, with what makes them comparable (nature, size, components touched)
- [ ] The table and what it shows are in `agent-context.md`; if delegation costs more than it saves, a story or an open question says what to change (a cheaper sub-agent model, fewer delegations, or none)

## Tasks

## Notes

- Follow-up of S-0175, agreed in TH-0042.
- ADR-0051's usage reading includes sub-agents in a run's `result`; a sub-agent's own calls carry `parent_tool_use_id` in the log.

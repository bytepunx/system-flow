---
id: S-0430
type: story
nature: feature
title: flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:58:39Z
updated: 2026-10-16T06:05:07Z
transitions: []
tags: [cli, template]
topics: [cli, conventions, orchestration]
after: [S-0429]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0430 flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's

## Goal

What a past story decided, or the operator answered, reaches the next agent that meets the same question without anyone writing it into a design document first. `flai prime --story` and the MCP `prime` add a digest step: the active and compacted entries whose topics meet the story's topics or whose paths meet its claim, ranked against the story with the rest within the budget, each marked as the digest's with its reason. An absorbed entry is never primed. `flai prime --role orchestrate` gives the orchestrator the `answer` and `recommendation` entries ranked against the threads awaiting the operator, so that its recommendation can cite the answer the operator gave last time as its `source`. The conventions tell an agent to cite an entry's ID when it relies on one, which is what keeps an entry from being compacted.

## Acceptance criteria

- [ ] An ADR, refining ADR-0049 and ADR-0075, records the digest step: which entries a story's pack selects (topics, paths against the claim), where they rank among the ranked sections, how a compacted entry loads (its one-line title, as a brief), that an absorbed entry never loads, how much of the budget the step may take, and what the orchestrator's pack selects; the sub-agents' packs take entries as briefs only.
- [ ] `flai/internal/context` adds the step: each entry carries `step: digest`, its reason (`topics: ...`, `path: ...`, `rank n`), and its size in text and `--json`; a test packs a fixture story whose topics select one entry, whose touches select another, and whose budget leaves out a third, and a fixture with an absorbed entry that never loads.
- [ ] `flai prime --role orchestrate`, as text, `--json`, and the MCP `prime` with role `orchestrate`, includes the answer and recommendation entries ranked against the titles and first entries of the threads awaiting the operator, within the budget; a test shows an entry ranked first for a thread that asks its question again.
- [ ] A `source` on `thread_reply` may name an entry's path, and the orchestrator's prompt and the convention `strategic-agents.md` say to cite the entry when one answers the question.
- [ ] The convention `agent-narrative.md` (template and here) says an agent cites `D-nnnn` in its narrative or thread when it relies on an entry, and asks the operator only what no entry answers; `design/system/agent-context.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` describe the step.
- [ ] Measured on this repository after the seed: the pack sizes of three archived stories with and without the step, recorded in `agent-context.md`, and no pack over its budget by the step alone.

## Tasks
- T-1818 An ADR records the pack's digest step: what a story's pack and the orchestrator's select, how entries rank and load, and what part of the budget they may take
- T-1819 A story's pack and a sub-agent's add the digest step: entries selected by topics and claim, ranked within the budget, each marked as the digest's
- T-1820 The orchestrator's pack adds the answer and recommendation entries ranked against the threads awaiting the operator
- T-1821 flai prime and the MCP prime carry the digest step, a thread_reply source may cite an entry, and the orchestrator's prompt says to cite the answer the operator gave last time
- T-1822 The conventions tell an agent to cite a digest entry it relies on and to ask only what no entry answers, and agent-context.md describes and measures the step

## Notes

- Drafted by the planner for E-0025, after S-0429 so that there are entries to prime with and measure.
- Parts (ADR-0104) are unchanged: entries are items like the rest and go into parts as they fit.

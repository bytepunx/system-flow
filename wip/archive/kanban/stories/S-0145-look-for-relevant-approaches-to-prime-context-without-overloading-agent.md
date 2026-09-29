---
id: S-0145
type: story
nature: research
title: Look for relevant approaches to prime context without overloading agent
status: done
parent: E-0010
owner: alex
created: 2026-09-29T01:57:32Z
updated: 2026-09-29T02:29:19Z
transitions:
  - to: ready
    at: 2026-09-29T01:57:54Z
    by: alex
  - to: in-progress
    at: 2026-09-29T02:12:32Z
    by: agent-S-0145
  - to: review
    at: 2026-09-29T02:23:03Z
    by: agent-S-0145
  - to: done
    at: 2026-09-29T02:29:19Z
    by: alex
tags: [cli]
touches: [design/adrs, design/issues, design/system/agent-context.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
---
# S-0145 Look for relevant approaches to prime context without overloading agent

## Goal

In order to prime agents with the relevant context for the story, we have to pull a lot of text from disparate sources. Even with some improvements to filtering context can be well over 537 KB (several hundred tokens). The agents can't load that much context in a single run.

Let's do some additional research on different ways to make sure only relevant documentation is fed to the agent. Make sure to include all approaches (RAG, MCP, Memory, etc.) with details about their trade-offs. The output should be an ADR draft.

## Acceptance criteria
- [x] Multiple potential approaches identified and analyzed
- [x] Recommendation for approach made as apart of the ADR draft
- [x] The approach has the potential to significantly reduce the number of tokens fed to the agent

## Tasks
- T-0517 Record I-0050's second instance and the narrowed claim
- T-0518 Measure where the context pack's size comes from
- T-0519 Survey approaches to prime an agent with relevant context and their trade-offs
- T-0520 Write the ADR draft with the recommendation
- T-0521 Close out: criteria, design pointer, narrative, check

## Notes

The output is [ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md), proposed, refining ADR-0047: twelve approaches compared (heading topics, a budget, BM25, embeddings, on-demand MCP tools, harness hooks, memory, model-written digests, a scout sub-agent, prompt caching, a `context:` list, splitting documents) and a recommendation estimated to bring the S-0138 pack from 541 KB to about 80 KB, about 20k tokens, with today's conventions. `design/system/agent-context.md` points at it. If accepted, the designer sets ADR-0047 as superseded in part (the draft uses `refines` so that a proposed ADR does not mark an accepted one superseded) and writes the three build stories the ADR lists. The claim was narrowed from `flai/cmd` to `design/adrs`, `design/issues`, and `design/system/agent-context.md` to clear the hold behind S-0138 (I-0050, second instance).

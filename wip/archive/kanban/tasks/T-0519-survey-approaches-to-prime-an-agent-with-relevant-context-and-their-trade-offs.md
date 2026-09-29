---
id: T-0519
type: task
nature: feature
title: Survey approaches to prime an agent with relevant context and their trade-offs
status: done
parent: S-0145
owner: alex
created: 2026-09-29T02:12:50Z
updated: 2026-09-29T02:20:52Z
transitions:
  - to: ready
    at: 2026-09-29T02:16:46Z
    by: agent-S-0145
  - to: in-progress
    at: 2026-09-29T02:16:47Z
    by: agent-S-0145
  - to: done
    at: 2026-09-29T02:20:52Z
    by: agent-S-0145
stream: S-0145
tags: []
---

# T-0519 Survey approaches to prime an agent with relevant context and their trade-offs

## Work

Survey the ways to give an agent only the documentation its story needs: selection by metadata (today's topics and links), lexical and embedding retrieval (RAG) at prime time and on demand, MCP tools the agent calls to fetch or search (progressive disclosure), a memory or summary layer kept across sessions, tiered packs with a budget, and model-written digests. For each: what it removes from the pack, what it costs to build and to keep true, what it depends on, how it fails, and how it fits flai (Go, offline, deterministic, explainable). Cite sources already in `design/system/agent-context.md` and add new ones.

## Done when

Every approach the story names (RAG, MCP, memory) and the others found is described with its trade-offs in the ADR draft's alternatives, with a comparison table, and none is dismissed without a stated reason.

## Notes

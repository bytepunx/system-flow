---
id: S-0147
type: story
nature: feature
title: "Agents fetch a design section on demand: doc_search, and a heading on doc_get and flai doc show"
status: ready
parent: E-0010
owner: alex
created: 2026-09-29T03:08:15Z
updated: 2026-09-29T03:16:05Z
transitions:
  - to: ready
    at: 2026-09-29T03:16:05Z
    by: alex
tags: [cli]
touches: [flai/internal/mcpserver, flai/internal/search, flai/cmd, docs/users, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0147 Agents fetch a design section on demand: doc_search, and a heading on doc_get and flai doc show

## Goal

An agent primed with briefs reads the body it needs, and only that, as [ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) decides: it searches sections, and fetches one section by its heading rather than the whole file.

## Acceptance criteria
- [ ] The MCP tool `doc_search` and `flai doc search` rank sections of `design/`, `docs/`, and the conventions by BM25 against a query, with `flai/internal/search` over the same section cuts the pack makes, and return path, heading path, the first lines, and size, at most twenty.
- [ ] `doc_get` and `flai doc show` take a heading and return that section with the sections below it; without one they return the whole document as today; an unknown heading is refused naming the headings there are.
- [ ] The MCP tool list tests and the server's instructions name `doc_search`; `docs/users/flai.md`, the reference, and `design/system/flai-cli.md` document both; all three test tiers and `flai check --strict` pass.

## Tasks

## Notes

Second of the three stories that build ADR-0049, from TH-0030. It can be worked alongside the budget story; the two touch `flai/internal/search`, so one waits for the other.

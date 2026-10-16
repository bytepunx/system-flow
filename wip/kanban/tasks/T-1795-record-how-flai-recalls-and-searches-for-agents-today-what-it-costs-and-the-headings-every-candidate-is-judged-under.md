---
id: T-1795
type: task
nature: research
title: Record how flai recalls and searches for agents today, what it costs, and the headings every candidate is judged under
status: in-progress
parent: S-0426
owner: alex
created: 2026-10-16T05:57:39Z
updated: 2026-10-16T06:02:48Z
transitions:
  - to: ready
    at: 2026-10-16T06:02:48Z
    by: agent-S-0426
  - to: in-progress
    at: 2026-10-16T06:02:48Z
    by: agent-S-0426
stream: S-0426
tags: []
touches: [design/system/agent-memory.md, design/system/README.md]
---
# T-1795 Record how flai recalls and searches for agents today, what it costs, and the headings every candidate is judged under

## Work

Create `design/system/agent-memory.md`, the finding of S-0426, in the shape of `agent-context.md` (S-0125) and `agent-adapters.md` (S-0339), and add its row to `design/system/README.md`. S-0413 holds the README too; add only this document's row, and agree the split in the two stories' conversation if flai asks.

Write its first section, `## Today`: how an agent recalls the corpus now, from the code and the design, with paths.

- Context packs: `flai/internal/context/` (`Build`, `ForStory`, `LoadDocs`, `rank.go`, `budget.go`, `parts.go`), the 80 KB budget and 40,000-byte parts ([ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), [ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md), [ADR-0104](../../../design/adrs/0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)).
- Search: the hand-written BM25+ in `flai/internal/search/`, rebuilt from the files on every `prime` and `doc_search` call; the in-memory index only the dashboard's `search.query` keeps (`flai/internal/hostapi/searching.go`); `design/tech/search.md`.
- What it costs, measured on this repository now: `prime` and `doc_search` time, pack size, corpus size in files and bytes, and the figures already recorded in `agent-context.md` § Measured and `server-performance.md`, each with its date. Say what is not measured, such as a pack's cache-read cost.
- Where recall fails today: what an agent cannot find without a link or a matching word (synonyms, past decisions in the archive, thread answers), with examples from this repository.

Then write `## What a solution must meet`: the headings every candidate in the later tasks is judged under. At least recall quality, search latency, index build and refresh cost, footprint on disk and in memory, files staying the source of truth ([ADR-0020](../../../design/adrs/0020-files-plus-mcp.md)), running locally with no paid service, fit with Go as a library or behind an optional configured adapter, licence, and maturity. Name the consumers it serves: the pack, `doc_search`, E-0025's digest, and S-0427's synchronized store.

Waits for nothing: it creates the document and the yardstick the later tasks use.

## Done when

- `design/system/agent-memory.md` exists with front matter, an introduction, `## Today`, and `## What a solution must meet`.
- Every figure in `## Today` comes from a command run on this repository or cites the document it is taken from, with its date.
- `design/system/README.md` lists the document.

## Notes

---
id: T-1796
type: task
nature: research
title: Evaluate open source agent memory and retrieval solutions under the document's headings
status: done
parent: S-0426
owner: alex
created: 2026-10-16T05:57:59Z
updated: 2026-10-16T06:21:12Z
transitions:
  - to: ready
    at: 2026-10-16T06:10:35Z
    by: agent-S-0426
  - to: in-progress
    at: 2026-10-16T06:10:35Z
    by: agent-S-0426
  - to: done
    at: 2026-10-16T06:21:12Z
    by: agent-S-0426
stream: S-0426
tags: []
touches: [design/system/agent-memory.md]
after: [T-1795]
usage:
  source: log
  seconds: 637
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 102
      output: 516
      cache_read: 4434858
      cache_write: 144679
      cost: 2.1309
---
# T-1796 Evaluate open source agent memory and retrieval solutions under the document's headings

## Work

Add `## Open source candidates` to `design/system/agent-memory.md`: one comparison table under the headings T-1795 wrote, then a short subsection per candidate. Each claim cites the project's documentation or repository with the date it was read.

Cover at least three families, plus any better candidate found on the way:

- Agent memory systems: Mem0, Letta (MemGPT), Zep with Graphiti, Cognee, LangMem.
- Embedded vector and hybrid stores a Go binary could use or run beside it: sqlite-vec, SQLite FTS5, LanceDB, Chroma, Qdrant, chromem-go.
- Text search engines and libraries: Bleve, Tantivy, Zoekt, txtai.

For each, say what it would replace or add in flai, whether it can be adopted behind an optional configured adapter or only adapted as a design, which embedding models it needs and whether one runs locally, and what it would cost to keep in step with the files. Mark each candidate tried or not tried; T-1798 tries the most promising.

Waits for T-1795: it judges every candidate under the headings that task writes, in the same document.

## Done when

- `## Open source candidates` holds a table covering every candidate under every heading of `## What a solution must meet`, with a cell left blank only where the documentation says nothing, marked so.
- Every candidate cites its source with the date read.

## Notes

---
id: T-1797
type: task
nature: research
title: Survey state-of-the-art memory and retrieval approaches flai cannot adopt directly but could learn from
status: backlog
parent: S-0426
owner: alex
created: 2026-10-16T05:58:08Z
updated: 2026-10-16T05:58:08Z
transitions: []
stream: S-0426
tags: []
touches: [design/system/agent-memory.md]
after: [T-1796]
---
# T-1797 Survey state-of-the-art memory and retrieval approaches flai cannot adopt directly but could learn from

## Work

Add `## State of the art` to `design/system/agent-memory.md`: approaches from papers and products that may not be available to flai as a library, each with what it does, the evidence for it (the paper or post, its date, and the benchmark it reports), and what it would mean for flai's speed and overhead.

Cover at least:

- Retrieval: hybrid lexical and dense retrieval with reciprocal rank fusion, learned sparse retrieval (SPLADE), late interaction (ColBERT), contextual retrieval and reranking.
- Structure over the corpus: hierarchical summaries (RAPTOR), knowledge graphs (GraphRAG, HippoRAG), and how flai's own links, topics, and ADR refinements already form a graph.
- Agent memory: tiered memory and paging (MemGPT), self-organising notes (A-MEM), reflection and consolidation, decay and compaction as E-0025 proposes.
- Cost: prompt caching and stable pack prefixes, incremental indexes that update on file change instead of rebuilding per call.

For each, say which part of flai it bears on (the pack, `doc_search`, E-0025's digest, S-0427's store) and whether it needs a model at query time.

Waits for T-1796: it writes the next section of the same document, and names where an open source candidate already implements an approach.

## Done when

- `## State of the art` covers every approach above with its source and date, and names the part of flai each bears on.
- Each approach says whether it needs a model or a paid service at query time.

## Notes

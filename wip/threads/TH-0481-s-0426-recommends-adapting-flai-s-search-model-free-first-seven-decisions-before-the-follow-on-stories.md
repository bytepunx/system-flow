---
id: TH-0481
title: "S-0426 recommends adapting flai's search, model-free first; seven decisions before the follow-on stories"
anchor:
  path: wip/kanban/stories/S-0426-agent-memory-system.md
  item: S-0426
status: resolved
participants: [agent-S-0426, orchestrator, alex]
created: 2026-10-16T07:51:37Z
updated: 2026-10-16T20:56:07Z
---

# TH-0481 S-0426 recommends adapting flai's search, model-free first; seven decisions before the follow-on stories

On wip/kanban/stories/S-0426-agent-memory-system.md.

## Entries

### 2026-10-16T07:51:37Z agent-S-0426
The finding is `design/system/agent-memory.md` § Recommendation, on story/S-0426.

**Recommendation: adapt, do not adopt.** flai keeps its own index, derived from the files and deletable. No agent memory system and no search server is adopted: the memory systems fail ADR-0020 (their store is the record) and determinism (an LLM writes it).

**What the trial measured** (25 questions agents asked here, recall at 10):

| Approach | R@10 | Query |
|----------|------|-------|
| flai BM25, CLI (today) | 15 | 233 ms, about 200 ms of it the rebuild |
| flai BM25, index kept in memory | 15 | 6.9 ms |
| SQLite FTS5, porter | 18 | 2.6 ms |
| BM25 + local dense model, RRF | 19 | 57 ms; 271 s to build |

**Follow-on stories, in order:** (1) a recall set and a script that measures it; (2) sections over a size bound indexed in pieces, and `doc_get` returning a large section in pieces; (3) `flai mcp` and `flai serve` keep the section index in memory; (4) a stemmer and a project glossary; (5) room in a story's pack under its budget, and the same bytes for the same inputs; (6) an optional embedding adapter, off by default; (7) an experiment with personalized PageRank over flai's links.

**Decisions, recommended answer first:**

1. Adapt or adopt? Adapt, as above.
2. Bound section size in flai, or split large documents by hand? In flai, for every project; story 2 picks the bound (2, 4, or 8 KB) by the recall set.
3. Where does the glossary live? `design/system/glossary.md`, a table of a term and its synonyms, which flai reads when it builds the index; the template ships it empty.
4. Should `doc_search` read the archive and threads? Not directly: through E-0025's digest (S-0428 puts entries under `design/`). Revisit if S-0435 shows the digest misses what agents ask.
5. Offer the embedding adapter, and when? Yes, off by default, only if stories 1 to 4 leave the set below 19 of 25. It calls a process the operator runs (Ollama's `/api/embed` first) and needs an ADR refining ADR-0049.
6. Take a SQLite driver to persist the index? Not now: the in-memory index answers in 6.9 ms, and S-0427 recommended dropping the read model it would share.
7. How does a pack get room under its budget? Narrow the conventions by heading topics (TH-0032) and keep the 80 KB budget.

No story or ADR is written: tell me which stories to create, or create them at acceptance.

### 2026-10-16T07:51:54Z orchestrator (recommendation)
**Recommendation:** take all seven answers as agent-S-0426 gives them, and create follow-on stories 1 to 3 now, under a new epic for search recall:

1. a recall set and a script that measures it;
2. sections over a size bound indexed in pieces;
3. the section index kept in memory by `flai mcp` and `flai serve`.

Leave stories 4 to 7 until story 1's recall set can judge them.

## Why

- **Decisions 1 and 6 follow a recorded decision.** ADR-0020 says "Durable state stays in files", so a memory system whose store is the record cannot be adopted. A persisted SQLite index adds a driver for no gain: the in-memory index answers in 6.9 ms. I accepted S-0427 at 07:30Z, and it recommended dropping the SQLite read model this index would share.
- **Stories 1 to 3 need no new decision:**
  - Story 1 measures.
  - Story 3 alone takes `doc_search` from 233 ms to 6.9 ms at the same recall.
  - Story 2 picks its bound by story 1's figures.
- **Decisions 3 and 5 can wait.** The glossary and the embedding adapter are worth deciding only once story 1 shows whether the stemmer reaches 19 of 25.
- **Decision 4 waits on E-0025's digest.** It is planned on TH-0480 and awaits your inputs on TH-0476.

## Why you

Decisions 2 to 5 and 7 are not yet recorded, and which stories to create is scope. Neither is mine to settle.

## How this affects S-0426's acceptance

None of these answers is needed to accept S-0426. Its two criteria ask for evaluations, not decisions. While this thread is open, though, flai's dry run blocks the acceptance. If you would rather decide later, resolve this thread or answer it at the acceptance. I will then verify and accept S-0426 under accept_reviews once it reaches review.

Source: design/adrs/0020-files-plus-mcp.md § Decision

### 2026-10-16T20:56:07Z alex
Resolved: S-0426 was accepted

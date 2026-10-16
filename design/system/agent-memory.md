---
title: "Agent memory: recall and search over a project's documents"
updated: 2026-10-16
status: active
topics: [cli]
---

# Agent memory: recall and search over a project's documents

The finding of S-0426. The operator asked for research into open source agent memory: what makes an agent recall a corpus of text better and find it faster, so that flai adapts the design and the technology, or adopts an existing solution behind an optional, configured adapter. This document says how an agent recalls this project's documents today, what that costs, and where it fails, and sets the yardstick every candidate is judged by. It then evaluates open source candidates, the state of the art that flai cannot take directly, a trial on this repository, and ends with a recommendation. Nothing here is built.

Every figure about this repository comes from a command run on it, or names the document and the date it is taken from. Today's figures were taken on 2026-10-16 between 06:03Z and 06:07Z at commit `4282b31c`, on the operator's host (WSL2, 24 cores, load average 1.3 to 2.3), with `bin/flai` built from that commit by `scripts/flai.sh`; the build is not counted. Times are wall-clock from bash's `time`.

## Today

An agent reaches the project's documents three ways: the context pack flai pushes at the start of a story, the `doc_search` and `doc_get` tools it pulls with, and its harness's own file reads. The operator has a fourth, the dashboard's search. None of them keeps anything on disk: every index is built from the files.

### The context pack

`flai prime --story` and the MCP `prime` tool both call `ForStory` (`flai/internal/context/story.go`). It reads the conventions, the documents under `design/system`, `design/tech`, and `design/adrs` (`LoadDocs` in `design.go`), the story's topics, and the story, its epic, and its tasks, listing every work item, the archive included, to find the tasks. `Build` (`context.go`) keeps each convention's sections whose topics match the story's and adds the open-issues table. Then come what the story names, whole, the design system, briefs of what topics and one link step select, sections ranked by BM25 against the story's title, goal, and criteria to fill what is left (`rank.go`), and a catalog of the rest ([ADR-0047](../adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), [ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)). The budget is 80 KB unless `prime.budget` says otherwise (`budget.go`). Conventions, what is named, and briefs are never cut. The MCP tool returns the pack in parts of at most 40,000 bytes of JSON (`parts.go`, [ADR-0104](../adrs/0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)). [agent-context.md](agent-context.md#as-built) records how it was built.

Measured today:

| Command | Runs | Result |
|---------|------|--------|
| `bin/flai prime --story S-0426 --json` | 5 | 0.70 to 0.71 s; 0.72 to 0.76 s through `scripts/flai.sh` |
| `bin/flai prime --story S-0261 --json`, `--story S-0386` | 3 each | 0.68 to 0.71 s |
| `bin/flai prime --cat` (conventions only, no story) | 3 | 27 to 29 ms |
| `bin/flai show S-0426` (lists every work item too) | 3 | 0.53 to 0.55 s |
| `bin/flai version` | 3 | 8 to 11 ms |

S-0426's pack is 128,182 bytes in 1,115 lines, 152,383 bytes of JSON, five parts. Its fifteen conventions, the README among them, are 84,931 bytes against a budget of 81,920, so the header reports `exceeded: conventions`, and the pack is the conventions, the 8 open issues (2,706 bytes), and a catalog of 199 documents (38,202 bytes). It names, briefs, and ranks nothing. The packs of S-0261 and S-0386 are the same shape: 128,461 and 128,685 bytes, no items. For the three stories measured, no design reaches the agent through the pack; it reaches it only through `doc_search`, `doc_get`, and the agent's own reads.

The comparison with `prime --cat` and `flai show` says most of a CLI pack's 0.7 s goes on listing the work items, not on search. That is inferred, not profiled. `flai mcp` keeps parsed items between requests (S-0156, [server-performance.md](server-performance.md#cause-1-work-items-kept-between-requests-s-0156)), so the MCP tool pays the listing once per process. It still reads, parses, and cuts every document on every call.

### Search

`flai/internal/search` is a hand-written BM25+ index ([search.md](../tech/search.md)): lower-cased runs of letters and digits, prefix and fuzzy matching with MiniSearch's weights, terms combined with OR (`search.go`). It has two users.

- **Sections, for agents.** `doc_search` and `flai doc search` call `SearchSections` (`flai/internal/context/sections.go`). It reads every Markdown file under `design/` and `docs/`, folder READMEs and the ADR template left out, cuts each at every heading, indexes each section as one document with common English words dropped (`search/sections.go`), and returns at most 20 hits. The pack's ranked step uses the same cut over `design/system`, `design/tech`, and `design/adrs`. Both build the index on every call and keep nothing.
- **Files, for the operator.** The dashboard's `search.query` indexes whole files under `design/`, `wip/`, and `docs/` (`flai/internal/hostapi/searching.go`). `flai serve` keeps the index in memory and rebuilds it when an FNV hash of every file's path, modification time, and size, taken by a walk that only stats, changes. It covers the archive and the threads; the agents' search does not.

`bin/flai doc search "<query>"` took 0.21 to 0.22 s for each of four queries ("context pack budget", "agent memory", "embedding retrieval", "where does the operator keep the release signing key"), three runs each. It indexes 379 files, 3,665,698 bytes, with 2,841 heading lines between them (counted with `grep -c '^#{1,6} '`, which also counts a heading-like line in a code block).

Earlier figures, on this host:

| Figure | Value | Source and date |
|--------|-------|-----------------|
| MCP `prime`, median of 3 | 148 ms, pack 98 KB | [server-performance.md](server-performance.md#mcp-tools), 2026-09-29 |
| MCP `doc_search`, median of 3 | 68 ms | same |
| `search.query`, cold | 256 ms; its `search.index` phase 297 ms cold, as recorded | same, [§ Dashboard methods](server-performance.md#dashboard-methods) |
| Dashboard index rebuild, 544 documents | about 60 ms | [search.md](../tech/search.md), 2026-09-20 |
| A cli story's pack with no budget | 541 KB, about 135k tokens | [agent-context.md](agent-context.md#fitting-the-pack-to-a-budget), 2026-09-29 |
| `prime` results as JSON, eight stories | 161 to 198 KB | [agent-context.md](agent-context.md#delivered-in-parts), 2026-10-06 |

The corpus, counted today with `find <folder> -name '*.md' | wc -l` and `... | xargs -0 cat | wc -c` (folder READMEs included):

| Folder | Files | Bytes | `doc_search` | The pack | Dashboard search |
|--------|-------|-------|--------------|----------|------------------|
| `design/system` | 26 | 1,398,601 | yes | ranked, briefed, catalogued | yes |
| `design/adrs` | 163 | 773,864 | yes | same | yes |
| `design/tech` | 13 | 32,558 | yes | same | yes |
| `design/conventions` | 16 | 97,947 | yes | by topic, whole sections | yes |
| `design/issues` | 148 | 356,148 | yes | the open-issues table | yes |
| `design/analysis`, `design/experiments` | 5 | 22,303 | yes | no | yes |
| `docs` | 18 | 1,038,721 | yes | no | when toggled on |
| `wip/archive` | 2,390 | 6,517,808 | no | the story's own epic and tasks, for their links | yes |
| `wip/threads` | 469 | 1,426,023 | no | no | yes |
| `wip/messages` | 127 | 364,489 | no | no | yes |
| `wip/kanban`, `wip/agents` | 157 | 571,951 | no | the story, its epic and tasks | yes |

Two files are a fifth of what `doc_search` reads: `design/system/flai-cli.md` (399,965 bytes) and `docs/users/flai.md` (370,233 bytes). `flai-cli.md § Commands` is one section of 340,847 bytes, because it has no heading below it.

### What is not measured

- What a pack costs once it is in the context. It is read again from the cache on every turn; [agent-context.md](agent-context.md#measured) records 87k to 188k cache-read tokens per model call for story agents on 2026-10-01 and 2026-10-02, but not the pack's share of them.
- MCP against the CLI today. Today's times are the CLI's, process start included; `flai mcp` was not timed again.
- The dashboard's `search.query` today: `flai guard` refuses `flai hostapi` to a task sub-agent, so the 2026-09-29 figures stand.
- Memory: the resident size of `flai mcp` and `flai serve` with their indexes.
- Recall: no set of questions with known answers exists. The replays in [agent-context.md](agent-context.md#as-built) measure which ADRs the pack's selection reached (93 of 93 on 2026-09-29, 47 by ranking alone), not whether `doc_search` finds an answer.
- How often agents call `doc_search`, and how many calls it takes them to find a section.

### Where recall fails

Queries run today with `bin/flai doc search --limit 20 "<query>"`:

| Query | Where the answer is | What came back |
|-------|---------------------|----------------|
| how many stories may run concurrently | [workflow.md](workflow.md) § Columns and WIP limits | All of the top 5 from `agent-coordination.md`; `workflow.md` not in the 20 |
| concurrency cap | same | `workflow.md` not in the 20; "WIP limit" finds it fifth |
| token spend per story | [ADR-0051](../adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md), the tokens and cost an agent spent | ADR-0018 and ADR-0033, on login tokens, first; ADR-0051 not in the 20 |
| anti-pattern low level git commands hand-edited files | The operator's rule on TH-0420, one sentence of [ADR-0155](../adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md) § Context | `flai-cli.md § Commands` (340,847 bytes) first; ADR-0155 fifteenth |
| branch ref empty file after host crash | TH-0390, and I-0137 | I-0137's sections first, second, and third |

What fails, and why:

- **Words, not meaning.** BM25 matches tokens. Prefix and fuzzy matching widen the spelling, not the sense, so "concurrently" never reaches "WIP", and "token" ranks login tokens with spent ones.
- **Sections too large to read.** A section ends only at a heading. `flai-cli.md § Commands` matches most words because it holds most of them, and no `doc_get` result can carry it: Claude Code saves a result over 50,000 characters to a file ([ADR-0104](../adrs/0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)).
- **Out of reach.** The 381 archived narratives hold 2,199 `## Decisions` bullets (counted with `awk` over `wip/archive/agents/*.md`), and 462 of the 468 threads are resolved. `doc_search` reads none of them. What an agent can find is what someone wrote into a design document, an ADR, or an issue. In the three threads checked (TH-0390, TH-0420, TH-0429) the answer had reached I-0137, ADR-0155, and `agent-adapters.md` with I-0144, but only as one clause in TH-0420's case, and nothing checks that it always does. That is E-0025's case. The operator's search covers the archive and the threads; the agent's does not.
- **No room to push.** Links and topics select design only when the conventions leave room in the budget, and today they leave none.

## What a solution must meet

A solution serves these consumers:

| Consumer | Today | What it needs from memory |
|----------|-------|---------------------------|
| The pack, `prime` | Built from the files on every call; conventions only for the stories measured | Rank sections against a story, within the budget, the same way every time |
| `doc_search` and `doc_get` | BM25 over sections of `design/` and `docs/`, built on every call | Find the section that answers an agent's question in its own words, and a section small enough to read |
| E-0025's digest | Not built (E-0025) | Entries drawn from the archive and the threads, each citing its source, that `doc_search` finds and the pack primes by topic |
| S-0427's synchronized store | Not built (S-0427) | A database populated from the files deterministically and kept current with them, which could hold the index too |
| The dashboard's `search.query` | Whole-file index kept in `flai serve` | The same recall for the operator, over the same files |

Every candidate in the sections that follow is judged under these headings, used as table columns:

| Heading | What it means for flai | The bar today |
|---------|------------------------|---------------|
| Recall | Finds the section that answers a question whose words differ from the document's, and reaches the archive and the threads when asked to | Misses the queries in [Where recall fails](#where-recall-fails); the pack's selection reached 93 of 93 ADRs (2026-09-29) |
| Latency | Time to answer one query, cold and warm, in the CLI and in `flai mcp` | `doc search` 0.21 to 0.22 s in the CLI; MCP `doc_search` 68 ms (2026-09-29); `prime` 0.70 s in the CLI, 148 ms over MCP (2026-09-29) |
| Build and refresh | Time to index the corpus from nothing, and to take in a changed file; whether a query can ever see a stale index | Built from the files on every call, so never stale; the dashboard's index rebuilt in 297 ms cold (2026-09-29) |
| Footprint | Bytes on disk, resident memory in `flai mcp` and `flai serve`, and bytes added to the binary | Nothing on disk; memory not measured; `bin/flai` is 31,310,038 bytes (development build, `ls -l`) |
| Files as source | The files stay the record ([ADR-0020](../adrs/0020-files-plus-mcp.md), which rejected a second store to keep consistent with them). An index or store is a cache derived from the files that can be deleted and rebuilt, and nothing is recorded only in it | Met: every index is derived |
| Local and free | Runs on the operator's host or in the dashboard's container with no paid service and no API key; any model runs locally | Met: no model, no network |
| Go fit | Either a Go library that builds with `CGO_ENABLED=0` for linux, darwin, and windows on amd64 and arm64, as releases do (`flai/.goreleaser.yaml`), or an optional adapter, set in the manifest or the config, that flai calls as a separate process, with today's search as the fallback when it is off or down | Met: about 470 lines of Go in `flai/internal/search`, no dependency |
| Licence | flai ships under the Bytepunx Shield License ([LICENSE.md](../../LICENSE.md)). A library linked into flai needs a permissive licence (MIT, BSD, Apache-2.0); a copyleft one (GPL, AGPL) fits only as a separate program behind the adapter, installed by the operator | Met: no third-party code |
| Maturity | Releases, maintainers, activity in the last year, users, and a storage format that stays readable across versions | flai's own code, tested in `search_test.go` and `sections_test.go` |
| Determinism | The same files and the same query give the same result, and a result says why it ranked, as every pack item carries its reason ([agent-context.md](agent-context.md#candidate-designs)) | Met: BM25 with ties broken by path |

Conventions are rules, and a rule that is missed is a rule that is broken ([agent-context.md](agent-context.md#what-the-evidence-says)): they stay selected by topic, whatever the memory. A candidate competes for the design, the decisions, the archive, and the threads.

## Open source candidates

Twenty candidates in four families were judged on 2026-10-16 from their READMEs and documentation and from the GitHub, PyPI, and Go module proxy records of their latest release, all read that day. None was run. The families are agent memory systems, which keep what an agent learns in a store of their own; embedded vector and hybrid stores; text search engines and libraries; and local embedding runtimes, which every vector or hybrid design needs to turn a section and a query into vectors. In the tables, "not stated" marks a heading the documentation read says nothing about, and "not confirmed" marks something it implies but that was not checked. Recall and latency are what each project says of itself; none of its figures was taken on this corpus. Stars, releases, and dates are GitHub's, PyPI's, or the Go proxy's on 2026-10-16.

The headings of [What a solution must meet](#what-a-solution-must-meet) are split over two tables with the same rows.

| Candidate | Family | Recall | Latency | Build and refresh | Footprint | Files as source |
|-----------|--------|--------|---------|-------------------|-----------|-----------------|
| Mem0 | Agent memory | Semantic, BM25, and entity signals fused; LoCoMo 92.5 on the managed platform, not the open source | p50 0.88 to 1.09 s, managed platform | An LLM call per add; add-only, nothing overwritten | Python or Node, an LLM, Qdrant at `/tmp/qdrant`, SQLite history | Fails: memories live only in its store |
| Letta Code (MemGPT) | Agent memory | The agent rewrites its own memory blocks; message search | Not stated | The agent and its "dreaming" rewrite memory | Node.js harness; Letta Cloud by default | Fails: memory the agent writes, tracked in git |
| Graphiti (Zep) | Agent memory | Semantic, BM25, and graph traversal | Sub-200 ms claimed for Zep's platform, not for Graphiti | Incremental; LLM extraction per episode, bound by provider rate limits | Python and Neo4j, FalkorDB, or Neptune | Fails: an LLM-extracted graph |
| Cognee | Agent memory | Graph, vector, or code context chosen per query | Not stated | Graph extraction per document, by GLiNER locally or by an LLM | Python, SQLite, a vector store, a graph store | Fails as memory; rebuildable from files only by re-extraction |
| LangMem | Agent memory | Vector search over what the agent chose to store | Not stated | A background manager extracts and consolidates | Python on LangGraph; in-memory or Postgres store | Fails: the agent decides what is stored |
| SQLite FTS5 | Embedded store | BM25 over words; porter and trigram tokenizers | Not stated; the cgo-free port is slower than C in CPU-bound work | Incremental insert and delete; external-content tables | One SQLite file; driver size not stated | Met as a derived cache |
| sqlite-vec | Embedded store | Vectors only | Not stated | Insert and delete rows; one embedding per section | One small C extension | Met as a derived cache |
| SQLite vec1 | Embedded store | Vectors only, approximate (IVFADC with OPQ) | Not stated | The index is trained (its roadmap) | One C file | Met as a derived cache |
| LanceDB | Embedded store | Vector, BM25 full text, hybrid fused by RRF | Not stated at this scale | A full-text index must be built before search | Lance columnar files; a native library for Go | Met as a derived cache |
| Chroma | Embedded store | Vector; hybrid and full text named only for Chroma Cloud | Not stated | Add, update, delete | A server, model files | Met as a derived cache |
| Qdrant | Embedded store | Dense, sparse, hybrid fused by RRF or DBSF | Not stated | Collections updated in place | A server (Docker) | Met as a derived cache |
| chromem-go | Embedded store | Vectors only, exact cosine; `$contains` filter | 1,000 documents 0.3 ms, 100,000 in 40 ms, 2020 laptop CPU | Add; one file per document written at once | In memory; gob files, gzip optional | Met as a derived cache |
| kelindar/search | Embedded store | Vectors only, exact cosine | Not stated | Add; scans every vector | In memory; one file | Met as a derived cache |
| coder/hnsw | Embedded store | Vectors only, approximate (HNSW) | Search O(log n) | Insert O(log n), delete O(M² log n) | In memory | Met as a derived cache |
| Bleve | Text search | BM25 or tf-idf, fuzzy, synonyms; hybrid by RRF or RSF only with vectors | Not stated | Index and delete per document | An index directory | Met as a derived cache |
| Tantivy | Text search | BM25, phrases, stemming in 17 languages | Start-up under 10 ms | Incremental, multithreaded | An mmap index directory | Met as a derived cache |
| Zoekt | Text search | Substring and regexp over trigrams, ranked for code; BM25 experimental | Not stated | Re-index per repository or directory | Index shards on disk | Met as a derived cache |
| txtai | Text search | Dense, BM25, or hybrid | Not stated | Not stated | Python; Faiss index by default | Met as a derived cache |
| Ollama | Embedding runtime | Embeddings only; as good as the model | Not stated; an HTTP call per query | Re-embed every changed section | A server; models 46 MB (all-minilm) to 4.7 GB (qwen3-embedding) | Met: vectors derived; a rebuild needs the same model |
| hugot | Embedding runtime | Embeddings only (feature extraction) | Go backend best at about 32 inputs per call | Re-embed every changed section | ONNX model files; binary size not stated | Met: vectors derived; a rebuild needs the same model |

| Candidate | Local and free | Go fit | Licence | Maturity | Determinism | Tried |
|-----------|----------------|--------|---------|----------|-------------|-------|
| Mem0 | With Ollama for the LLM and the embedder; defaults to OpenAI | Python or Node; a server with a REST API: adapter only | Apache-2.0 | v2.2.1, 2026-09-25; 67k stars | No: an LLM writes the memories | Not tried |
| Letta Code (MemGPT) | A local backend is offered; Cloud is the default; own LLM keys | Node harness: design only | Apache-2.0 | v0.34.15, 2026-10-15; the V1 Python server is archived | No | Not tried |
| Graphiti (Zep) | Through Ollama's OpenAI-compatible endpoint; defaults to OpenAI; Zep itself is a paid cloud | Python, with an MCP server: adapter only | Apache-2.0 | v0.31.1, 2026-10-15; 31.8k stars | No: LLM extraction | Not tried |
| Cognee | Yes: GLiNER and a local embedding model with no key | Python, with an API and an MCP server: adapter only | Apache-2.0; its production Postgres graph store is a licensed product | 1.6.4 on PyPI, 2026-10-13; 32.4k stars | Not confirmed for GLiNER; no with an LLM | Not tried |
| LangMem | Its README asks for an LLM provider's key; local not confirmed | Python library, no server: design only | MIT | 0.0.30, 2025-10-27 | No: the agent decides | Not tried |
| SQLite FTS5 | Yes | Library: `modernc.org/sqlite` or `ncruces/go-sqlite3`, both cgo-free | SQLite public domain; modernc BSD-3-Clause; ncruces MIT | SQLite 3.53.4 in modernc v1.60.1, 2026-09-29; ncruces v0.35.6, 2026-09-23 | Yes: a `bm25()` score per row | Not tried |
| sqlite-vec | Yes; needs a local model | Cgo, or ncruces Wasm bindings of 2024 whose fit with ncruces today is not confirmed | Apache-2.0 | v0.1.9, 2026-03-31; pre-v1, breaking changes expected | A distance; exact or approximate not stated | Not tried |
| SQLite vec1 | Yes; needs a local model | Library through ncruces `ext/vec1`, cgo-free | Extension not confirmed; ncruces MIT | Version 0.7; "testing is insufficient" | No: approximate | Not tried |
| LanceDB | Yes; needs a model | Go SDK needs cgo and native libraries: fails | Apache-2.0 | v0.40.0, 2026-10-07; Go SDK v0.1.2 | Approximate index; RRF rank | Not tried |
| Chroma | Yes; default model all-MiniLM-L6-v2 runs locally | A server and a Go client (`chroma-go`, MIT): adapter only | Apache-2.0 | 1.5.9, 2026-05-05; 29.5k stars | Not stated | Not tried |
| Qdrant | Yes, self-hosted; needs a model | A server and a Go client: adapter only | Apache-2.0 | v1.19.2, 2026-10-05; 35.1k stars | Fused scores; index exactness not stated | Not tried |
| chromem-go | With Ollama or LocalAI; defaults to OpenAI | Pure Go, no dependencies: library | MPL-2.0, file-level copyleft | v0.7.0, 2024-09-01; beta until v1 | Yes: exact cosine | Not tried |
| kelindar/search | With its llama module | Core pure Go, but needs Go 1.27 and flai is on 1.26; the llama module needs native llama.cpp | MIT | v0.5.0, 2026-09-19 | Yes: exact cosine | Not tried |
| coder/hnsw | Yes; needs a model | Pure Go: library | CC0-1.0 | v0.6.1, 2024-06-14 | No: approximate | Not tried |
| Bleve | Yes | Go; text search without cgo not confirmed; vectors need FAISS (C++) | Apache-2.0 | v2.6.1, 2026-08-24; 11.2k stars | BM25 score; hybrid fused | Not tried |
| Tantivy | Yes | Rust: a separate process (`tantivy-cli`) or design only | MIT | 0.26.1, 2026-05-10; 16.2k stars | BM25 score | Not tried |
| Zoekt | Yes | Go commands and a JSON API; library use not confirmed | Apache-2.0 | No releases; pseudo-version of 2026-10-15 | Score with `DebugScore` detail | Not tried |
| txtai | Yes; default model all-MiniLM-L6-v2 runs locally | Python API with a Go client (`txtai.go`): adapter only | Apache-2.0 | v9.14.0, 2026-10-06; 13k stars | Not stated | Not tried |
| Ollama | Yes | A separate process with an HTTP API: adapter | MIT | v0.40.5, 2026-10-14; 183k stars | Not confirmed | Not tried |
| hugot | Yes | Go backend without cgo; ONNX Runtime and XLA backends need cgo | Apache-2.0 | v0.8.3, 2026-10-14 | Not stated | Not tried |

### Agent memory systems

All five keep what an agent learns in a store of their own, written by an LLM or by the agent. That fails [Files as source](../adrs/0020-files-plus-mcp.md): the memory is the record, not a cache derived from the files. Each could serve only as a derived index, rebuilt from the files, and then their extraction step costs an LLM call per document and is not deterministic. What they offer flai is design: fused retrieval signals, entity linking, and time-aware facts, which E-0025's digest could take.

#### Mem0

A memory layer that extracts facts from conversations with an LLM and retrieves them by semantic, BM25, and entity signals fused in one pass; BM25 needs its `nlp` extra and spaCy. It would replace E-0025's digest, not `doc_search`: it stores what an LLM extracted, not sections of the files. It defaults to OpenAI's `gpt-5-mini` and `text-embedding-3-small`; its documentation shows Ollama as a local embedder, and a local Qdrant at `/tmp/qdrant` is its default store. flai could reach its self-hosted server over REST behind an optional adapter, but keeping it in step means an LLM call per changed file, and its benchmark figures are for the managed platform, "not available in the open-source SDK". Sources: [README](https://github.com/mem0ai/mem0) and [open source overview](https://docs.mem0.ai/open-source/overview), read 2026-10-16.

#### Letta Code (MemGPT)

Letta's current product is a stateful agent harness in TypeScript; the Python V1 server is retired to an `archive` branch. Its agents rewrite their own memory blocks and skills, "dream" between sessions, and search past messages; MemFS tracks all of an agent's context in git. It would replace the agent, not a part of flai, so it fits only as design. MemFS is the idea worth taking: memory as files under version control is close to [ADR-0020](../adrs/0020-files-plus-mcp.md), though here the agent, not the project, writes them. Sources: [letta-ai/letta](https://github.com/letta-ai/letta) and [letta-ai/letta-code](https://github.com/letta-ai/letta-code), read 2026-10-16.

#### Graphiti (Zep)

Graphiti builds a temporal knowledge graph from "episodes" with an LLM, invalidates facts that a later episode contradicts, and searches by embeddings, BM25, and graph traversal. Zep is the paid cloud built on it; Zep's Community Edition is no longer supported. It needs Neo4j, FalkorDB, or Neptune, and defaults to OpenAI for extraction and embedding; local models are reached through their OpenAI-compatible endpoints. It ships an MCP server, so it could sit behind an optional adapter, but the graph is extracted, not derived, and ingestion is bound by the LLM's rate limits. Its temporal facts, each with when it held, are the design to take for E-0025's digest of the archive and the threads. Sources: [getzep/graphiti](https://github.com/getzep/graphiti) and [getzep/zep](https://github.com/getzep/zep), read 2026-10-16.

#### Cognee

Cognee turns documents into a knowledge graph and embeddings and answers by choosing graph, vector, or code context per query. With no LLM key it extracts with the local GLiNER model and embeds with a local model, both downloaded on first use, so it runs locally and free. It keeps SQLite or Postgres, a vector store (LanceDB, PGVector, ChromaDB), and a graph store (Kuzu, Neo4j, Memgraph). It offers a CLI, an API on port 8000, and an MCP server on 8001, so it fits behind an optional adapter; keeping it in step with the files means extracting every changed document again. Sources: [README](https://github.com/topoteretes/cognee) and [setup and configuration](https://docs.cognee.ai/setup-configuration/overview), read 2026-10-16.

#### LangMem

A Python library of memory tools for LangGraph agents: the agent decides what to store and searches it by embeddings, and a background manager consolidates. Its README configures OpenAI's embedding model and asks for an LLM provider's key. There is no server, so flai could take it only as design, and what it stores is what an agent chose, not the files. It has had no PyPI release since 2025-10-27. Source: [langchain-ai/langmem](https://github.com/langchain-ai/langmem), read 2026-10-16.

### Embedded vector and hybrid stores

Every store here can be a deletable cache derived from the files, which meets [Files as source](../adrs/0020-files-plus-mcp.md). What separates them is whether flai can link them and whether they bring their own keyword search. Only the SQLite family links into flai with cgo off and holds both BM25 and vectors in one file, which is also a shape S-0427's synchronized store could take.

#### SQLite FTS5

SQLite's full-text extension, with a built-in `bm25()` ranking, porter and trigram tokenizers, `snippet()` and `highlight()`, and external-content tables that index rows held elsewhere. It would replace flai's hand-written BM25 for `doc_search`, the pack's ranked step, and the dashboard's `search.query`, and persist the index so a call need not rebuild it. It adds no meaning: it is the same class of recall as today. Two cgo-free drivers carry it: `modernc.org/sqlite`, a transpiled port that lists linux, darwin, and windows on amd64 and arm64 and is slower than C in CPU-bound work, and `ncruces/go-sqlite3`, a Wasm build translated to Go, tested on the same six targets, with an `ext/fts5` package. Keeping it in step means re-indexing the sections of changed files, found by the hash `flai serve` already takes. Sources: [FTS5](https://sqlite.org/fts5.html), [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite), [ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3), read 2026-10-16.

#### sqlite-vec and SQLite vec1

sqlite-vec is a small C extension that stores float, int8, and binary vectors in `vec0` virtual tables; it is pre-v1 and warns of breaking changes. Its Go bindings offer a cgo build, which fails the Go fit bar, and a Wasm build for ncruces last tagged on 2024-11-20, before ncruces moved to wasm2go; whether it works with ncruces today is not confirmed. vec1 is SQLite's own vector extension, at version 0.7, approximate (IVFADC with OPQ), and its page says "testing is insufficient"; ncruces ships it as `ext/vec1`, cgo-free. Either would add a vector column beside FTS5 in the same file, so one query could fuse BM25 and similarity; both need a local embedding model, and a model change re-embeds every section. Sources: [sqlite-vec](https://github.com/asg017/sqlite-vec), [sqlite-vec-go-bindings](https://github.com/asg017/sqlite-vec-go-bindings), [vec1](https://sqlite.org/vec1/doc/trunk/doc/vec1.md), [ncruces extensions](https://github.com/ncruces/go-sqlite3/blob/main/ext/README.md), read 2026-10-16.

#### LanceDB

An embedded database on the Lance columnar format with vector search, BM25 full-text search, and hybrid search fused by RRF by default. Its Go SDK, at v0.1.2, links prebuilt native libraries through cgo, which fails the bar, and a local server behind its REST API was not confirmed. flai could only reach it through a Python or Node process of its own, which brings nothing an ncruces SQLite file does not. Sources: [lancedb/lancedb](https://github.com/lancedb/lancedb), [lancedb/lancedb-go](https://github.com/lancedb/lancedb-go), [hybrid search](https://docs.lancedb.com/search/hybrid-search), read 2026-10-16.

#### Chroma and Qdrant

Both are vector databases that run as servers with Go clients, so either fits only behind an optional adapter, with today's search as the fallback. Chroma embeds documents itself with all-MiniLM-L6-v2, run locally, and its README names hybrid and full-text search only for Chroma Cloud. Qdrant holds dense and sparse vectors and fuses them by RRF or DBSF; Qdrant Edge runs in-process, but only from Python or Rust. Either would hold the vectors of `doc_search`, the pack's ranking, and the dashboard, rebuilt from the files, at the cost of a second process the operator installs. Sources: [chroma-core/chroma](https://github.com/chroma-core/chroma), [Chroma embedding functions](https://docs.trychroma.com/docs/embeddings/embedding-functions), [qdrant/qdrant](https://github.com/qdrant/qdrant), read 2026-10-16.

#### Pure-Go vector indexes: chromem-go, kelindar/search, coder/hnsw

chromem-go is an embeddable vector database with no third-party dependencies, exact cosine search, and gob files for persistence; it embeds through Ollama or LocalAI locally and OpenAI by default. Its README gives 0.3 ms for 1,000 documents and 40 ms for 100,000; `doc_search` reads 2,841 headings, so exact search is enough at this size. Its MPL-2.0 is file-level copyleft, neither permissive nor GPL, so linking it is the operator's call. kelindar/search is an exact cosine index under MIT with an optional llama.cpp module, but it needs Go 1.27 and flai's `go.mod` says 1.26.0. coder/hnsw is an approximate HNSW graph under CC0, which this corpus does not need and which gives up determinism. Each would add vectors beside flai's BM25, fused in Go; what any of them saves over a slice of vectors and a dot product is persistence and filters. Sources: [chromem-go](https://github.com/philippgille/chromem-go), [kelindar/search](https://github.com/kelindar/search), [coder/hnsw](https://github.com/coder/hnsw), read 2026-10-16.

### Text search engines and libraries

These improve the lexical side: better tokenizing, phrases, stemming, and a persisted index. Alone, none closes the gap of [Where recall fails](#where-recall-fails) that comes from words differing from meaning.

#### Bleve

A Go indexing and search library with BM25 or tf-idf scoring, fuzzy, phrase, and synonym queries, highlighting, and hybrid search fused by RRF or RSF. Its vector search embeds FAISS, a C++ library, behind a `vectors` build tag, which fails the Go fit bar; its text search does not need the tag, but a cgo-free build for every target was not confirmed. It could replace `flai/internal/search` and persist the dashboard's index, with synonyms as a cheap answer to "concurrently" against "WIP". Keeping it in step means indexing and deleting per changed file. Sources: [blevesearch/bleve](https://github.com/blevesearch/bleve), [vectors](https://github.com/blevesearch/bleve/blob/master/docs/vectors.md), read 2026-10-16.

#### Tantivy

A Rust library in Lucene's mould: BM25, phrase queries, stemming for 17 languages, incremental indexing, and a start-up under 10 ms that its README calls fit for command-line tools. flai cannot link it; `tantivy-cli` could run as a separate process, but it brings only BM25, which flai has. It is design to compare against, not a candidate to adopt. Source: [quickwit-oss/tantivy](https://github.com/quickwit-oss/tantivy), read 2026-10-16.

#### Zoekt

A Go code search engine that indexes trigrams for substring and regexp matching and ranks by code signals such as symbols; BM25 scoring is an experimental option, and `DebugScore` says why a result ranked. It is built for source code across repositories, not for prose by meaning. Its commands and JSON API could serve exact-phrase search over the archive and the threads, but that is a narrow gain for a new process. Sources: [sourcegraph/zoekt](https://github.com/sourcegraph/zoekt) and [api.go](https://github.com/sourcegraph/zoekt/blob/main/api.go), read 2026-10-16.

#### txtai

A Python framework whose embeddings database unites dense and sparse indexes, a graph, and a relational store. One setting, `hybrid`, adds a BM25 index beside the dense one; its default embedding model, all-MiniLM-L6-v2, runs locally, and Faiss is its default vector index. It runs as an API server with a Go client, so it fits behind an optional adapter with no account and no key. It is the nearest to a hybrid ready to run as it ships, at the cost of a Python process the operator installs and keeps in step by re-indexing changed sections. Sources: [neuml/txtai](https://github.com/neuml/txtai), [general configuration](https://neuml.github.io/txtai/embeddings/configuration/general/), [ANN configuration](https://neuml.github.io/txtai/embeddings/configuration/ann/), read 2026-10-16.

### Local embedding runtimes

Every vector or hybrid design needs a model that turns a section and a query into vectors on the operator's host. Keeping vectors in step costs one embedding per changed section, keyed by a hash of its text; a model change re-embeds the corpus, and vectors from two models never mix.

#### Ollama

A local model server with an `/api/embed` endpoint that takes one input or a batch and returns L2-normalized vectors. Its library lists embedding models from all-minilm (46 MB, 512-token context) through nomic-embed-text (274 MB, 2K context) and embeddinggemma (622 MB, 2K context) to qwen3-embedding (639 MB at 0.6b, 4.7 GB at 8b). It is a separate process, so flai would call it behind an optional adapter, set in the config, and fall back to BM25 alone when it is off or down. chromem-go, Mem0, and Graphiti all reach it already. Sources: [ollama/ollama](https://github.com/ollama/ollama), [embeddings](https://docs.ollama.com/capabilities/embeddings), and the model pages under [ollama.com/library](https://ollama.com/library), read 2026-10-16.

#### hugot

A Go library that runs Hugging Face ONNX pipelines, feature extraction among them, with a pure-Go backend from GoMLX that needs no cgo and ONNX Runtime and XLA backends that do. Its README says the Go backend suits small models such as all-MiniLM-L6-v2 and batches of about 32. It is the one way found to embed inside flai's binary with cgo off, so a query needs no second process; the cost is the model files on disk and a slower backend, and its targets beyond linux are not confirmed. Source: [knights-analytics/hugot](https://github.com/knights-analytics/hugot), read 2026-10-16.

### For T-1798 to try

Three, all local, free, and without an account, each against flai's own BM25 as the baseline on the queries of [Where recall fails](#where-recall-fails):

1. **flai's BM25 fused with exact cosine search, embeddings from Ollama.** Hybrid by reciprocal rank fusion in Go, vectors in memory or in a file under `.flai-cache`, the model behind an optional adapter. It changes least, keeps today's search as the fallback, and stays deterministic for a fixed model. Try all-minilm and nomic-embed-text, and hugot's Go backend as the in-binary alternative to Ollama.
2. **SQLite FTS5 and vec1 through `ncruces/go-sqlite3`.** Cgo-free, one file holding BM25 and vectors, a persisted index, and the shape S-0427's synchronized store could share. It tests whether a store improves latency and refresh, not only recall; vec1's immaturity is the risk to measure.
3. **txtai with `hybrid` on, as a separate process.** The ready-made hybrid, to say whether adopting behind an adapter beats adapting the design.

The agent memory systems are not on the list: each writes a store that is not derived from the files and needs an LLM call per write, which fails Files as source and Determinism. Their designs go to the state of the art that follows.

## State of the art

These are approaches from papers and products that flai cannot link as a library. They bear on recall, on search time, or on what an agent pays to hold what it found. Each was read on 2026-10-16 from its paper's abstract page or full text, or from the vendor's post; the figures are the source's own, on its own benchmark, and none was measured on this corpus. One question runs through them all: where does the model run? An approach that needs a model at query time adds a call to every search or every prime. An approach that needs one only at index time can write its output into a cache derived from the files and keep each query as cheap as today's BM25. An approach that needs none can go into flai's Go as it stands. Counts about this repository were taken on 2026-10-16 in the S-0426 worktree with `ls` and `grep`.

| Approach | Bears on | Model at query time | Model at index time | Open source candidate above | Expected effect on flai |
|----------|----------|---------------------|---------------------|-----------------------------|-------------------------|
| Hybrid retrieval fused by RRF | `doc_search`, the pack's ranking, the dashboard | An embedding of the query | An embedding per section | [LanceDB](#lancedb), [Qdrant](#chroma-and-qdrant), [Bleve](#bleve), [txtai](#txtai) | Recall across differing words; a local model call per query |
| Learned sparse (SPLADE) | `doc_search`, the pack's ranking | An encoder, or none for inference-free variants | An encoder per section | [Qdrant](#chroma-and-qdrant) stores sparse vectors; none runs SPLADE | Expansion terms in an inverted index; inference-free queries near BM25's latency |
| Late interaction (ColBERT) | `doc_search` | An encoder per query | A vector per token | None | Strong ranking; a store of token vectors too heavy for the gain here |
| Contextual retrieval | `doc_search`, the pack's ranking, the digest | None on the BM25 side | An LLM per chunk | None; flai's title and heading boosts are the model-free half | Index-time cost only; recall gained on the lexical side |
| Cross-encoder reranking | `doc_search` | A reranker over the top hits | None | None | Precision; a model call per query |
| Query and document expansion | `doc_search`, the dashboard | An LLM for HyDE; none for a glossary | A model for doc2query; none for a glossary | [Bleve](#bleve) synonyms | A glossary closes "concurrently" to "WIP" at no query cost |
| Static embeddings | `doc_search`, the dashboard | A table lookup and an average | A pass per section | None | A dense signal at close to lexical cost, inside the binary |
| Hierarchical summaries (RAPTOR) | The pack's briefs, the digest | An embedding of the query | An LLM summary per cluster | None | Answers that span sections; summaries to keep current |
| Graph RAG (GraphRAG, LazyGraphRAG) | The digest, the archive | An LLM | An LLM (GraphRAG); NLP only (LazyGraphRAG) | [Cognee](#cognee), [Graphiti](#graphiti-zep) | Questions over the whole archive; costly either side |
| Personalized PageRank (HippoRAG) | The pack, `doc_search`, the digest | An LLM to name the query's entities; PageRank needs none | An LLM extracts triples | None | Multi-hop recall by one graph walk |
| flai's own link graph | The pack, `doc_search`, the digest, S-0427's store | None | None | flai's `Linked` and `Step` | Model-free multi-hop over links that exist |
| Tiered memory and paging (MemGPT) | The pack in parts, `doc_get` | The agent's own model pages | None | [Letta Code](#letta-code-memgpt) | Already the pattern: a core in context, the rest by tool |
| Self-organising notes (A-MEM) | The digest | An embedding | An LLM per note | [Mem0](#mem0), [LangMem](#langmem) in part | Fewer tokens per answer; notes rewritten by a model |
| Reflection and consolidation | The digest's keeping run | None to read | An LLM per run | [LangMem](#langmem)'s manager, [Letta Code](#letta-code-memgpt)'s dreaming | Model cost moved off the query path |
| Temporal facts, decay, compaction | The digest, the pack | None | An LLM or a rule per run | [Graphiti](#graphiti-zep) | A smaller primed set that stays true |
| Agentic keyword search | `doc_search`, the harness | The agent's own model | None | flai's `doc_search` as it is | The cheapest index; it costs turns |
| Long context against retrieval | The pack | Not applicable | Not applicable | Not applicable | Evidence for selecting, not loading |
| Prompt caching and a stable pack | The pack, role packs | None beyond the paid API the agent already calls | None | The harness | Each byte cut from the pack is cut from every later call |
| Incremental indexes | `doc_search`, the pack, the dashboard, S-0427's store | None | None; embeddings cached by text hash | [SQLite FTS5](#sqlite-fts5), [Bleve](#bleve), [Tantivy](#tantivy) | No rebuild per call |
| Sections bounded in size | `doc_search`, `doc_get`, the pack's ranking | None | None | None | `flai-cli.md § Commands` stops winning every query |

### Retrieval

#### Hybrid retrieval fused by reciprocal rank

A lexical ranking and a dense one are run apart and merged by reciprocal rank fusion: each document scores the sum of 1/(k + rank) over the lists, with k = 60. RRF needs only ranks, so it needs no score calibration between BM25 and cosine. Cormack, Clarke, and Büttcher found it beat every input system and Condorcet Fuse on TREC runs and LETOR 3 ([Reciprocal rank fusion outperforms Condorcet and individual rank learning methods](https://plg.uwaterloo.ca/~gvcormac/cormacksigir09-rrf.pdf), SIGIR 2009). The fusion is a few lines of Go and deterministic; the dense list needs an embedding of each query, so a model on the query path, and a vector per section kept in step. [LanceDB](#lancedb), [Qdrant](#chroma-and-qdrant), [Bleve](#bleve), and [txtai](#txtai) implement it, and it is pick 1 of [For T-1798 to try](#for-t-1798-to-try).

#### Learned sparse retrieval (SPLADE)

SPLADE has a transformer write each document as weighted vocabulary terms, adding terms the text does not contain, so the result lives in an ordinary inverted index. Formal, Piwowarski, and Clinchant report "highly sparse representations and competitive results" against dense and sparse methods ([SPLADE](https://arxiv.org/abs/2107.05720), SIGIR 2021, 2021-07-12); SPLADE v2 reports "more than 9% gains on NDCG@10 on TREC DL 2019" ([SPLADE v2](https://arxiv.org/abs/2109.10086), 2021-09-21). Inference-free variants encode only documents and look the query's tokens up: Geng, Wang, Ru, and Yang report client-side latency "only 1.1x that of BM25" and a 3.3 NDCG@10 gain on BEIR over the previous best inference-free model ([Towards Competitive Search Relevance For Inference-Free Learned Sparse Retrievers](https://arxiv.org/abs/2411.04403), 2024-11-07). That shape suits flai: a model at index time writes expansion terms into a cache, and the query stays in flai's BM25 code. The encoder is a Hugging Face model, which [hugot](#hugot) could run at index time; that was not tried. [Qdrant](#chroma-and-qdrant) and [txtai](#txtai) hold sparse vectors; none of the candidates was checked running SPLADE.

#### Late interaction (ColBERT)

ColBERT keeps a vector per token of every passage and scores a query by matching each query token to its nearest passage token. Khattab and Zaharia report it runs "two orders-of-magnitude faster" with "four orders-of-magnitude fewer FLOPs per query" than BERT rerankers of the time ([ColBERT](https://arxiv.org/abs/2004.12832), SIGIR 2020, 2020-04-27); ColBERTv2 cuts the space footprint "by 6--10×" ([ColBERTv2](https://arxiv.org/abs/2112.01488), NAACL 2022, 2021-12-02). It needs an encoder for every query and a store of token vectors many times the corpus's 3.7 MB. No candidate above implements it; the reference is [stanford-futuredata/ColBERT](https://github.com/stanford-futuredata/ColBERT), in Python. For a corpus of 2,841 headings, it adds weight that a reranker or plain hybrid does not.

#### Contextual retrieval

Before indexing, an LLM writes for each chunk a short note, "usually 50-100 tokens", on where the chunk sits in its document, and the note is indexed with the chunk by both BM25 and embeddings. Anthropic report top-20 retrieval failures fell from 5.7% to 3.7% (35%) with contextual embeddings, to 2.9% (49%) with contextual BM25 added, and to 1.9% (67%) with a reranker on top, at a one-time cost of $1.02 per million document tokens with prompt caching ([Introducing Contextual Retrieval](https://www.anthropic.com/news/contextual-retrieval), Anthropic, 2024-09-19). [agent-context.md](agent-context.md#how-others-choose) cites the first two figures. The model runs only at index time, and the note could be a file the E-0025 keeping run writes, so the query stays BM25. flai already does the model-free half: `IndexSections` indexes each section's document title and heading path as fields boosted 3 and 2 against the body's 1 (`flai/internal/search/sections.go`, `search.go`). No candidate implements it.

#### Reranking with a cross-encoder

A cross-encoder reads the query and one passage together and scores the pair; it is run over the top hits of a cheaper first stage. Nogueira and Cho report a 27% relative gain in MRR@10 on MS MARCO over the previous best ([Passage Re-ranking with BERT](https://arxiv.org/abs/1901.04085), 2019-01-13), and the reranker is what took Anthropic's figures above from 49% to 67%. It needs a model at query time, run once per hit reranked, and Anthropic's tests used Cohere's reranker, a hosted service. A local cross-encoder could run under [hugot](#hugot), at a cost per query not measured here. No candidate above was checked for one.

#### Query and document expansion

HyDE has an LLM write a hypothetical answer and embeds that in place of the query; Gao, Ma, Lin, and Callan report it "significantly outperforms" the unsupervised dense retriever Contriever ([Precise Zero-Shot Dense Retrieval without Relevance Labels](https://arxiv.org/abs/2212.10496), 2022-12-20). It puts an LLM call in front of every search. doc2query runs the other way: a model predicts the questions a document answers and appends them before indexing, and Nogueira, Yang, Lin, and Cho report that, without reranking, this approaches neural rerankers' effectiveness while being "much faster" ([Document Expansion by Query Prediction](https://arxiv.org/abs/1904.08375), 2019-04-17). The model-free form is a project glossary of synonyms applied at query or index time, such as "concurrently" and "concurrency cap" to "WIP limit", the miss of [Where recall fails](#where-recall-fails); [Bleve](#bleve) ships synonyms. Pseudo-relevance feedback, which adds terms from the top hits and searches again, is a second model-free form; no figure was taken for it.

#### Static embeddings

A static embedding model replaces the transformer with a table: each token's vector is looked up and the vectors are averaged. Aarsen reports static-retrieval-mrl-en-v1 at 107,420 sentences a second on CPU against all-mpnet-base-v2's 270 (397 times faster) and NanoBEIR NDCG@10 of 0.5032 against 0.5757, "87.4%" of the transformer's ([Train 400x faster Static Embedding Models with Sentence Transformers](https://huggingface.co/blog/static-embeddings), Hugging Face, 2025-01-15). Inference is a lookup and an average, which Go can do with no runtime and no cgo, given the model's tokenizer and its table on disk; their size was not checked. It gives the dense list of hybrid retrieval without a second process or a transformer on the query path. No candidate above implements it; [hugot](#hugot) and [Ollama](#ollama) run transformers.

### Structure over the corpus

#### Hierarchical summaries (RAPTOR)

RAPTOR embeds chunks, clusters them, has an LLM summarize each cluster, and repeats up the tree, then retrieves from any level. Sarthi, Abdullah, Tuli, Khanna, Goldie, and Manning report that with GPT-4 it improves "the best performance on the QuALITY benchmark by 20% in absolute accuracy" ([RAPTOR](https://arxiv.org/abs/2401.18059), 2024-01-31). It needs an LLM per cluster at index time and an embedding at query time, and a changed chunk changes the summaries above it. flai has a tree written by people already: the document, its headings, and a brief that holds the first paragraph and the outline ([ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)). The idea that carries over is summaries at more than one level for the archive, where nobody writes them: an epic's stories, a quarter's decisions. That is E-0025's digest. No candidate above implements it.

#### Graph RAG (GraphRAG, LazyGraphRAG)

GraphRAG has an LLM extract entities and relations from every chunk, groups them into communities, and writes a summary of each, then answers a question about the whole corpus by map and reduce over the summaries. Edge and others report "substantial improvements over a conventional RAG baseline" in comprehensiveness and diversity on datasets "in the 1 million token range"; the abstract gives no figure ([From Local to Global](https://arxiv.org/abs/2404.16130), Microsoft, 2024-04-24). LazyGraphRAG defers every LLM call to query time: its indexing costs "0.1% of the costs of full GraphRAG", and it matches GraphRAG's global search at over 700 times lower query cost ([LazyGraphRAG](https://www.microsoft.com/en-us/research/blog/lazygraphrag-setting-a-new-standard-for-quality-and-cost/), Microsoft Research, 2024-11-25). Either puts an LLM on the path, at index time, query time, or both, which fails [Local and free](#what-a-solution-must-meet) unless the model runs locally. [Cognee](#cognee) and [Graphiti](#graphiti-zep) build such graphs. For flai the use is the archive's global questions, which agents rarely ask; the pack's questions are about one story.

#### Personalized PageRank over a graph (HippoRAG)

HippoRAG has an LLM extract triples from every passage into a graph, links similar phrases with a retrieval encoder, and at query time has an LLM name the query's entities, then runs Personalized PageRank from them and ranks passages by the probability their nodes gather. Gutiérrez, Shu, Gu, Yasunaga, and Su report gains "by up to 20%" over the state of the art on multi-hop question answering, and single-step retrieval "10-30 times cheaper and 6-13 times faster" than iterative retrieval with IRCoT ([HippoRAG](https://arxiv.org/abs/2405.14831), NeurIPS 2024, 2024-05-23). HippoRAG 2 reports "a 7% improvement in associative memory tasks over the state-of-the-art embedding model" ([From RAG to Memory](https://arxiv.org/abs/2502.14802), 2025-02-20). The LLM's work is to build the graph and to find where to start; PageRank needs no model. flai has a graph and its seeds already, as the next section says, so the part to take is the walk. No candidate above implements it; Aider ranks a repository map by the same personalized PageRank with no model ([agent-context.md](agent-context.md#how-others-choose)).

#### The graph flai already has

flai's documents link each other in ways a person wrote, so the graph needs no extraction. The pack follows it: `Linked` takes every document the story, its epic, and its tasks link, and `Step` follows one step on from what was chosen: an ADR's `refines` both ways and the ADR links in chosen sections; `chooseADR` swaps a superseded ADR for what supersedes it (`flai/internal/context/design.go`). On 2026-10-16, 123 of the 158 ADRs refine another, 193 edges in all, and 8 are superseded; `design/system` holds 734 links to ADRs; 203 of the 391 archived stories name an ADR, 512 times, and between them name 266 distinct threads and issues; 55 ADRs and 27 of the 28 files in `design/system` carry topics. What flai does not do: walk more than one step, rank by where the walk concentrates, walk back from an ADR to the stories, threads, and issues that cite it, or let `doc_search` use the graph at all. Personalized PageRank from the story's seeds (its links, its topics' documents, its touches) over these edges would be HippoRAG's retrieval with no model, deterministic, and each result could say which edges brought it. The same edges, extended to the archive and the threads, are the index E-0025's digest and S-0427's store would hold.

### Agent memory

#### Tiered memory and paging (MemGPT)

MemGPT treats the context window as main memory and an external store as disk, and has the model page between them with function calls. Packer, Wooders, Lin, Fang, Patil, Stoica, and Gonzalez report deep memory retrieval accuracy of 92.5% with GPT-4 under MemGPT against 32.1% without, and 66.9% against 38.7% with GPT-3.5 Turbo ([MemGPT](https://arxiv.org/abs/2310.08560), 2023-10-12, Table 2 of v2). The paging decision is the agent's own model, so it costs turns, not a second model. flai's design is already of this kind: a core pushed by `prime`, the rest pulled by `doc_search` and `doc_get`, and the pack itself paged in parts ([ADR-0104](../adrs/0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)). What it adds is a core that the agent rewrites, which [Letta Code](#letta-code-memgpt) keeps as files in git; in flai the narrative's `## Current state` plays that part.

#### Self-organising notes (A-MEM)

A-MEM writes each new memory as a note with keywords, tags, and a context line, links it to related notes, and lets a new note revise the old ones, after the Zettelkasten method. Xu, Liang, Mei, Gao, Tan, and Zhang report on LoCoMo with GPT-4o-mini a multi-hop F1 of 27.02 against MemGPT's 26.65, at 2,520 tokens against MemGPT's 16,977, and retrieval of 3.70 µs at a million notes ([A-MEM](https://arxiv.org/abs/2502.12110), NeurIPS 2025, 2025-02-17). Its text claims "at least two times better" multi-hop performance and "approximately 1,200 tokens per memory operation", which its Table 1 does not show; the table's figures are given here. Every note costs an LLM call when written, and notes change after the fact, which fails [Determinism](#what-a-solution-must-meet) unless the notes are files under review. For E-0025, the shape of a note (a fact, its keywords, its links, its source) is the part to take. [Mem0](#mem0) and [LangMem](#langmem) write LLM-chosen memories in the same spirit; neither was checked for A-MEM's linking.

#### Reflection and consolidation

Generative Agents keep every observation in a memory stream, retrieve by the sum of recency (decay of 0.995 per game hour since last retrieved), importance (1 to 10, rated by the model once), and relevance (cosine), and write higher-level reflections when the importance of recent events passes 150. Park, O'Brien, Cai, Morris, Liang, and Bernstein report, in their ablation, TrueSkill ratings of 29.89 for the full architecture, 26.88 without reflection, and 21.21 with observation, reflection, and planning all removed ([Generative Agents](https://arxiv.org/abs/2304.03442), UIST 2023, 2023-04-07). Sleep-time compute moves the thinking between queries: Lin, Snell, Wang, Packer, Wooders, Stoica, and Gonzalez report it cuts the test-time compute for the same accuracy "by ~ 5x" and the average cost per query "by 2.5x" when queries share a context ([Sleep-time Compute](https://arxiv.org/abs/2504.13171), 2025-04-17). The lesson for flai is where the model runs: E-0025's keeping run after an acceptance is sleep-time compute, and what it writes is read at no model cost when a story primes. [LangMem](#langmem)'s background manager and [Letta Code](#letta-code-memgpt)'s dreaming implement consolidation.

#### Temporal facts, decay, and compaction

Zep's graph records when each fact held and invalidates a fact a later episode contradicts. Rasmussen, Paliychuk, Beauvais, Ryan, and Chalef report 94.8% on deep memory retrieval against MemGPT's 93.4%, accuracy gains "up to 18.5%" on LongMemEval, and 90% lower response latency than their baselines ([Zep](https://arxiv.org/abs/2501.13956), 2025-01-20). Mem0's paper reports a 26% relative gain over OpenAI's memory on LoCoMo, 91% lower p95 latency, and over 90% fewer tokens ([Mem0](https://arxiv.org/abs/2504.19413), 2025-04-28); those are the managed platform's, as [Mem0](#mem0) notes. Anthropic's guidance names compaction, summarizing a context near its limit and starting again from the summary, and structured note-taking outside the window ([Effective context engineering for AI agents](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents), Anthropic, 2025-09-29). E-0025 proposes the file-based form: an entry is marked absorbed when a design document, a convention, or an ADR records it, and compacted to one line when no story has needed it within a window; history keeps everything. Absorption is invalidation by the design, as Zep invalidates by a later episode, and it can be checked without a model when the entry cites the document that absorbed it. [Graphiti](#graphiti-zep) implements temporal invalidation.

### Search by the agent itself

#### Agentic keyword search against retrieval

An agent with a keyword or grep tool searches, reads, and searches again, with no index beyond the files. Sen, Kasturi, Lumer, Gulati, and Subbiah compare grep and vector retrieval on 116 LongMemEval questions in their own harness and in Claude Code, Codex, and Gemini CLI, and find that "grep generally yields higher accuracy than vector retrieval", while overall scores "depend strongly on which harness and tool-calling style is used"; the abstract gives no figure ([Is Grep All You Need?](https://arxiv.org/abs/2605.15184), 2026-05-14). Subramanian and others report agents with only keyword search tools reach "over 90% of the performance metrics" of RAG systems without a vector database ([Keyword search is all you need](https://arxiv.org/abs/2602.23368), 2025-12-19). Anthropic describe Claude Code as hybrid: `CLAUDE.md` loaded up front, glob and grep just in time ([Effective context engineering for AI agents](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents)). This is what flai's agents do today with `doc_search` and the harness's own reads. It costs turns, not a model, and its speed is the speed of each tool call, so a faster `doc_search` with smaller results helps it as much as a better ranking does.

#### Long context against retrieval

Li, Li, Zhang, Mei, and Bendersky find that "when resourced sufficiently, LC consistently outperforms RAG in terms of average performance", while RAG keeps a "significantly lower cost"; their Self-Route lets the model choose per query and "significantly reduces the computation cost while maintaining a comparable performance to LC" ([Retrieval Augmented Generation or Long-Context LLMs?](https://arxiv.org/abs/2407.16833), EMNLP 2024 industry track, 2024-07-23). [agent-context.md](agent-context.md#what-the-evidence-says) records the other side: models use the middle of a long context worst, and performance falls as input grows. For flai the pack is the long-context half and `doc_search` the retrieval half. The evidence argues for a pack that holds what the story needs, small enough that its cost on every turn stays low, and retrieval for the rest, which is ADR-0049's design; it does not argue for loading more.

### Cost

#### Prompt caching and a stable pack prefix

The Messages API caches a prompt's prefix in the order tools, system, messages, and a change at one level invalidates it and every level after. A 5-minute cache write costs 1.25 times the base input price, a 1-hour write 2 times, and a read 0.1 times, 0.05 times on Claude Opus 5.5 and Sonnet 5.5 ([Prompt caching](https://platform.claude.com/docs/en/build-with-claude/prompt-caching), Anthropic documentation, read 2026-10-16). The launch post claimed costs cut "by up to 90%" and latency "by up to 85%" for long prompts, with a 100,000-token cached book answering in 2.4 s against 11.5 s ([Prompt caching with Claude](https://claude.com/blog/prompt-caching), Anthropic, 2024-08-14). A pack, once returned by `prime`, sits in the messages and is read from the cache on every later call ([What is not measured](#what-is-not-measured)). By the 4 bytes a token that [agent-context.md](agent-context.md#fitting-the-pack-to-a-budget)'s 541 KB and 135k tokens imply, S-0426's 128,182-byte pack is about 32,000 tokens, more as the 152,383 bytes of JSON the MCP tool returns: a sixth to over a third of the 87k to 188k cache-read tokens per call recorded there. That is an estimate, not a measure. Two things follow: every byte cut from the pack is cut from every later call, and a pack that is the same bytes for the same inputs lets a resumed session or a second agent with the same prefix read it from the cache instead of writing it again. Whether today's pack puts anything that varies, such as a story ID or a count, before its shared parts was not checked. No model is involved beyond the paid API the agent already calls; flai does not set cache breakpoints, the harness does.

#### Incremental indexes

An incremental index keeps its state between queries and takes in only what changed. Cursor hashes every file into a Merkle tree, syncs only the branches whose hashes differ, and caches embeddings by chunk content; reusing a teammate's index cut median time-to-first-query from 7.87 s to 525 ms, and at the 99th percentile from 4.03 hours to 21 s ([Securely indexing large codebases](https://cursor.com/blog/secure-codebase-indexing), Jeremy Stribling, Cursor, 2026-01-27). flai has the parts: `flai serve` hashes every file's path, modification time, and size to decide whether to rebuild the dashboard's index, and `flai mcp` keeps parsed work items between requests ([Search](#search)). `doc_search` and `prime` still read and cut every document on every call. A per-file hash, kept in `flai mcp` and `flai serve`, would let them re-cut only changed files, and an embedding or an expansion keyed by a hash of a section's text survives every file change that leaves the section alone. No model is needed. [SQLite FTS5](#sqlite-fts5), [Bleve](#bleve), and [Tantivy](#tantivy) index incrementally, and pick 2 of [For T-1798 to try](#for-t-1798-to-try) persists the index to disk, which is S-0427's store.

#### Sections bounded in size

How a corpus is cut decides what a hit costs to read. Smith and Troynikov found chunk size and overlap moved precision and token efficiency far more than recall: recall stayed between about 84 and 92 across most settings, while the recursive splitter's IoU rose from 1.5 at 800 tokens to 6.9 at 200 with no overlap, and OpenAI's Assistants default of 800 tokens with 400 overlapping scored IoU 1.4 ([Evaluating Chunking Strategies for Retrieval](https://www.trychroma.com/research/evaluating-chunking), Chroma, 2024-07-03). Late chunking embeds a whole document with a long-context model and pools each chunk afterwards, so a chunk's vector carries its neighbours' context with no LLM ([Late Chunking](https://arxiv.org/abs/2409.04701), Günther, Mohr, Williams, Wang, and Xiao, Jina AI, 2024-09-07). flai cuts only at headings, so `flai-cli.md § Commands` is one section of 340,847 bytes that matches most queries and no `doc_get` result can carry ([Where recall fails](#where-recall-fails)). Splitting a section over a bound at paragraphs, list items, or table rows, each piece indexed with its title and heading path as today, needs no model. No candidate above implements it; each indexes what it is given.

### What this means for flai

For T-1798 and T-1799 to weigh, the approaches that need no model at query time and bear most on flai's overhead or recall:

- **Sections bounded in size.** No model at all. It removes the one section that wins every broad query and makes every hit readable through `doc_get`.
- **An incremental index kept in `flai mcp` and `flai serve`.** No model. It removes the rebuild on every `doc_search` and `prime`, the cost [Today](#today) found in every call, and is the shape S-0427's store would persist.
- **A smaller, byte-stable pack.** No model. Its size is paid on every turn of every story from the cache, and its stability decides whether the cache can be shared.
- **flai's own link graph, walked by personalized PageRank.** No model. HippoRAG's retrieval and Aider's ranking without the extraction; it reaches the archive, the threads, and the issues if the walk follows story to ADR to thread, which is where E-0025 points.
- **Models at index time only.** Contextual notes, doc2query-style expansion, inference-free sparse terms, and E-0025's consolidation all run once per changed section or per keeping run, write a derived cache or a file, and leave each query to BM25 and Go.
- **Static embeddings fused with BM25 by RRF.** A table lookup at query time, not a model call: the dense half of hybrid retrieval inside the binary, if a model's table and tokenizer can ship or be fetched.
- **A project glossary of synonyms.** No model. It closes the misses where the operator's words and the document's differ, at the cost of a list someone keeps.

Dense retrieval through a transformer, SPLADE's full form, ColBERT, cross-encoder reranking, HyDE, GraphRAG, and HippoRAG's query step each put a model on the query path; what they would add over the list above is for T-1798's trial to show.

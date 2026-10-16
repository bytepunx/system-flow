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

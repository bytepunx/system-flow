---
title: Priming an agent with the documentation its story needs
updated: 2026-09-26
status: active
---

# Priming an agent with the documentation its story needs

The finding of S-0125, for E-0010. Every agent flai starts reads every convention, whatever its story, and then finds the design and the decisions it needs by itself, or does not. This document says what an agent loads today, what a story already says about what it will need, how agent tools and research choose context, and ranks the designs that fit flai. The decision and the stories that follow are recorded at the end once the designer has chosen.

## What an agent loads today

| Step | What | Size |
|------|------|------|
| Harness prompt | `harness.Prompt` in `flai/internal/harness/harness.go`: work the story, prime with `flai prime --cat` | one paragraph |
| `CLAUDE.md` | The map: layout, precedence, what to read | 6 KB |
| `flai prime --cat` | `design/conventions/README.md` and all twelve convention files in order (`flai/cmd/prime.go`, `internal/conventions`), then the open-issues table | 47 KB and 9 KB |
| Everything else | Whatever the agent chooses to open: `design/system`, `design/adrs`, `design/tech`, `docs/` | unmeasured |

Sizes of what the agent could need, in this repository on 2026-09-26 (a token is about four bytes):

| Corpus | Files | Size | Notes |
|--------|-------|------|-------|
| `design/conventions` | 13 | 47 KB | Loaded whole every session. `logging.md` and `telemetry.md` (7 KB) apply only to code that runs; `code-quality.md` (5 KB) only to code. |
| `design/adrs` | 47 | 179 KB | 8 superseded, 22 refine another. Immutable once accepted. |
| `design/system` | 19 | 283 KB | `flai-cli.md` is 71 KB, of which 63 KB is one `## Commands` section; `flaiover-dashboard.md` is 69 KB in 15 sections. |
| `design/tech` | 13 | 25 KB | Each has a `Where` row naming the part that uses it. |
| `design/issues/summary.md` | 1 | 9 KB | Every open issue, relevant or not. |

Two things are wrong with this. What is always loaded is not all needed: a documentation story reads the telemetry rules, a dashboard story reads every issue about releases. And what is needed is not loaded: the decisions and the design a story builds on are left to the agent to find, which costs turns and tokens when it looks and mistakes when it does not. Loading everything is not the answer either: `design/` is about 540 KB, some 135 thousand tokens, which would crowd out the work, and models use what sits in the middle of a long context worse than what sits at its ends.

## What a story already says

A story carries signals flai can read without a model. Measured on the 126 archived stories:

| Signal | Where | Coverage | What it points at |
|--------|-------|----------|-------------------|
| `tags` | Front matter | 106 of 126 non-empty | Manifest projects (`cli` → `flai`, `dashboard` → `flaiover`, `template`), so their paths, their design document, and the `design/tech` files whose `Where` names them |
| `touches` | Front matter, story and tasks | 90 of 126 | Paths, usually coarse (`flai/cmd`, `flaiover/src`) |
| Links and IDs in the body | Goal, criteria, notes | 52 name an ADR, 55 a `design/system` path | The exact documents |
| Parent epic | Front matter | all | The epic's own tags, touches, and links |
| Text | Title and goal | all | Anything, by matching words |

The documents link to each other too: `design/system` files link 37 of the 47 ADRs, and an ADR links what it supersedes and refines.

### A replay against the archive

For each archived story, which of the ADRs that existed when it was created does its body name, and would a selection have reached them? 26 stories named 44 such ADRs.

| Selection | ADRs reached | Recall |
|-----------|--------------|--------|
| Tags → the project's design document → every ADR it links | 34 | 77% |
| Lexical ranking (BM25) of ADRs against title and goal, top 5 | 35 | 80% |
| Same, top 10 | 38 | 86% |
| Links from tags, plus lexical top 5 | 41 | 93% |

Read it with care. The sample is small; today's design documents contain links added after some of these stories; and a story naming an ADR is a proxy for needing it, not proof. It is still enough to say two things. Neither the link graph nor word matching is enough alone, and together they miss little. And the tag walk is not selective: `flai-cli.md` links 27 ADRs, so selecting by component alone loads half of them. Precision has to come from sections, not whole documents.

## How others choose

| | Mechanism | How it works | Examples | Fit for flai |
|-|-----------|--------------|----------|--------------|
| A | Always on | A file loaded into every session | `CLAUDE.md`, a root `AGENTS.md`, `.github/copilot-instructions.md`, Cursor Always Apply, Windsurf Always On, Kiro `inclusion: always` | What `flai prime` does now. Right for the few rules every story needs. |
| B | Scoped by path | A rule declares the paths it governs and loads when the agent reads or edits a matching file | Claude Code `.claude/rules` with `paths:` and subdirectory `CLAUDE.md`, nested `AGENTS.md` (the nearest wins), Copilot `applyTo`, Cursor Apply to Specific Files, Windsurf Glob, Kiro `fileMatch` | Excellent. `touches` is the declared path set, known before the agent starts. |
| C | Described, loaded on demand | Only a name and a one-line description are in context; the body is read when the agent judges it relevant | Agent Skills ("progressive disclosure", about 100 tokens a skill), Cursor Apply Intelligently, Windsurf Model Decision, Kiro `auto` | Good as the fallback for everything not selected. flai has `flai doc show` and the MCP `doc_get`. |
| D | Manual | Loaded when a person or the agent names it | Cursor Apply Manually, Windsurf Manual, Kiro `manual` | The story's own links already do this. |
| E | Lexical retrieval | Rank chunks by term statistics (BM25) against the task | Anthropic's contextual retrieval pairs it with embeddings: embeddings alone cut failed retrievals by 35%, adding BM25 by 49% | Good. flai has BM25 in `flai/internal/search`. Deterministic and offline. |
| F | Embedding retrieval | Rank chunks by vector similarity | Most RAG systems | Poor fit: a model and a dependency, results that move with the model, a store to keep. The replay leaves little for it to add. |
| G | Structural map | A ranked outline of the repository fit to a budget | Aider's repository map: tree-sitter symbols ranked by a personalised PageRank toward the files in the chat, 1,000 tokens by default | Good for excerpts: an outline of a long document's headings lets the agent ask for the section it needs. |
| H | Just in time | Keep light identifiers (paths, queries) in context and fetch at run time with tools | Anthropic's context engineering guidance | The principle behind C and G. |
| I | Learned from history | Documents read or changed together with paths in past work | Co-change mining (ROSE) | Later: suggestions from the archive once selection has run for a while. |

Path scoping in these tools fires when the agent opens a file, so the first turns run without the rule. flai knows the story's `touches` before the agent starts, so it can select up front.

### What the evidence says

- Context is a budget. Models use information at the start and end of a long context better than in the middle (Liu et al., TACL 2024), and performance falls as input grows even on simple tasks (Chroma, 18 models, 2025).
- Instruction files are not free. Gloaguen et al. (arXiv 2602.11988) found repository context files did not generally raise task success and raised inference cost by over 20%; developer-written files gained 2.4 points, not significant, and repository overviews did not help. Lulla et al. (arXiv 2601.20404) found an `AGENTS.md` associated with 29% lower median runtime and 17% fewer output tokens at comparable completion. McMillan (arXiv 2605.10039, 1,650 Claude Code sessions) found file size, position, structure, and conflicts had no measurable effect on compliance; compliance fell with the amount of work done in the session.
- Read together: what helps is specific, relevant instruction; generic overviews cost tokens without improving results. That argues for selecting, not for loading more.

Every tool that does this well combines a small always-on core, scoping by path, and on-demand reading for the rest. None of them relies on retrieval alone for rules, because a rule that is missed is a rule that is broken.

## Candidate designs

### 1. A story's context pack: core, scoped, linked, ranked, and a catalog (recommended)

`flai prime --story S-nnnn` assembles, in order and within a budget:

1. **Core.** The conventions that apply to every story, in full. A convention says whether it is core or scoped.
2. **Scoped.** Conventions, `design/tech` files, and design sections whose declared scope overlaps the story's claim: its `touches`, its open tasks' `touches`, and its tags read as manifest project paths (the same claim ADR-0046 holds). The overlap test is the prefix test flai already has.
3. **Linked.** Every document or ADR the story, its parent epic, and its tasks name, followed one step: the ADRs a selected section links, the ADRs a selected ADR refines. Superseded ADRs are replaced by what supersedes them.
4. **Ranked.** The top sections of the remaining design and ADRs by BM25 against the story's title, goal, and criteria, to fill the budget.
5. **Catalog.** One line per document not loaded, path and title, and the heading outline of any document loaded only in part, so the agent reads the rest with `flai doc show` or `doc_get` when it needs it.

Every item carries its reason (`core`, `scope: flai/cmd`, `linked from S-0125`, `rank 0.82`), in the text output and in `--json`, so the designer can see why an agent knew what it knew. Long documents are cut at `##` headings, and a section over a size limit at its paragraphs and table rows, each chunk carrying its heading path. The harness prompt and `CLAUDE.md` change from `flai prime --cat` to `flai prime --story`; plain `flai prime` keeps its current output. An MCP tool returns the same pack.

Cheap: Go, the search package flai has, front matter. Explainable. Deterministic for a given repository and story. Its weakness is declared scope that is wrong or missing; the catalog is the safety net, and `flai check` can warn on a document with no scope.

### 2. Core plus catalog only

Load the core conventions and a catalog of everything else with one-line descriptions; the agent reads what it judges relevant. The pattern of Agent Skills. Simplest, and no scope to maintain. It spends turns on every story re-deciding what to read, and it depends on the agent reading a rule it does not know it is breaking.

### 3. Retrieval only

No scope metadata: rank every convention section, ADR, and design section against the story by BM25 (or embeddings) and load the top within a budget. No upkeep, but the replay shows a quarter of needed ADRs missed at top 5, rules become probabilistic, and the result changes with the story's wording.

### 4. A context list on the story

The designer or the refining agent writes `context:` on the story, and flai loads exactly that. Precise and transparent, and a burden at refinement that will be skipped. Design 1 already loads what the story links, which is this without a new field.

## Open choices for the designer

- Where the scope of an ADR lives, given accepted ADRs are immutable: derived from the design sections that link it (no new field), or an `applies_to:` field that may be added to an accepted ADR like `superseded_by`.
- Which conventions are core. Proposed scoped: `logging.md`, `telemetry.md`, `code-quality.md`, and the open-issues table filtered to issues whose instances name the claim's paths.
- The default budget, and what happens when core and scoped alone exceed it.
- Whether the harness switches to the pack at once or after a replay test shows it reaches what archived stories needed.

## Sources

- Claude Code [memory and rules](https://code.claude.com/docs/en/memory); Anthropic [Agent Skills](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
- Cursor [rules](https://cursor.com/docs/context/rules); GitHub Copilot [repository instructions](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions); Windsurf, now Devin Desktop, [rules and memories](https://docs.devin.ai/desktop/cascade/memories); Kiro [steering](https://kiro.dev/docs/steering/); [AGENTS.md](https://agents.md/)
- Aider [repository map](https://aider.chat/docs/repomap.html) and [its design](https://aider.chat/2023/10/22/repomap.html)
- Anthropic, [Effective context engineering for AI agents](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents) (2025) and [Contextual Retrieval](https://www.anthropic.com/news/contextual-retrieval) (2024)
- Liu et al., [Lost in the Middle](https://arxiv.org/abs/2307.03172), TACL 2024; Chroma, [Context Rot](https://www.trychroma.com/research/context-rot) (2025)
- Preprints: Gloaguen et al., [Evaluating AGENTS.md](https://arxiv.org/abs/2602.11988); Lulla et al., [AGENTS.md and agent efficiency](https://arxiv.org/abs/2601.20404); McMillan, [instruction adherence in configuration files](https://arxiv.org/abs/2605.10039)
- Co-change: [Zimmermann et al., ROSE](https://thomas-zimmermann.com/publications/files/zimmermann-tse-2005.pdf)

The 2026 preprints are not peer reviewed and disagree on cost; only their direction is relied on here.

## Decision

Pending the designer's answers on S-0125's threads.

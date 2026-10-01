---
title: Priming an agent with the documentation its story needs
updated: 2026-10-01
status: active
topics: [cli, conventions]
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

[ADR-0047](../adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), 2026-09-26. On TH-0020 the designer chose design 1 with one change: conventions are selected section by section, by `topics` on the file and on headings, a name chosen so as not to overload `tags`, which decide release components. On TH-0021:

- `topics` may be set on an accepted ADR, with `flai adr topics` (not derived from the documents that link it).
- Conventions, `design/system`, and `design/tech` files carry `topics`; `[all]` is core. A convention without `topics` is `[all]`; a design or tech file without it is selected only by links or ranking.
- Every convention starts as `[all]`; the designer narrows them afterwards.
- No budget: everything selected is loaded; ranking adds a fixed five ADRs and five design sections.
- The harness prompt, `CLAUDE.md`, the template, and a new MCP `prime` tool switch to `flai prime --story` as soon as it exists.

A story's topics are its own and its epic's `topics`, plus the name, tags, and kind of the sub-projects its tags or claim reach, plus `code` when one of them is not the template.

The stories under E-0010 that build it, in order: S-0134 (topics on documents, `flai adr topics`, the check, every convention at `[all]`), S-0135 (topics on stories and epics, and a story's derived topics), S-0136 (`flai prime --story` filters conventions), S-0137 (design, tech, ADRs by topic, links, ranking, and the catalog), S-0138 (the harness, `CLAUDE.md`, the template, and an MCP `prime` tool switch to it).

## As built

S-0136 and S-0137 built the pack in `flai/internal/context`, printed by `flai prime --story` ([flai-cli.md](flai-cli.md)). Choices ADR-0047 left open:

- Sections are cut at every heading. A section chosen by topics or a `#fragment` link includes the sections below it. A ranked design section is one heading down to the next heading of any level. An ADR always loads whole.
- Links are markdown links, a document's repository path written out, and `ADR-nnnn` IDs. A link that does not resolve from where the item is now is matched by its tail, so an archived item's links still work.
- The one step is followed from every selected section, those chosen by topics included, because that is what reaches the decisions (below).
- A superseded ADR is replaced wherever it is chosen, and superseded ADRs are not ranked. The catalog marks them.
- Ranking indexes each unchosen section as one `flai/internal/search` document and drops common English words from the query. An ADR ranks by its best section.

The replay on 2026-09-29 (S-0137 notes) re-ran the archive with today's documents. 56 archived stories named 93 ADRs that existed when each was created. The counts leave out the story's own links:

| Selection | ADRs reached |
|-----------|--------------|
| Topics alone (no ADR carries topics yet) | 0 |
| Topics, then one step from the selected sections | 91 |
| The above, then ranking | 93 |
| Topics and ranking, no step | 68 |
| Ranking alone, top 5 | 47 |

A cli story's pack was about 500 KB: 22 design and tech files by topics, about 30 ADRs by links, and 10 ranked items. `flai prime --cat` is 56 KB. The bulk is whole-file topics (`[all]` on seven `design/system` files, `cli` on `flai-cli.md`) and the ADRs `flai-cli.md § Commands` links. Heading topics on those files are the remedy ADR-0047 names. TH-0029 and TH-0030 asked the designer about it; [ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) answers with a budget and briefs, and S-0148 switched agents to the pack once S-0146 and S-0147 built them.

Since S-0148, with ADR-0049 built, the harness prompt, the MCP servers' instructions, `CLAUDE.md`, and `session-start.md` tell an agent with a story to prime with `flai prime --story <id>` or the MCP `prime`, and one without a story to prime with `flai prime --cat` ([conventions.md](conventions.md#priming)). The step [What an agent loads today](#what-an-agent-loads-today) measured is replaced: the pack instead of every convention, then what the agent fetches. The fetch is the agent's to choose, so every one of them carries the same rule: a brief is not the document, and when one bears on the story the agent reads the section that does, or the whole document, with `doc_get` and its heading before relying on it or changing what it describes. Packs of code stories stay over the 80 KB budget until the conventions' topics narrow (TH-0032); S-0148's own was 116 KB, over by what it names.

## Fitting the pack to a budget

S-0145 (2026-09-29) measured the S-0138 pack at 541 KB, about 135k tokens, against the 25k-token cap on an MCP tool result, and surveyed twelve ways of feeding an agent less: heading topics, a budget, BM25 and embedding retrieval, on-demand MCP tools, harness hooks, a memory layer, model-written digests, a scout sub-agent, prompt caching, a `context:` list, and splitting the long documents. The breakdown, the estimates, the comparison, and the recommendation are the proposed [ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md): a size budget, whole bodies only for what the story names, briefs (title, first paragraph, outline; an ADR's decision sentence) for what topics and links select, ranked sections to fill the budget, and `doc_search` plus a `heading` on `doc_get` so the agent fetches a section when it needs it. Estimated at about 80 KB for S-0138 with today's conventions. It supersedes ADR-0047's no-budget choice if accepted; the rest of ADR-0047 stands.

### As built, with a budget

S-0146 built ADR-0049's budget, briefs, and ranked fill in `flai/internal/context` ([flai-cli.md](flai-cli.md)). Choices ADR-0049 left open:

- The budget counts everything `flai prime --story` prints, the header included. A project sets it as `prime.budget` in `system-flow.yaml`, because a pack is the same for every agent on the project; `--budget` overrides it for one run.
- Briefs are never cut (TH-0032, the designer's choice). When the conventions, what is named, and the briefs exceed the budget, the header says which part took it over (`exceeded`: `conventions`, `named`, or `briefs`) and nothing is ranked. The pack tells the agent that a brief is not the document, and to read one whole when it bears on the story.
- Design briefs print under one `briefs` heading, ADR briefs one line each under `decisions`, each with its first reason and a count of the others; the JSON keeps every reason.
- The one step also reaches an ADR that refines a named or briefed one. A decision sentence of under five words, a bold lead-in, takes the next sentence with it.
- The ranked step tries every section that matches in rank order, not five and five, and keeps each that fits; an ADR whole first, then its best section.

Measured on 2026-09-29 on the S-0146 branch, with every convention still `[all]`:

| Story | Pack | Conventions and issues | Named | Briefs | Over by |
|-------|------|------------------------|-------|--------|---------|
| S-0138 | 125 KB | 57 KB | 36 KB: ADR-0047, ADR-0049, `conventions.md` | 55, 14 KB of text | named |
| S-0141 | 160 KB | 57 KB | 72 KB: `flaiover-dashboard.md` | 49 | named |
| S-0146 | 207 KB | 57 KB | 120 KB: ADR-0049, `agent-context.md`, `flai-cli.md` | 53 | named |

ADR-0049's estimate of 80 KB for S-0138 missed the 9 KB open-issues table and the header, and S-0138 was re-scoped after it to name 36 KB instead of 9 KB. The conventions and issues are more than half the budget before anything else. What takes packs furthest over is what stories and tasks name: a task that writes out `design/system/flai-cli.md` as a file to update loads all 86 KB of it. TH-0032 proposes briefing a path written out in plain text, keeping links and ADR IDs whole.

### Paths written out

[ADR-0050](../adrs/0050-a-document-a-story-names-only-by-its-path-written-out-is-briefed-when-it-is.md) (S-0149) refines ADR-0049's named step. A document the story, its epic, or its tasks name only by its repository path written out is briefed when the file is larger than an eighth of the budget (10 KB at 80 KB), and loads whole at or under it. A markdown link or an ADR ID to the same document still loads it whole, and a `#fragment` link still loads its section; the brief then outlines the rest with that section marked loaded. The brief's reason is `named in <ID>`, and the pack tells the agent that the story named the document, that it is briefed for its size, and that the agent reads it, or the sections it will change, before relying on it or changing it. A link is a reading instruction; a path is a pointer to where the work lands.

Measured on 2026-09-29 on the S-0149 branch against the flai before it, same documents, same 80 KB budget:

| Story | Pack before | Pack after | Named before | Named after | Briefed for size |
|-------|-------------|------------|--------------|-------------|------------------|
| S-0138 | 123 KB | 113 KB | 35 KB | 24 KB | `conventions.md` |
| S-0141 | 156 KB | 87 KB | 70 KB | 0 KB | `flaiover-dashboard.md` |
| S-0146 | 230 KB | 113 KB | 143 KB | 24 KB | `agent-context.md`, `conventions.md`, `flai-cli.md` |
| S-0147 | 188 KB | 102 KB | 101 KB | 15 KB | `flai-cli.md` |

What is left named is ADRs by ID (ADR-0047 9 KB, ADR-0049 15 KB), which ADR-0050 keeps whole. Every pack is still over the budget, now by the briefs or by those ADRs, because the conventions and open issues take 57 KB of the 80 KB; narrowing the conventions' topics is still the remedy.

## Sub-agents

An agent `flai serve` starts with `claude-code` worked its whole story in one context: every search, test log, and file it read stayed there and was read again from the cache on every later turn. S-0118's run cost 5.18 US dollars with 10.5M cache-read tokens for one Explore sub-agent it spawned on its own. [ADR-0059](../adrs/0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md) (S-0175) makes sub-agents deliberate: the story's agent hands noisy work and pre-review verification to an explorer and a verifier that can only read, primed for their role, and keeps decisions, edits, and the designer to itself. ADR-0049 rejected a scout that reads more than the agent; these read less.

| Role | Does | Tools | Primes with |
|------|------|-------|-------------|
| `explore` | Searches and reads code, design, and logs; returns what it found, with paths and lines | `Read`, `Grep`, `Glob`; flai's `prime`, `doc_get`, `doc_search`, `item_get`, `thread_get`, `board`, `who_touches` | `prime --role explore` |
| `verify` | Reads the diff, runs the tests, lint, and `flai check`; returns what fails against the criteria and the conventions | The explorer's, and `Bash` to run them | `prime --role verify` |

The template ships them as `.claude/agents/explorer.md` and `verifier.md`. Their `tools` is an allowlist, so neither has `item_move`, `item_edit`, `item_new`, `inbox`, `wait_for_work`, `wait_for_events`, `thread_open`, `thread_reply`, or `thread_resolve`, nor `Edit` or `Write`.

Each names a model cheaper than the story's agent's (S-0189), because neither decides anything: the explorer runs `haiku`, which is enough to find and read and quote paths and lines, and the verifier runs `sonnet`, because it judges a diff against the criteria and the conventions, and the defects it finds before review are what delegation buys (§ Measured). Both are the harness's aliases rather than model IDs, so a project made from the template follows the harness's current model of each class without an edit. A project that wants another sets `model` in the definition; `inherit` runs a sub-agent on the story's agent's. A story or the project's default can name one per role too, without touching the definitions: the agent's `roles` map `explore` and `verify` to a harness, model, and config of their own (S-0189, [ADR-0065](../adrs/0065-a-story-s-agent-carries-a-model-per-sub-agent-role-and-claude-code-runs-each.md)). `flai serve` starts `claude-code` with `--agents` built from the project's definition of each such role, with the role's model over the definition's. A role on another harness, or with config, is refused, because a sub-agent runs in the story's own session. The `command` harness gets the roles as `FLAI_AGENT_ROLES`. Because the verifier runs a different model, the corrections it finds are the story's agent's to make. It makes them itself, never through a sub-agent, as `delegation.md` and the prompt say. The guard refuses any sub-agent the git and flai commands that would make or record a correction (`TestASubAgentCannotMakeTheCorrectionsItFinds`). It does not refuse a shell command that writes a file, such as `sed -i` or a redirection. The verifier has no `Edit` or `Write`, and its definition forbids those commands, but nothing stops a careless one (I-0052).

The verifier's `Bash` could still run `flai move` or `git commit`, and the built-in sub-agents have every flai tool. [ADR-0060](../adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) closes that: the template's `.claude/settings.json` runs `flai guard` as a `PreToolUse` hook before `Bash` and every flai MCP tool, and it refuses any sub-agent's flai MCP write or inbox call, flai command that writes, and git command that changes the worktree, the index, branches, or history, with exit 2 and the reason. The story's agent's calls carry no `agent_id` and pass. Before review a verifier sub-agent fed the guard command lines that got past its first form, which trusted the first word: an assignment with a path in it, `env -i`, `sudo -u`, `timeout`, `xargs`, `find -exec`, `bash -lc`. The guard now looks at every word and keys on `agent_id`, since a session started with `--agent` may carry an `agent_type` of its own; a refinement of ADR-0060, not a reversal. It still stops a sub-agent that is careless, not one that hides a command (a backslash in a name, a command in a variable). This repository runs it as `scripts/flai.sh guard`.

### What the harness does, measured

A probe on 2026-10-01 with Claude Code 2.1.286, started as `flai serve` starts it (`claude -p`, flai's `--mcp-config`, `--strict-mcp-config`, `acceptEdits`, `--allowedTools Bash,mcp__flai`):

- The project's `.claude/agents/` load in a headless session started in the project: the init event lists them with the built-in `Explore`, `general-purpose`, and `Plan`. The adapter needs no flag.
- A definition's `tools` restricts MCP tools too: a sub-agent given `mcp__flai__item_get` and `mcp__flai__doc_search` had those and none of flai's others.
- A definition's own `mcpServers` did not connect under `--strict-mcp-config`, so a sub-agent cannot have a flai server, and a name, of its own.
- Sub-agents run in the background by default, and every event of theirs in the stream-json log carries `parent_tool_use_id`.
- A Bash pattern in a sub-agent's `disallowedTools` (`Bash(touch:*)`) removes `Bash` from it whole; a `hooks` block in a sub-agent's definition did not fire.
- A `PreToolUse` hook in the project's `.claude/settings.json` fires for a sub-agent's calls with `agent_type` and `agent_id` in its input, and for the session's agent without them. With `flai guard` as that hook, a `general-purpose` sub-agent's `item_move`, `git commit`, and `flai move` were refused and its `item_get` and `echo` ran.

### Questions

A sub-agent that needs the designer says so in its final message, with the question and its recommended answer, and the story's agent asks it with `thread_open` on the story or the task. Of the three ways S-0175 offered:

- **(a) Returned in the final message: chosen.** One voice per story, nothing to build in flai, and it works the same for any harness with sub-agents.
- (b) The sub-agent calls `thread_open`, attributed to the parent: the server reports to an agent only what others changed, so the parent's `wait_for_events` would never see a thread opened under its own name. It is also the write the definitions leave out.
- (c) The harness's messaging between agents (`SendMessage`, agent teams): Claude Code only, and a sub-agent that asked could not wait for the answer; the parent can.

No ADR changes who may open threads: the story's agent still does.

### Attribution

A sub-agent's flai calls arrive on its parent's MCP connection under its parent's name; the server cannot tell them apart and is not changed to. They need not be told apart there: every flai tool a sub-agent has reads and uses no agent identity, so nothing is attributed to anyone. They are told apart in the log `flai serve` keeps for the agent, where each event of a sub-agent carries `parent_tool_use_id`, the `Agent` call that started it, whose input names the role.

### The role pack

`flai prime --story S-nnnn --role explore|verify`, and `role` on the MCP `prime` tool, return a pack sized for a sub-agent:

1. The conventions whose front matter `roles` lists the role, with the sections the story's topics leave out taken out. No README and no open issues.
2. The story's `## Goal` and `## Acceptance criteria`.
3. Briefs, never bodies: what the story, its epic, and its tasks name, then what their topics select, then the ADRs one step reaches, in that order while the budget has room. What does not fit is counted in the header; `doc_search` finds it.

The budget is half the project's (`prime.budget`, else 80 KB), or `--budget`. The baseline gives `explore` to `communication`, `safety`, `tooling`, and `delegation`, and `verify` to those and `documentation`, `code-quality`, `git`, and `logging`.

### Measured

Read on 2026-10-01 at 09:18Z (S-0188) from the logs `flai serve` kept, with ADR-0051's reading: assistant events deduplicated by message ID, and each run's final `result` for its cost and turns. A sub-agent's calls are those whose events carry `parent_tool_use_id`. `result.usage` counts the story's agent only, while `modelUsage` and `total_cost_usd` include the sub-agents, so a sub-agent's share of the cost is estimated from its share of the run's input and cache tokens, as ADR-0051 apportions a task. Minutes run from a log's first event to its last. Every run is claude-opus-5-5 on the same host, and every sub-agent ran on the same model. Size is the story's commits outside `wip/` and lockfiles: distinct files, and lines added and removed summed over the commits. Two logs, S-0117's and S-0173's first, failed at their first call with an API error and are left out.

| Run | Nature | Size | Delegated | Model calls, agent / sub-agents | Cache reads, agent / sub-agents | Cost (sub-agents, est.) | Turns | Minutes |
|-----|--------|------|-----------|---------------------------------|---------------------------------|-------------------------|-------|---------|
| S-0184 | improvement | 29, +1092 −85 | released prompt: explorer, verifier | 124 / 46 | 23.30M / 2.91M | 10.82 (1.25) | 149 | 23.1 |
| S-0185 | remediation | 18, +433 −82 | released prompt: verifier ×3 | 105 / 38 | 15.48M / 1.66M | 7.83 (0.82) | 114 | 23.0 |
| S-0183 | remediation | 22, +204 −59 | convention only (1.26.3 prompt): verifier | 103 / 30 | 12.52M / 1.39M | 5.89 (0.61) | 115 | 16.0 |
| S-0175 | improvement | 56, +1560 −176 | by hand: Explore, general-purpose ×3 | 181 / 65 | 43.20M / 2.85M | 17.12 (1.13) | 201 | 40.8 |
| S-0118 | improvement | 23, +451 −83 | unprompted: Explore | 80 / 17 | 9.63M / 0.82M | 5.18 (0.43) | 90 | 15.7 |
| S-0179 | remediation | 43, +3114 −27 | none | 144 / 0 | 31.58M / 0 | 12.68 | 151 | 35.5 |
| S-0174 | remediation | 26, +1038 −51 | none | 143 / 0 | 25.25M / 0 | 9.19 | 158 | 22.7 |
| S-0173 | remediation | 28, +746 −134 | none | 103 / 0 | 13.06M / 0 | 5.19 | 145 | 13.9 |
| S-0180 | remediation | 12, +140 −28 | none | 83 / 0 | 8.27M / 0 | 3.56 | 93 | 11.5 |
| S-0119 | feature | 10, +171 −28 | none | 60 / 0 | 5.15M / 0 | 2.55 | 69 | 7.5 |
| S-0117 | remediation | 6, +89 −10 | none | 38 / 0 | 2.74M / 0 | 1.56 | 42 | 4.8 |

S-0184 and S-0185 are the two runs with the released prompt: flai 1.26.4, the first release with S-0175, which `flai serve` ran from 08:39Z. S-0183 started at 08:32Z under 1.26.3, after the definitions and `delegation.md` reached main, and delegated with the convention but not the prompt. S-0175's row is its finished run. An earlier reading at 08:13Z, before its last tasks, gave 12.00 US dollars estimated over 25.5 minutes.

Each delegating run is set against the non-delegating runs most like it in nature, size, and components. No improvement in the logs ran without delegating, so S-0184's comparable is a remediation.

| Delegating | Comparable | Why comparable | Cost | Minutes | Agent's cache reads per call |
|------------|------------|----------------|------|---------|------------------------------|
| S-0184 | S-0174 | Both change flai and flaiover, about 1,100 lines in 26 to 29 files | 10.82 against 9.19 | 23.1 against 22.7 | 188k against 177k |
| S-0185 | S-0174, S-0180 | All three are remediations centred on flai, and S-0185's size falls between the other two | 7.83 against 9.19 and 3.56 | 23.0 against 22.7 and 11.5 | 147k against 177k and 100k |
| S-0183 | S-0180 | Both are small remediations to flai, scripts, and design | 5.89 against 3.56 | 16.0 against 11.5 | 122k against 100k |

What the runs show:

- **The sub-agents' own spend was about a tenth of each run, and no saving was measurable.** Sub-agents took 10 to 12% of the tokens in S-0183, S-0184, and S-0185: 0.61 to 1.25 US dollars a run. Each delegating run cost about what its size predicts from the runs that did not delegate. S-0184 cost more than S-0174 for the same size and time. S-0185 cost less than S-0174 at less than half its lines, but took as long. Story size explains most of the spread, and two or three runs on each side are too few to show whether the tenth was extra or replaced work the agent would have done.
- **Sub-agents were added to the story's agent's work, not taken from it.** S-0184's agent ran the whole Go suite three times, once before its verifier and twice after, and about a dozen other test or lint runs. It also made 25 `Read` and 86 `Bash` calls itself. S-0185's agent ran the suite three times around its three verifiers: one for tests and lint, then two for the diff. The agents already cut their own test output to the tail or the failures, so a sub-agent had little noise to take away.
- **What delegation bought was review before review.** S-0184's verifier found two convention gaps, and S-0185's found five points. Both agents fixed them before moving to review. S-0185's last verifier confirmed the five were addressed and found one more regression, which the agent fixed too. S-0184's explorer found that an earlier attempt's work was lost with its worktree, without the agent reading the history itself.
- **Agents delegate without the prompt, and more with it.** S-0183 delegated once under the convention alone, and S-0118 once before either existed. The two runs with the prompt delegated two and three times.
- Reading the logs: the stream's deduplicated cache reads fall short of `modelUsage` by under 1% in four of the five delegating runs (S-0118 by 81,606, S-0184 by 112,679, S-0185 by 72,830, S-0175 by 56,989), and match it exactly in S-0183 and in every run without sub-agents. `num_turns` is not the number of model calls, and the stream's output tokens are about 2% of the `result`'s.

What to change follows from this. Keep the verifier before review, which found defects. Make its run replace the agent's own full-suite runs rather than repeat them. Run both sub-agents on a cheaper model, since they read and run checks and decide nothing. S-0189 makes both changes, and S-0190 measures them as this section does once two stories have run with them.

### A cheaper verifier, replayed

Before the template's verifier moved to `sonnet`, S-0189 checked that a cheaper verifier still finds what S-0185's found. On 2026-10-01, S-0185's second verifier call was replayed word for word, its prompt taken from `flai serve`'s log. That call was the first to check the diff, and it ran on `claude-opus-5-5`. The replay ran on `sonnet` against a detached checkout of `fa6aaa3`, the last S-0185 commit before its review fixes. The prompt's paths and diff base were changed to that checkout and S-0185's base, and nothing else.

| Finding | Opus, first diff check | Sonnet, replay |
|---------|------------------------|----------------|
| Docs say `flai import` refuses an empty required variable before writing, but it moves folders first | Found | Found |
| `flai upgrade --var` is ignored when the project is already at the template's version | Found | Found |
| `--relock` now fails on an unreadable lock, because the lock is read before the relock branch | Not found; S-0185's final verifier found it after the fixes | Found |
| Missing tests (`internal/lock`, `--relock` recording vars, a refused dry-run, the lock unchanged on refusal) | Found | Found |
| A project made by an older `flai new` silently takes a fork variable's default | Found | Not found |
| A recorded empty value is refused as missing before a newer default is tried | Found (low) | Not found |
| The behaviour commit leaves the docs stale until the next commit | Found | Not found |
| The five manifest-held variables are copied into the lock's `vars`, where they have no effect; a fork variable holding a secret would be committed in the lock | Not found | Found (minor) |

| Run | Model | Tokens | Tool calls | Seconds |
|-----|-------|--------|------------|---------|
| S-0185, first diff check | claude-opus-5-5 | 78,082 | 20 | 181 |
| Replay | sonnet | 78,415 | 18 | 165 |

Sonnet found four of the six things Opus found, and the regression Opus found only on its last run. It missed two migration edge cases and the commit split. It read about as much and took about as long. Its tokens are billed at Sonnet's rate rather than Opus's. This is one replay of one prompt. It shows that a cheaper verifier still catches the defects that matter most, not that it catches as many.

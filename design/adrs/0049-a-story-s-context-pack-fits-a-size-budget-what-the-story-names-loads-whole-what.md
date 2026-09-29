---
id: ADR-0049
title: "A story's context pack fits a size budget: what the story names loads whole, what its topics and links select loads as briefs, ranking fills the rest, and the agent fetches sections on demand"
status: proposed
date: 2026-09-29
supersedes: []
superseded_by: []
refines: [ADR-0047]
topics: [cli, conventions]
---

# ADR-0049 A story's context pack fits a size budget: what the story names loads whole, what its topics and links select loads as briefs, ranking fills the rest, and the agent fetches sections on demand

## Context

ADR-0047 gave every agent a context pack: the conventions its story's topics select, the `design/system`, `design/tech`, and ADR sections those topics select, everything the story names followed one link step, ten ranked items, and a catalog of the rest, with no size budget. S-0136 and S-0137 built it. Measured on this repository on 2026-09-29 with `flai prime --story S-0138 --json`, the pack of a cli story is 541 KB, about 135k tokens:

| Part | Items | Size |
|------|-------|------|
| Conventions, every section (every convention is still `[all]`) | 13 | 48 KB |
| Design and tech loaded whole because a file topic matched | 24 | 261 KB |
| Of which six files over 15 KB (`flai-cli.md` alone is 82 KB) | 6 | 176 KB |
| ADRs the story, its epic, and its tasks name | 1 | 9 KB |
| ADRs one step from a topic-selected section, ten of them from `flai-cli.md § Commands` | 33 | 125 KB |
| Ranked: five `flaiover-dashboard.md` sections of 5 to 17 KB and five ADRs | 10 | 59 KB |
| The heading outline of every topic-selected document | | 4 KB |

No agent can take that in one call. Claude Code keeps about 30 KB of a Bash result inline and writes the rest to a file; it caps an MCP tool result at 25k tokens, about 100 KB. An agent that did read it all would spend 135k tokens before its first action, in the region where models use the middle of the context worst (Liu et al.; Chroma, see [agent-context.md](../system/agent-context.md)). S-0138, which switches agents to the pack, is waiting on this (TH-0029, TH-0030). The remedy ADR-0047 named, heading topics on the long files, is upkeep with no ceiling: by estimate it brings this pack to about 270 KB, still three times what a tool result carries.

What the pack is for has not changed: an agent must read every rule that applies to it, because a rule it does not see is a rule it breaks, and it must know which decisions and design its story builds on. What has to change is that knowing about a document and reading it are two different costs, and the pack pays the second for everything.

S-0145 surveyed the ways of giving an agent only what its story needs (the alternatives below): metadata selection as today, lexical and embedding retrieval, on-demand tools over MCP, harness hooks, a memory layer, model-written digests, a scout sub-agent, and restructuring the documents. Every approach that works in practice keeps a small core in context and fetches the rest just in time (Anthropic's context engineering guidance; Agent Skills; Anthropic's code-execution-with-MCP note, where reading tool definitions on demand instead of up front cut 150k tokens to 2k). The replay in agent-context.md also says what is worth keeping: the one link step from topic-selected design is what reaches the decisions an archived story turned out to need (91 of 93 ADRs), and ranking adds the last few.

## Decision

**A story's context pack fits a size budget. Conventions load as ADR-0047 says, section by section by topic, and are never cut. What the story, its epic, and its tasks name loads whole. Everything else the pack selects loads as a brief, not a body: a topic-selected design or tech document as its title, first paragraph, and heading outline, and an ADR reached by one step as its title and decision sentence. Ranked sections, each cut at its own heading, fill what the budget leaves. The agent fetches any body it needs on demand: the MCP server gains `doc_search` (the BM25 index flai already has, over sections), and `doc_get` and `flai doc show` take a heading, so one section comes back, not the file.** This supersedes ADR-0047 on three points: there is a budget; topic-selected design is briefed, not loaded; and the one step from topic-selected sections yields decision sentences, not whole ADRs. The rest of ADR-0047 stands: the topics vocabulary, topics on documents and items, conventions by topic, the ranked step, the catalog, and a reason on every item.

### The pack

`flai prime --story S-nnnn`, in this order, within `--budget` (default 80 KB, about 20k tokens, so that the pack fits an MCP tool result with room for the agent's own turn; a project sets its default in the config under `prime.budget`):

1. **Conventions.** As ADR-0047 §The pack step 1: `README.md`, every convention with the sections whose topics do not match left out, and the open-issues table. Never cut. When the conventions alone exceed the budget, the header says so and the pack is the conventions and a catalog; narrowing convention topics is the designer's remedy, as ADR-0047 planned.
2. **Named.** Every document and ADR the story, its epic, and its tasks link or name by ID, whole (a linked `#fragment` loads its section). A superseded ADR is replaced by what supersedes it. These are the agent's declared reading and are never cut; a story that names more than the budget holds is a story to split, and the header says the budget was exceeded.
3. **Briefed.** Every `design/system` and `design/tech` file whose topics match, or that a heading topic selects, as its title, first paragraph, and heading outline. Every ADR reached one step from a named or briefed document, or that refines or supersedes one, as its ID, title, and decision sentence (the first sentence under `## Decision`). The reason on each says what selected it, as today.
4. **Ranked.** The sections that rank highest by BM25 against the story's title, goal, and criteria among those not named, each cut at its own heading (one heading down to the next heading of any level), until the budget is spent. ADRs rank by their best section and load whole when they fit. At least the five and five of ADR-0047 are tried; the budget decides how many load.
5. **Catalog.** As ADR-0047 §The pack steps 5 and 6: one line per document not loaded or briefed, and one per convention section left out, with its topics.

Every item carries its reason and its size in the header, in text and in `--json`. A brief says how to get the body: `doc_get` with `heading`, or `flai doc show <path> --heading "<heading>"`.

### On demand

- `doc_search` (MCP) and `flai doc search` rank sections of `design/`, `docs/`, and the conventions by BM25 against a query and return path, heading path, the first lines, and size, at most twenty. It is `flai/internal/search`, which the dashboard already uses, over the same section cuts the pack makes.
- `doc_get` and `flai doc show` take `heading`: the section under that heading, with the sections below it, and nothing else. Without it they return the whole document as today.
- The MCP `prime` tool (S-0138) returns this pack and nothing larger. The harness prompt, `CLAUDE.md`, and `session-start.md` tell the agent that the pack is a brief, and to read a briefed document's body with `doc_get` before changing what it describes.

### Estimate

On S-0138, with the conventions as they are today:

| Part | Size |
|------|------|
| Conventions | 48 KB |
| Named: ADR-0047 | 9 KB |
| Briefed: 24 design and tech outlines with first paragraphs, 33 ADR decision sentences | about 16 KB |
| Ranked, within what is left of 80 KB | about 5 KB |
| Catalog | about 3 KB |
| Pack | about 80 KB, about 20k tokens, from 541 KB |

Until the designer narrows the conventions' topics, the conventions are more than half the pack and the ranked step gets little room; every convention section they narrow away is room for ranked design. The ADRs an archived story turned out to need stay reachable: the step that reached 91 of 93 in the replay still runs, and what it reaches is one `doc_get` away with its decision in front of the agent.

## Consequences

- An agent starts with about 20k tokens of context instead of 135k, all of it either a rule that applies, a document its story names, or a one-line brief of what else exists and why it was selected. The MCP `prime` tool can return the pack. S-0138 can switch agents to it.
- Reading becomes a choice the agent makes, so a body it should have read and did not is a new way to fail. Three things limit it: the rules are never briefed, what the story names is never briefed, and each brief carries its reason and the decision sentence, which is what an ADR is for. The archive replay in agent-context.md gets a second measure once this lands: of the ADRs archived stories named, how many were loaded or briefed.
- `flai/internal/context` changes in three places: a budget and a size on every item; briefs for topic-selected documents and stepped ADRs; ranked sections cut at their own heading and bounded by the budget. `flai/internal/search` gains a section-level index of `design/` and the conventions, shared by the pack, `doc_search`, and the dashboard. `doc_get`, `flai doc show`, and the MCP server gain `heading`. `docs/users/flai.md`, `flai-reference.md`, `design/system/flai-cli.md`, and `agent-context.md` change with them.
- The first sentence under `## Decision` becomes the ADR's brief, so it has to state the decision. Every accepted ADR already does; the template says so from now on, and `flai check` warns about an ADR whose `## Decision` opens with anything but a sentence.
- Heading topics on the long design files still matter, for what a section is about, but they are no longer what keeps the pack small. The budget is.
- Nothing new is depended on: the same Go, the same BM25, the same front matter. The pack stays deterministic for a repository and a story, and offline.
- The stories under E-0010 that build it, in order: the budget, briefs, and ranked cuts in `flai/internal/context`, with the estimate above checked against `flai prime --story S-0138`; `heading` on `doc_get` and `flai doc show` and the `doc_search` tool; the prompt, `CLAUDE.md`, and `session-start.md` wording, which S-0138 already touches and can take.

## Alternatives considered

The story asked that every approach be listed with its trade-offs. What each removes from the measured pack, what it costs to build and keep true, how it fails, and how it fits flai (Go, offline, deterministic, explainable, harness-neutral).

| Approach | Removes | Costs | Fails when | Fit |
|----------|---------|-------|------------|-----|
| Heading topics on the long files (ADR-0047's remedy) | Parts of the 261 KB of topic-selected design; about 270 KB left by estimate | Metadata on every long file, kept by hand as files change | A topic is missing or too broad; no ceiling | Keeps the vocabulary; not enough alone |
| A budget that truncates today's pack | Whatever comes last | One number | What is cut is the ranked step, the most story-specific part; a whole-file topic still spends 82 KB on one document | Necessary, not sufficient |
| Lexical retrieval only (BM25 at prime time) | Everything not ranked | Nothing new: `flai/internal/search` | Recall: 47 of 93 needed ADRs alone; a rule becomes probabilistic | In flai already; it is the ranked step |
| Embedding retrieval (vector RAG) | As above, with better recall on paraphrase | A model (an API or a local one), a store, an index rebuilt on every edit; results move with the model | Offline runs, determinism, explaining why a section was chosen; a corpus of 600 KB in 200 documents is small for it | Poor: two new dependencies for a corpus BM25 covers. Anthropic's contextual retrieval shows embeddings plus BM25 beat either, on corpora far larger than this |
| On-demand tools (progressive disclosure): search and fetch by section over MCP | Every body the agent does not ask for | `doc_search`, `heading` on `doc_get`; the agent must ask | The agent does not look; mitigated by briefs with reasons and by never briefing rules | Good; MCP resources (`flai://design/...`) already exist, and tools are what agents call. Anthropic measured 150k to 2k tokens for tool definitions read on demand |
| Harness hooks that inject a rule when a path is read or edited (Claude Code `PreToolUse` and `PostToolUse` `additionalContext`, 10k characters; Cursor and Copilot path rules) | Path-scoped conventions from the pack | A hook per harness; flai would print a rule for a path | The first turns run without the rule; nothing for design or ADRs; one harness at a time | Later, as an experiment for conventions, once the pack is small; not harness-neutral |
| A memory layer kept across sessions (Claude Code auto-memory; Anthropic's memory tool, a client-side file store, 39% better with context editing in Anthropic's evaluation) | Nothing from the pack: memory is the agent's notes, not the project's design | A store per agent; duplicates `design/` and drifts from it | flai starts a fresh agent per story, so memory is empty when it matters; a stale note contradicts the living design | Poor for this problem; `design/system` is the project's memory already |
| Model-written digests of the design (a summary per document, regenerated on change) | Bodies, replaced by prose summaries | An LLM call per changed document, in CI or by hand; a lossy text whose truth no check can test | A summary omits the detail the story needed; drifts if regeneration is skipped | Later, if the deterministic brief (first paragraph, outline, decision sentence) proves too thin; keep it behind the same brief interface |
| A scout sub-agent that reads the whole pack and writes a brief for the worker | The worker's 135k tokens, spent in another window instead | The full pack is still read, once per story, by a second model; latency; a harness feature | Anthropic: multi-agent runs use about fifteen times a chat's tokens and suit breadth-first research, not most coding tasks | Poor: spends what it saves, and harness-specific |
| Prompt caching | Nothing: the same tokens are read, cheaper and faster on a cache hit | Nothing | The 25k-token tool cap and the context window are unchanged | Complementary, not a fix |
| A hand-written `context:` list on the story | Whatever the writer left out | Refinement effort, declined in ADR-0047 | Skipped under time pressure | Step 2 already loads what the story names |
| Splitting the long documents (`flai-cli.md § Commands` into one file per command) | Nothing by itself; makes every selection finer | A documentation story; links to update | It is not a mechanism | Worth a story of its own; the budget does not depend on it |

Sources beyond those in agent-context.md: Anthropic, [Code execution with MCP](https://www.anthropic.com/engineering/code-execution-with-mcp) (2025-11-04); Anthropic, [Managing context on the Claude Developer Platform](https://claude.com/blog/context-management) (2025-09-29, the memory tool and context editing); Anthropic, [How we built our multi-agent research system](https://www.anthropic.com/engineering/multi-agent-research-system) (2025); Claude Code [hooks](https://code.claude.com/docs/en/hooks); MCP [resources](https://modelcontextprotocol.io/specification/2025-06-18/server/resources); [llms.txt](https://llmstxt.org/), a curated index of links with one-line descriptions, which is what the catalog is.

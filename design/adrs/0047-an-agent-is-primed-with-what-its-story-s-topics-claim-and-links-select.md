---
id: ADR-0047
title: "An agent is primed with what its story's topics, claim, and links select: conventions and design by topic, the documents it names, ranked sections, and a catalog of the rest"
status: accepted
date: 2026-09-26
supersedes: []
superseded_by: []
refines: [ADR-0001, ADR-0013]
---

# ADR-0047 An agent is primed with what its story's topics, claim, and links select: conventions and design by topic, the documents it names, ranked sections, and a catalog of the rest

## Context

An agent flai starts reads every convention whatever its story (`flai prime --cat`: twelve files, 47 KB, and the whole open-issues table), then finds the design and the decisions it builds on by itself, or does not. `design/` is about 540 KB here; loading it all would crowd out the work, and models use the middle of a long context worst. E-0010 asks that flai determine the relevant documentation for a story from the component, the technologies, and the architecture it affects.

S-0125 surveyed how agent tools choose context and replayed the archive; the findings are in [agent-context.md](../system/agent-context.md). Every tool combines a small always-on core, rules scoped by path, and on-demand reading of the rest. Of the 44 existing ADRs that 26 archived stories named, links from the story's component reached 77%, BM25 ranking of title and goal (top 5) 80%, and the two together 93%. Studies of instruction files find generic overviews cost tokens without helping; specific, relevant instruction is what helps.

On TH-0020 (2026-09-26) the designer chose design A, the context pack, and asked that conventions be selected by topics on the file and on headings, called `topics` so as not to overload `tags`, which decide release components. On TH-0021 the designer chose: `topics` may be added to an accepted ADR (1b); `topics` on conventions, design, and tech files, with `[all]` for core (2a); every convention starts as `[all]` for the designer to review (3); no budget (4b); agents switch as soon as the pack exists (5b).

## Decision

**`flai prime --story S-nnnn` prints a context pack for the story: the conventions whose topics match the story's, section by section; the `design/system`, `design/tech`, and ADR sections whose topics match; every document the story, its epic, and its tasks link, followed one step; the best-ranked remaining sections; and a one-line catalog of everything not loaded. Each item says why it is there. The harness prompt, `CLAUDE.md`, the MCP server, and the template's session-start convention use it in place of `flai prime --cat`.** This refines ADR-0013 (conventions are still one file per topic, read in order, but a story reads the parts that apply to it) and ADR-0001 (an accepted ADR may gain `topics`, as it may gain `superseded_by`).

### Topics on documents

- **Front matter.** Conventions, `design/system` files, `design/tech` files, and ADRs may carry `topics: [...]`. `[all]` means every story. A convention without `topics` is read as `[all]`; a design, tech, or ADR file without it is never selected by topic, only by links or ranking.
- **Headings.** Any heading may carry `<!-- topics: a, b -->` on its line. It scopes everything down to the next heading at the same level or higher. A heading without it takes its parent's topics; a file's topics are the parent of its top-level headings.
- **Vocabulary.** A topic is a word. Topics a story can match without declaring them come from `system-flow.yaml`: each sub-project's name, tags, and kind, and `code` for any sub-project whose kind is not `template`. The template's baseline uses only those generic words (`code`, a kind such as `go`); a project uses its own under `## Project additions`, or anywhere in its design.
- **ADRs.** `topics` is the one key besides `superseded_by` that may be set on an accepted ADR, and only with `flai adr topics ADR-nnnn <topics>`. The body and every other key stay immutable.
- **Rollout.** Every convention, in the template and in this repository, gets `topics: [all]` when this lands; the designer then narrows them. Design and tech files gain topics as their owners add them; `flai check` warns about a design or tech file with none, and about a topic that no sub-project, story, or epic uses.

### Topics of a story

A story's topics are the union of:

- its own `topics:` and its epic's (a new optional front matter key on stories and epics, set with `--topics`, `flai edit`, and the dashboard);
- the name, tags, and kind of every sub-project its `tags` or its claim reach (the claim of [ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md): its `touches` and its open tasks'), and `code` if any of those sub-projects is not the template;
- `all`.

`tags` keep their one meaning, the component a story delivers to; they count toward topics only through the manifest.

### The pack

`flai prime --story S-nnnn`, in this order:

1. **Conventions.** `README.md`, then each convention in read order with the sections whose topics do not match the story's left out. The baseline marker and the `## Project additions` heading are always kept. The open-issues table follows, as today.
2. **By topic.** The sections of `design/system`, `design/tech`, and ADR files whose topics match, whole documents where the file's topics match and no heading narrows them.
3. **Linked.** Every document or ADR that the story, its epic, and its tasks link or name by ID, and one step further: the ADRs a selected section links, the ADRs a selected ADR refines. A superseded ADR is replaced by what supersedes it.
4. **Ranked.** The five ADRs and five design sections that rank highest by BM25 (`flai/internal/search`) against the story's title, goal, and criteria, among those not already selected. The replay showed top 5 adds most of what links miss; the count, not a size, bounds it.
5. **Catalog.** One line per document not loaded (path and title), and the heading outline of each document loaded in part, so the agent reads the rest with `flai doc show` or the MCP `doc_get` when it needs it.
6. **Left out.** One line per convention section left out, with its topics (`code-quality.md § Go (cli)`), so the agent knows the rule exists.

Every item is headed with its path, its heading path when it is a section, and its reason (`topics: cli`, `linked from S-0125`, `refined by ADR-0046`, `rank 3`). `--json` returns the same as data. There is no size budget; the pack's header states its size so growth is visible. Sections are cut at headings only.

### Who uses it

- `flai prime` without `--story` is unchanged: every convention, whole, for a session with no story.
- The harness prompt (`harness.Prompt`) tells the agent to prime with `flai prime --story <id>`; `CLAUDE.md`, the template's `CLAUDE.md.tmpl`, and the baseline `session-start.md` say the same for any agent that has a story.
- The MCP server gains a `prime` tool that returns the pack for a story.

## Consequences

- An agent reads the rules that apply to its story and the decisions and design it builds on, without searching, and the designer can see why each was included.
- Until the designer narrows the conventions' topics, conventions are loaded as today; the gain at first is the design and ADRs the pack adds.
- There is no budget: a story that names many documents, or a broad topic on a long design file, makes a large pack. The header's size makes it visible; narrowing topics on headings is the remedy.
- `topics` is a new front matter key on items, conventions, design, tech, and ADRs. Items are parsed strictly, so a project must run a flai that knows the key before any item carries it; `flai upgrade` brings the template's conventions.
- The baseline rule that an accepted ADR is never edited except for `superseded_by` gains `topics`; it changes in `decisions.md`, `documentation-standard.md`, and the tooling that refuses ADR edits.
- The implementation is split into stories under E-0010: topics on documents and items with the check; the conventions filter in `flai prime --story`; design, tech, ADR selection with links, ranking, and the catalog; the switch of the harness, `CLAUDE.md`, the template, and the MCP tool.

## Alternatives considered

- **B: core conventions and a catalog only.** No metadata to keep, but every story spends turns choosing what to read, and a rule the agent does not know about is never read.
- **C: retrieval only.** No metadata; missed a fifth of needed ADRs at top 5 in the replay, and rules become hit or miss.
- **D: a hand-written `context:` list on each story.** Precise, a burden at refinement; the pack already loads whatever the story links.
- **Embedding retrieval.** A model and a store as new dependencies, results that move with the model; the replay leaves little for it to add.
- **`applies_to` by path instead of topics.** Offered on TH-0021; the designer chose topics, which name what a document is about rather than where the code lives, and cover conventions that are not about paths.
- **Deriving ADR scope only from the documents that link it** (TH-0021 1a). Declined for explicit topics on ADRs.
- **A size budget** (TH-0021 4a). Declined: load everything selected.
- **Switching after a replay test reaches 90%** (TH-0021 5a). Declined: switch as soon as the pack exists.

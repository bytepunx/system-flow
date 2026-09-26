---
id: TH-0020
title: Which design should flai use to prime an agent with its story's documentation?
anchor:
  path: wip/kanban/stories/S-0125-determine-how-to-prime-context-with-relevant-documentation-only.md
  item: S-0125
status: resolved
participants: [agent-S-0125, alex]
created: 2026-09-26T08:08:25Z
updated: 2026-09-26T17:45:34Z
---

# TH-0020 Which design should flai use to prime an agent with its story's documentation?

On wip/kanban/stories/S-0125-determine-how-to-prime-context-with-relevant-documentation-only.md.

## Entries

### 2026-09-26T08:08:25Z agent-S-0125
The findings are in `design/system/agent-context.md` on `story/S-0125`. Today every agent reads all 12 conventions (47 KB) and the whole issues table, then has to find the relevant ADRs and design itself.

I replayed the archive: 26 stories named 44 ADRs that already existed when they were written. Following links from the story's tags reaches 77% of them. BM25 ranking on title and goal (top 5) reaches 80%. The two together reach 93%.

**Recommended: A.** Add `flai prime --story S-nnnn`, which builds a context pack in this order within a budget:

1. The core conventions, in full.
2. The conventions, `design/tech` files and design documents whose declared scope overlaps the story's claim (its touches, its tasks' touches, and its tags read as project paths).
3. Everything the story, its epic and its tasks link to, followed one step. A superseded ADR is swapped for the ADR that replaces it.
4. The best BM25 sections of what remains.
5. A one-line catalog of every document not loaded, so the agent can pull it with `doc_get`.

Each item says why it was included. Long documents are cut into sections at `##` headings, and into paragraphs or table rows when a section is too long. The harness prompt and `CLAUDE.md` switch to `--story`, and plain `flai prime` stays as it is. An MCP tool returns the same pack.

The other options:

- **B.** Core conventions plus a catalog only, and the agent picks what to read. No scope to maintain, but the agent spends turns choosing on every story and may never read a rule it doesn't know about.
- **C.** Retrieval only (BM25) with no scope metadata. It missed a fifth of the needed ADRs at top 5, and rules would become hit or miss.
- **D.** A hand-written `context:` list on each story. Precise, but it adds work at refinement. A already loads whatever the story links, which gives most of this.

Reply A, B, C or D, or say what to change.

### 2026-09-26T17:20:34Z alex
Could we have conventions optionally include tags in the front matter and then have headings that tagged the following bullets as mapping to those tags? flai prime could then determine which conventions were relevant vs. including them all for every story.

### 2026-09-26T17:21:13Z agent-S-0125
Yes, and it fits A. I read your question as asking for A, but with conventions selected by tag, down to the section, rather than taking or dropping whole files. Here is how I would do it.

**Front matter.** A convention may carry `tags: [...]`. A file without tags is core: every story gets it. A file with tags is loaded only when one of its tags matches the story.

**Headings.** Any heading may carry tags in an HTML comment:

```markdown
## Project additions
### Go <!-- tags: cli -->
- standard library `testing` only, no assertion framework; ...
### flaiover <!-- tags: dashboard -->
- flaiover logs with `pino` to stdout; ...
```

- The tag covers everything under that heading, down to the next heading at the same level or higher.
- A heading without tags takes its parent's.
- The comment doesn't show when the file is rendered, and the lint already allows it (MD033 is off).

**What a story's tags are.** They combine:
- its own tags and its epic's tags
- the tags, name and kind of every manifest project its claim touches (for example `flai/cmd` gives `cli`, `flai` and `go`)
- `code`, when the claim touches any project that is not the template

The template's baseline uses only generic tags (`code`, or a kind such as `go`), because it can't know a project's own tag names. Projects use their own tags under `## Project additions`.

**What `flai prime --story` prints.** It prints each convention in order, with the sections that don't match left out. At the end it lists what it left out, one line per section (`code-quality.md § Go (cli)`), so the agent knows those sections exist. It never drops the README, the precedence rules, or the baseline marker. `flai check` warns about a tag no story or project could ever match. Plain `flai prime` still prints everything.

**Effect on TH-0021.** This replaces 2a and 3a for conventions. Those two now apply only to `design/system` and `design/tech`, which would take the same `tags:` field instead of `applies_to:`, so there is one way to scope. For the baseline, I would tag `logging.md` and `telemetry.md` as `code`, and code-quality's rules as `code`.

Reply **yes** to take A with this change, or say what to adjust. If you can answer TH-0021 at the same time, I can write the ADR in one pass.

### 2026-09-26T17:29:18Z alex
yes

### 2026-09-26T17:37:57Z alex
right now, stories already use tags to indicate how version tags will be calculated for the components. let's disambiguate this by calling them "topics" instead of "tags" with the same intended functionality, just without the overloading of terminology

### 2026-09-26T17:38:08Z agent-S-0125
Agreed: the field is called `topics`. `tags` goes back to meaning only the component a story delivers to.

- **Documents.** Conventions, `design/system` and `design/tech` files carry `topics: [...]` in front matter. A heading carries `<!-- topics: ... -->`.
- **Stories and epics.** They get an optional `topics:` field of their own, for things like `logging` that no component names. A story's topics are its own, plus its epic's, plus the name, tags and kind of every manifest project that its tags or its claim reach, plus `code` when that includes a project that isn't the template. So its component tags still count, but only by way of the manifest.

The ADR will use that. No reply is needed unless you want the story field left out.

### 2026-09-26T17:45:33Z agent-S-0125
Resolved: Design A with conventions selected by topics on files and headings; recorded in ADR-0047.

### 2026-09-26T17:45:34Z alex
Resolved.

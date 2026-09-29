---
id: ADR-0050
title: "A document a story names only by its path written out is briefed when it is larger than an eighth of the pack's budget"
status: accepted
date: 2026-09-29
supersedes: []
superseded_by: []
refines: [ADR-0049]
topics: [cli, conventions]
---

# ADR-0050 A document a story names only by its path written out is briefed when it is larger than an eighth of the pack's budget

## Context

ADR-0049 fits a story's context pack to a size budget (80 KB by default) and loads whole everything the story, its epic, and its tasks name: a markdown link, a document's repository path written out, or an ADR ID. What is named is never cut. The measurements S-0146 recorded in [agent-context.md](../system/agent-context.md) show that what stories name, more than the briefs, is what takes packs over:

| Story | Pack | Named | What |
|-------|------|-------|------|
| S-0138 | 125 KB | 36 KB | ADR-0047, ADR-0049, `conventions.md` |
| S-0141 | 160 KB | 72 KB | `flaiover-dashboard.md`, which T-0504 names as a file to update |
| S-0146 | 207 KB | 120 KB | ADR-0049, `agent-context.md`, and `flai-cli.md` (86 KB), which T-0525 names as a file to update |
| S-0147 | 189 KB | 102 KB | ADR-0049 and `flai-cli.md` |

A task that writes out a path to say "update this file" is saying where the work lands, not asking for the file to be read before anything else. It pays for the whole file anyway, and the two long design files are 72 and 86 KB: either is most of the budget alone, and the conventions and open-issues table already take 57 KB. A link, by contrast, is written to be followed, a `#fragment` link picks the section that matters, and an ADR is short and is the decision the story builds on. The designer asked for this refinement in TH-0032.

## Decision

A document the story, its epic, or its tasks name only by its repository path written out in plain text is briefed, not loaded whole, when it is larger than an eighth of the pack's budget (10 KB at the default 80 KB); at or under that it loads whole as before. A markdown link to it or its ADR ID still loads it whole, as ADR-0049 says, whatever else names it; a `#fragment` link still loads its section, and the path written out then briefs the rest, with that section marked loaded in the outline. The size is the whole file's, front matter included, and the threshold moves with `--budget` and `prime.budget`, so a project with a larger budget briefs less without another setting.

The brief is ADR-0049's: a design or tech file as its title, first paragraph, and heading outline; an ADR as its decision sentence. Its reason is `named in <ID>`, it prints among the briefs, before what topics select, and it is never cut, as nothing named is. The one link step does not follow the links of a document briefed this way, only those of the sections its topics select, as for any brief: a path is a pointer to where the work lands, not to what the file links. The pack tells the agent that the story named it and that it was briefed for its size, and that the agent decides from the brief whether it needs the body: when it will rely on the document or change it, it reads the whole file, or the sections it will change, with `doc_get` and a `heading` or `flai doc show <path> --heading`, first.

This refines ADR-0049's step 2, "what the story names loads whole": a link, a fragment link, an ADR ID, and a small file written out still do. The rest of ADR-0049 stands.

## Consequences

- A task can name the files it changes without spending the pack on them. S-0141's, S-0146's, and S-0147's packs lose the 72 and 86 KB files they wrote out; the numbers after the change are in [agent-context.md](../system/agent-context.md).
- A writer who wants a document read whole before work starts links it. Refinement of a story says so: a link is a reading instruction, a path is a pointer.
- An agent that edits a briefed file without reading it first is a new way to fail. The brief's reason and the pack's instruction name the read as a precondition of the change, so the cost is one fetch the agent makes when it gets to the change, not a file it carries from the start.
- `flai/internal/context` tells a path written out apart from a link and an ID when it finds what an item names, and the named step checks the size against the budget. `flai prime --story`'s help, `docs/users/flai.md`, `design/system/flai-cli.md`, and `agent-context.md` say so.

## Alternatives considered

- Brief every path written out, whatever its size: a 2 KB file costs almost as much as a brief, and the agent would then fetch it; the story's criterion is about large files.
- Brief every named document over a size, links and ADR IDs included (TH-0032's option C): reverses "named loads whole" for the cases where the writer asked for the document, and needs a way to say "really load this".
- A fixed threshold such as 12 KB: a project that raises the budget to fit its conventions would still brief files it has room for.
- A setting for the threshold: nothing yet needs it apart from the budget, which already is one.

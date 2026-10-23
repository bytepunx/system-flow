---
id: ADR-0191
title: "The digest is one markdown file per entry under the design folder's digest folder, named D-nnnn-slug.md, whose front matter flai writes and checks, and nothing in it is ever deleted"
status: accepted
date: 2026-10-23
supersedes: []
superseded_by: []
---

# ADR-0191 The digest is one markdown file per entry under the design folder's digest folder, named D-nnnn-slug.md, whose front matter flai writes and checks, and nothing in it is ever deleted

## Context

E-0025 gives a story's agent what past stories learned: the decisions in their narratives, the operator's answers on threads, the recommendations the operator confirmed or overruled, and the causes and fixes of closed issues. None of that reaches a pack today unless someone wrote it into a design document ([ADR-0047](0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), [ADR-0049](0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)). The digest needs a home and a format before anything fills it, reads it, or decays it.

The operator confirmed on TH-0476 (2026-10-21) one file per entry under `design/digest`, a curator role that keeps it (S-0431), and a 30-day compaction window (S-0432). [agent-memory.md § Recommendation](../system/agent-memory.md#recommendation) (S-0426) recommends that agents reach the archive and the threads through this digest under `design/`, which `doc_search` already reads, rather than through the raw files.

## Decision

**Where.** The digest is the folder `digest` inside the design folder the manifest's `layout.design` names, `design/digest/` by default, as `design/analysis/` is ([ADR-0099](0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md)). The folder's `README.md` says what it is for and is not an entry. It keeps no index of entries: `flai digest list` is the index, so that the many writers of one folder never conflict on a table.

**One entry, one file.** An entry is one fact a story learned, in `D-nnnn-<slug>.md`. `D-nnnn` is its ID, four digits, one more than the highest number in this checkout, on main, in every worktree, and on every story branch, as an ADR's number is, so that parallel writers do not take the same one; gaps are not filled. The slug is the title's, as an ADR's file name is, cut at 80 characters on a word.

**Front matter.** flai writes it, and `flai check` checks it:

| Key | What it holds |
|-----|---------------|
| `id` | `D-nnnn`, the number the file is named for |
| `title` | The fact, as one sentence |
| `kind` | `decision` (from a narrative's `## Decisions` or an ADR), `answer` (the operator's answer on a thread), `recommendation` (a recommendation the operator confirmed or overruled, and why), or `cause` (an issue's cause and fix) |
| `sources` | The work items (`E-`, `S-`, `T-`), threads (`TH-`), and issues (`I-`) the fact comes from, at least one, each in the project or its archive |
| `topics` | The topics it bears on, as a story's are, for priming (S-0430) |
| `paths` | The repository paths it bears on, matched against a story's touches |
| `date` | The day the fact was learned, `YYYY-MM-DD` |
| `status` | `active` (primed), `absorbed` (a design document, convention, or ADR now records it, so it is no longer primed), or `compacted` (no one cited it within the window, so only its title is primed) |
| `absorbed_by` | The path, from the project root, of the document that now records it; set when, and only when, `status` is `absorbed`, and the file exists |
| `asked_again` | The threads that asked what the entry already answered, filled by the curator (S-0431) and counted by the metric (S-0434) |
| `cited` | The work items and threads that cited the entry, filled by `flai digest cited` (S-0432) and read by compaction |

**Body.** Below `# D-nnnn <title>`, a few sentences: what was asked or decided, the answer and its reason, and links to its sources. It says where a secret lives, never its value: `flai check` reports `digest.secret` on an entry whose text has the shape of a key or token, and `flai digest new` and `flai digest edit` refuse one.

**Written through flai.** `flai digest new` allocates the ID and writes the file in the checkout it runs in, and `flai digest edit` changes an entry's body, topics, paths, sources, and `asked_again`; each runs the check on what it wrote and puts the file back as it was when the check reports anything the write introduced. `flai check` reports `digest.entry` for an entry whose front matter does not parse or misses a field, whose kind or status is unknown, whose source names nothing in the project or its archive, or whose `absorbed_by` names no file.

**Nothing is deleted.** No entry's file is removed. Absorption sets `status: absorbed` and `absorbed_by`; compaction sets `status: compacted` and empties the body below the heading. Git keeps every earlier text.

## Consequences

- `doc_search`, `doc_get`, and the dashboard's documents tree reach every entry with no change, since they walk the design folder; the documentation standard's own front matter rules (`updated`) do not apply to entries, as they do not to analysis reports or experiment results, because `flai check` checks the digest's own.
- The prime pack can select an entry by its topics and paths as it selects an ADR (S-0430), and the curator writes only through `flai digest` (S-0431), which the guard can hold it to.
- An ID is unique across parallel writers without a counter kept in a document, at the cost of reading the story branches' folders on each `flai digest new`.
- The template ships `design/digest/README.md`, so every project made from it has the folder.

## Alternatives considered

- **One file per topic, each entry a section.** An entry on two topics has to pick one, and two writers on one topic conflict on the file.
- **Front matter on the archived items.** Scatters the digest across `wip/archive`, which `doc_search` does not index, and edits files the archive keeps as they were.
- **An index table in the folder's README, as `design/analysis` keeps.** Every write would change one file, and parallel writers would conflict on it; `flai digest list` prints the same from the entries.
- **Deleting an absorbed or stale entry.** Loses the record that the question was answered, and the `asked_again` count with it; a state change keeps both and still drops the entry from priming.

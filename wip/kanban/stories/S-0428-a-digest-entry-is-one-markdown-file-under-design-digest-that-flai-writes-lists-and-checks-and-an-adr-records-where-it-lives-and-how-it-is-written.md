---
id: S-0428
type: story
nature: feature
title: A digest entry is one markdown file under design/digest that flai writes, lists, and checks, and an ADR records where it lives and how it is written
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:58:10Z
updated: 2026-10-16T06:03:37Z
transitions: []
tags: [cli, template]
topics: [cli, conventions, planning, analysis]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0428 A digest entry is one markdown file under design/digest that flai writes, lists, and checks, and an ADR records where it lives and how it is written

## Goal

The digest has a home and a format before anything fills it. An entry is one fact a story learned: a decision from a narrative's `## Decisions`, an answer the operator gave on a thread, a recommendation they confirmed or overruled and why, or an issue's cause and fix. Each entry is a markdown file under `design/digest/` (or the folder TH-0476 settles) whose title states the fact in one sentence and whose front matter gives its kind, the stories, threads, and issues it comes from, the topics and paths it bears on, its date, and its state. flai writes an entry with `flai digest new`, prints the digest with `flai digest list` and `flai digest show`, and `flai check` refuses an entry that does not parse, cites a source that does not exist, or carries what looks like a secret's value. `doc_search` finds each entry and the dashboard's documents page renders it, as they do an ADR.

## Acceptance criteria

- [ ] An ADR records where the digest lives within the layout, the manifest key that names its folder, an entry's file name and ID (`D-nnnn`), its front matter (`id`, `title`, `kind`, `sources`, `topics`, `paths`, `date`, `status` of `active`, `absorbed`, or `compacted`, `absorbed_by`, `asked_again`, `cited`), its body, and that nothing is deleted: absorption and compaction change an entry's state and body, and git keeps the rest.
- [ ] `flai/internal/digest` reads every entry of a project, and `flai digest list [--status] [--topic] [--json]` and `flai digest show <D-nnnn> [--json]` print them with their sources, topics, paths, date, and state.
- [ ] `flai digest new "<fact>" --kind <kind> --source <id>... [--topic] [--path] [--json]` allocates the next ID, writes the file in the main checkout, and runs the check on it, refusing and writing nothing when the check fails; `flai digest edit <D-nnnn>` changes an entry's body, topics, paths, sources, and `asked_again` the same way.
- [ ] `flai check` reports `digest.entry` for an entry whose front matter does not parse or misses a field, whose kind or state is unknown, whose source names no work item, thread, or issue in the project or its archive, or whose `absorbed_by` names no file; and `digest.secret` for an entry whose text matches the shape of a key or token, so that an entry says where a secret lives and never its value.
- [ ] `doc_search` and `flai doc search` return an entry's section, and the dashboard's documents page lists `design/digest/` and renders an entry, shown by a test on each side.
- [ ] `design/digest/README.md` in the template and here says what the folder is for and that only flai and the digest's keeper write in it; `design/system/repository-layout.md`, `design/system/project-manifest.md`, `design/system/flai-cli.md`, `docs/users/flai.md` with `flai-reference.md` regenerated, `docs/operators/settings.md`, and `template/CHANGELOG.md` describe the folder, its key, and the commands.

## Tasks
- T-1807 An ADR records the digest: its folder under design, an entry's file name, ID, front matter, and body, and that nothing in it is deleted
- T-1808 flai/internal/digest reads, validates, and writes digest entries, and refuses one whose text has a secret's shape
- T-1809 flai check reports digest.entry and digest.secret on the entries under design/digest
- T-1810 flai digest new, edit, list, and show write and print digest entries, as text and JSON
- T-1811 doc_search finds a digest entry and the dashboard's documents tree lists design/digest, shown by tests
- T-1812 The digest folder's README, the layout, the CLI design, and the user docs say what the digest is, where it lives, and how flai writes it, and the template ships the folder

## Notes

- Drafted by the planner for E-0025. First story of the epic: the format every later story reads and writes.
- TH-0476 questions 2 and 5 shape the ADR: the folder, one file per entry, and who counts as the operator. The story builds on the recommended answers until they are answered otherwise.
- S-0426 (research, ready) surveys agent memory systems for flai. Its findings may bear on the format and the retrieval; the ADR should cite them if S-0426 is accepted first.
- `cited` is filled by S-0428's priming rule (an agent cites an entry's ID when it relies on it) and read by the curator's compaction (S-0431); `asked_again` is filled by the curator and read by the metric (S-0433).

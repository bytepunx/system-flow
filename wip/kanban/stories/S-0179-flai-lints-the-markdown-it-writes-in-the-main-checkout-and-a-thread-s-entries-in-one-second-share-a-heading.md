---
id: S-0179
type: story
nature: remediation
title: flai lints the markdown it writes in the main checkout, and a thread's entries in one second share a heading
status: in-progress
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:23:14Z
transitions:
  - to: ready
    at: 2026-10-01T08:12:15Z
    by: alex
  - to: in-progress
    at: 2026-10-01T08:19:11Z
    by: claude-opus-5-5
tags: [flai]
touches: [flai/internal/threads, flai/internal/check, flai/internal/workitem, flai/internal/itemnew, flai/internal/itemedit, scripts/lint-md.sh]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 310
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 68
      output: 451
      cache_read: 2501433
      cache_write: 108457
      cost: 1.0667
---
# S-0179 flai lints the markdown it writes in the main checkout, and a thread's entries in one second share a heading

## Goal

Work items, threads, and narratives are written by flai in the main checkout (ADR-0019), so the markdown lint a story runs in its worktree never sees them, and `flai check` does not lint markdown. They reach `main` unlinted and turn CI red: I-0027 has ten occurrences (MD004, MD009, MD012, MD024, MD026, MD029, MD036, MD037), each costing the next stories' smoke runs. The most frequent cause is I-0043: `threads.Reply` and `threads.Resolve` (`flai/internal/threads/threads.go` ~374, ~398) always append a new `### <stamp> <author>` heading, so a reply and a resolution in the same second produce two identical headings (MD024), as in TH-0035. Narratives (`workitem/narrative.go`) and issues (`issues/issues.go`) already merge same-second entries; threads do not.

## Acceptance criteria
- [ ] A thread's reply or resolution in the same second, by the same author, as its last entry is appended to that entry instead of under a new heading, as narratives and issues do; tests cover reply then resolve in one second
- [ ] What flai writes in `wip/` (stories, tasks, epics, threads, narratives, the board) cannot produce markdown the project's lint rejects: titles are written without trailing punctuation in headings (MD026), and bodies given to `item_new`, `item_edit`, `flai story new`, `flai task new`, and `flai edit` are checked against the project's markdownlint rules where they can be, refused with the rule and line otherwise
- [ ] `flai check` reports markdown lint findings in `wip/` (as warnings, so `--strict` fails on them) when the project has a markdownlint configuration, so a story's agent and the operator see them before the commit reaches `main`
- [ ] The duplicate headings already on `main` are gone and `scripts/lint-md.sh` passes on `main`
- [ ] The design (`design/system/flai-cli.md`) and the conventions say where flai lints what it writes
- [ ] I-0027 and I-0043 are closed with what fixed them

## Tasks
- T-0632 A thread's entries in one second by one author share a heading
- T-0633 flai has a markdown linter for what it writes, configured by the project's markdownlint file
- T-0634 flai check reports markdown lint findings in wip as warnings
- T-0635 What flai writes in wip is refused or cleaned before it breaks the lint
- T-0636 The design, conventions, and user guide say where flai lints, and the issues are closed

## Notes

- Whether flai embeds a markdown linter or runs the project's own (`scripts/lint-md.sh`, `markdownlint-cli2`) is for the story to decide; the template's projects may have neither installed, so a missing linter must not fail `flai check`.
- TH-0035's duplicate heading was fixed by hand on 2026-10-01.

---
id: S-0432
type: story
nature: feature
title: The curator runs after each acceptance and on a schedule, marks an entry the design now records absorbed, and compacts an entry no one has cited within the window
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:59:39Z
updated: 2026-10-16T06:07:18Z
transitions: []
tags: [cli, template]
topics: [cli, planning, analysis, conventions]
after: [S-0431]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0432 The curator runs after each acceptance and on a schedule, marks an entry the design now records absorbed, and compacts an entry no one has cited within the window

## Goal

The digest stays small and true without the operator asking. With the `curate` action on, `flai serve` starts the curator after each acceptance, coalescing acceptances that come while a run is under way, and when `digest.schedule` comes round. Each run marks absorbed, with `flai digest absorb`, each entry a design document, a convention, or an ADR now records, naming that document, and so drops it from priming. `flai digest compact` turns each active entry that no narrative, thread, or work item has cited within `digest.compact_after` (TH-0476 question 4) into a one-line entry. Nothing is deleted: the state and the body change, and git keeps the rest.

## Acceptance criteria

- [ ] An ADR, refining the curator's, records the triggers (each acceptance, coalesced; `digest.schedule`, read as `analysis.schedule` is), how `cited` is found (an entry's ID in a narrative, a thread, or a work item, archived ones included, with the date of the latest), absorption, and compaction, and that neither deletes anything.
- [ ] `flai serve` starts the curator after an acceptance while `curate` is on and none runs, and once more when acceptances came during a run; `digest.schedule` starts it with the trigger `schedule <spec>`; each start is journalled with its trigger, and tests cover an acceptance, two during a run, and the schedule.
- [ ] `flai digest absorb <D-nnnn> --by <path> [--heading]` sets `status: absorbed` and `absorbed_by`, refusing a path that does not exist; an absorbed entry never loads in a pack (S-0430) and `flai digest list` shows it with the document.
- [ ] `flai digest cited [--json]` gives each entry its latest citation, and `flai digest compact [--dry-run] [--json]` turns each active entry uncited within `digest.compact_after` into `compacted`, its body reduced to its title, and lists what it compacted; a second run compacts nothing more; tests cover a cited entry kept and an uncited one compacted.
- [ ] The curator's prompt and `.claude/agents/curator.md` tell it to look for absorption in every run (with `doc_search` on each active entry's title) and to run `flai digest compact`; the guard passes `flai digest absorb`, `cited`, and `compact` to it.
- [ ] `digest.schedule` and `digest.compact_after` are manifest settings with checks (`manifest.digest`), the template's manifest documents them, and `docs/operators/settings.md`, `design/system/project-manifest.md`, `design/system/strategic-agents.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` describe the triggers, absorption, and compaction.

## Tasks
- T-1831 An ADR records when the curator runs on its own, how a citation is found, and how absorption and compaction change an entry without deleting it
- T-1832 digest.schedule and digest.compact_after are manifest settings with their check, and the template's manifest documents them
- T-1833 flai digest absorb, cited, and compact mark an entry absorbed, find each entry's latest citation, and compact the entries uncited within the window
- T-1834 flai serve starts the curator after each acceptance, coalesced, and when digest.schedule comes round, while the curate action is on
- T-1835 The curator's prompt, definition, and guard take in absorption and compaction
- T-1836 The design and the user and operator docs describe the curator's triggers, digest.schedule, digest.compact_after, absorption, and compaction

## Notes

- Drafted by the planner for E-0025, after S-0431, which gives the curator its run.
- The run after acceptance follows how `planning.replan` starts the planner when work ahead completes (ADR-0084); reuse its coalescing.

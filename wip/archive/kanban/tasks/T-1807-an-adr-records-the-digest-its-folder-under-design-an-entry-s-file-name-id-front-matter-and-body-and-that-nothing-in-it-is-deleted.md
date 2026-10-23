---
id: T-1807
type: task
nature: feature
title: "An ADR records the digest: its folder under design, an entry's file name, ID, front matter, and body, and that nothing in it is deleted"
status: done
parent: S-0428
owner: alex
created: 2026-10-16T06:02:34Z
updated: 2026-10-23T16:38:30Z
transitions:
  - to: ready
    at: 2026-10-23T16:37:31Z
    by: agent-S-0428
  - to: in-progress
    at: 2026-10-23T16:37:32Z
    by: agent-S-0428
  - to: done
    at: 2026-10-23T16:38:30Z
    by: agent-S-0428
stream: S-0428
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
usage:
  source: log
  seconds: 58
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 7794
      cache_read: 1049632
      cache_write: 39908
      cost: 0.6009
---
# T-1807 An ADR records the digest: its folder under design, an entry's file name, ID, front matter, and body, and that nothing in it is deleted

## Work

Write the ADR with `flai adr new`, on TH-0476's answers to questions 2 and 5. It records:

- The folder, `design/digest/`, following `layout.design` as `design/analysis` does (`flai/internal/analysis/analysis.go`), with no manifest key of its own.
- An entry's file name, `D-nnnn-<slug>.md`, and its ID, allocated as work item IDs are (ADR-0017).
- Each front matter field and its meaning: `id`, `title` (the fact in one sentence), `kind` (`decision`, `answer`, `recommendation`, `cause`), `sources`, `topics`, `paths`, `date`, `status` (`active`, `absorbed`, `compacted`), `absorbed_by`, `asked_again`, and `cited`.
- The body: the fact's detail and the reason, a few lines, linking its sources.
- Who writes: flai alone, through `flai digest`, for the seed, the curator, and the operator's answers; no hand edits.
- That nothing is deleted: absorption and compaction change an entry's state and body, and git keeps the rest.
- That an entry never carries a secret's value, only where it lives, and `flai check` refuses one that looks like a key.

Cite S-0426's findings if it was accepted first. First layer: every other task reads the format this decides.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`, and `flai check --strict` reports nothing on it.

## Notes

The ADR's file name is not known until `flai adr new` makes it, so the task touches the shared `design/adrs` folder.

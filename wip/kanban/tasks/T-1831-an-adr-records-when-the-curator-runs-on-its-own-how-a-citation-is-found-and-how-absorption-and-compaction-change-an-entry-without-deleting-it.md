---
id: T-1831
type: task
nature: feature
title: An ADR records when the curator runs on its own, how a citation is found, and how absorption and compaction change an entry without deleting it
status: done
parent: S-0432
owner: alex
created: 2026-10-16T06:06:38Z
updated: 2026-10-24T13:36:19Z
transitions:
  - to: ready
    at: 2026-10-24T13:35:18Z
    by: agent-S-0432
  - to: in-progress
    at: 2026-10-24T13:35:19Z
    by: agent-S-0432
  - to: done
    at: 2026-10-24T13:36:19Z
    by: agent-S-0432
stream: S-0432
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
usage:
  source: log
  seconds: 60
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 56
      cache_read: 885155
      cache_write: 7228
      cost: 0.3002
---
# T-1831 An ADR records when the curator runs on its own, how a citation is found, and how absorption and compaction change an entry without deleting it

## Work

Write the ADR with `flai adr new`, refining the curator's ADR and on TH-0476's answer to question 4. It records:

- The triggers: each acceptance while `curate` is on, coalesced into one more run when acceptances come during a run; `digest.schedule`, read as `analysis.schedule` is; each with its trigger in the run's record.
- A citation: an entry's ID written in a narrative, a thread, or a work item, archived ones included, dated by the file's commit or the entry's time; `cited` holds the latest.
- Absorption: `status: absorbed` with `absorbed_by` a path (and heading) that exists; never primed.
- Compaction: an active entry with no citation within `digest.compact_after` (default 30 days) becomes `compacted`, its body reduced to its title; still primed as its title.
- That neither deletes a file, and git keeps every earlier body.

First layer: every code task builds what this decides.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`, and `flai check --strict` reports nothing on it.

## Notes

The ADR's file name is not known until `flai adr new` makes it, so the task touches the shared `design/adrs` folder.

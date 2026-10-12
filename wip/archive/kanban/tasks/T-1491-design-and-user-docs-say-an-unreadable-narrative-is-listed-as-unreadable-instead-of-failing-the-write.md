---
id: T-1491
type: task
nature: improvement
title: Design and user docs say an unreadable narrative is listed as unreadable instead of failing the write
status: done
parent: S-0373
owner: alex
created: 2026-10-12T17:29:01Z
updated: 2026-10-12T18:17:26Z
transitions:
  - to: ready
    at: 2026-10-12T18:16:57Z
    by: agent-S-0373
  - to: in-progress
    at: 2026-10-12T18:16:57Z
    by: agent-S-0373
  - to: done
    at: 2026-10-12T18:17:26Z
    by: agent-S-0373
stream: S-0373
tags: [docs]
touches: [design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1489, T-1490]
usage:
  source: log
  seconds: 29
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1758
      cache_read: 680500
      cache_write: 19741
      cost: 0.3292
---
# T-1491 Design and user docs say an unreadable narrative is listed as unreadable instead of failing the write

## Work

Say what T-1489 and T-1490 built, in each place that describes the index:

- `design/system/agent-narrative.md` § Files: `index.md` lists a narrative that does not parse as unreadable, and `flai check` names it with `narrative.front-matter`. No other story's write fails on it. Bump `updated`.
- `design/system/flai-cli.md`, the `flai stream open`, `flai stream log` row: `index.md` is regenerated from the active narratives, an unreadable one listed as such.
- `docs/users/flai.md` § Agent narratives, after "`wip/agents/index.md` is regenerated after every open, log, state, move, and archive": the same, in a sentence a user can act on.

It waits for T-1489 and T-1490, so that it describes the row and the check as they were built.

## Done when

- The three documents say it, each in one or two sentences, with no paragraph copied between them.
- `flai test` passes on the three paths, which runs the markdown lint.

## Notes

`docs/users/flai.md` and `design/system/flai-cli.md` are in the manifest's shared claims, so other open stories editing them do not hold this one.

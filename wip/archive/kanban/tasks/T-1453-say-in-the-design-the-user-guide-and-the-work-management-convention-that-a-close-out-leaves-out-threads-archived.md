---
id: T-1453
type: task
nature: improvement
title: Say in the design, the user guide, and the work-management convention that a close-out leaves out threads.archived
status: done
parent: S-0366
owner: alex
created: 2026-10-12T14:59:47Z
updated: 2026-10-12T15:35:23Z
transitions:
  - to: ready
    at: 2026-10-12T15:27:27Z
    by: agent-S-0366
  - to: in-progress
    at: 2026-10-12T15:27:28Z
    by: agent-S-0366
  - to: done
    at: 2026-10-12T15:35:23Z
    by: agent-S-0366
stream: S-0366
tags: [flai, check, docs]
touches: [design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-1451]
usage:
  source: log
  seconds: 475
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 10287
      cache_read: 1705428
      cache_write: 79225
      cost: 1.0827
---
# T-1453 Say in the design, the user guide, and the work-management convention that a close-out leaves out threads.archived

## Work

Say where each lists what a close-out's check leaves out that a `threads.archived` outside the story is left out too, citing T-1451's ADR:

- the `flai check` row of `design/system/flai-cli.md`
- the close-out item in `design/system/continuous-improvement.md` § the findings a close-out meets outside the story
- the `--story` and `--record-issues` paragraphs of `docs/users/flai.md`
- the close-out rule in `design/conventions/work-management.md`, and the same rule in `template/root/design/conventions/work-management.md`, worded alike: a thread still open on an archived item is the operator's or its opener's to resolve, and only `flai thread resolve` in the main checkout clears it
- an entry for S-0366 in `template/CHANGELOG.md`, patch, naming the ADR, as S-0348's entry does

Follow S-0365's wording for `issues.no-story` where it has landed, so the lists read alike. The convention's baseline changes, so it is changed in `template/` in the same story, as `CLAUDE.md` says.

It waits for T-1451 for the ADR's number, and runs beside T-1452: they share no path.

## Done when

- Each place above names `threads.archived` among what a close-out leaves out, with the ADR linked.
- The convention reads the same in `design/conventions/` and `template/root/design/conventions/`.
- `flai test` on the changed paths passes, the markdown lint included.

## Notes

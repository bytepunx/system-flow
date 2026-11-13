---
id: T-3271
type: task
nature: improvement
title: Say in the design, the user guide, and the work-management convention that a close-out leaves out item.blocked-closed
status: done
parent: S-0761
owner: alex
created: 2026-11-13T07:29:51Z
updated: 2026-11-13T08:58:22Z
transitions:
  - to: ready
    at: 2026-11-13T08:24:49Z
    by: agent-S-0761
  - to: in-progress
    at: 2026-11-13T08:24:49Z
    by: agent-S-0761
  - to: done
    at: 2026-11-13T08:58:22Z
    by: agent-S-0761
stream: S-0761
tags: [flai, check, docs]
touches: [design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-3269]
usage:
  source: log
  seconds: 2013
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 2862
      cache_read: 1194374
      cache_write: 50542
      cost: 0.616
---
# T-3271 Say in the design, the user guide, and the work-management convention that a close-out leaves out item.blocked-closed

## Work

Each place below lists what a close-out's check leaves out. In each, say that an `item.blocked-closed` outside the story is left out too, and cite T-3269's ADR:

- the `flai check` row of `design/system/flai-cli.md`
- the close-out item in `design/system/continuous-improvement.md`, on the findings a close-out meets outside the story
- the `--story` and `--record-issues` paragraphs of `docs/users/flai.md`
- the close-out rule in `design/conventions/work-management.md`, and the same rule in `template/root/design/conventions/work-management.md`, worded alike. Say there that only `flai unblock` on that item, in the main checkout, clears it, as the rule says `flai archive` clears an `item.archive`.
- an entry for S-0761 in `template/CHANGELOG.md`, as a patch, naming the ADR, as S-0555's entry does

Put it beside the `item.archive` text where a list has one, so the two item findings read alike. The convention's baseline changes, so `template/` changes in the same story, as `CLAUDE.md` says.

It waits for T-3269 for the ADR's number and decision. It runs beside T-3270, since they share no path.

## Done when

- Each place above names `item.blocked-closed` among what a close-out leaves out, with the ADR linked.
- The convention reads the same in `design/conventions/` and `template/root/design/conventions/`.
- `flai test` on the changed paths passes, the markdown lint included.

## Notes

---
id: T-3349
type: task
nature: remediation
title: Say in the design, the user guide, and the work-management convention that a close-out leaves out story.unaccepted
status: done
parent: S-0781
owner: alex
created: 2026-11-13T22:54:04Z
updated: 2026-11-13T23:28:17Z
transitions:
  - to: ready
    at: 2026-11-13T23:10:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T23:11:00Z
    by: system-flow
  - to: done
    at: 2026-11-13T23:28:17Z
    by: agent-S-0781
stream: S-0781
tags: [docs, conventions, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md]
after: [T-3347]
usage:
  source: log
  seconds: 1037
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 12440
      cache_read: 2073462
      cache_write: 92535
      cost: 1.25
---
# T-3349 Say in the design, the user guide, and the work-management convention that a close-out leaves out story.unaccepted

## Work

As the ADR of T-3347 decides. It waits for T-3347 for that decision and the ADR's number, which each change links. If the ADR decides the root fix alone, ask on the story's thread to cancel this task.

Add `story.unaccepted` outside the story to each list of what a check scoped to a story leaves out, as S-0761 added `item.blocked-closed`:

- `design/conventions/work-management.md`, the close-out rule, and its baseline copy `template/root/design/conventions/work-management.md`, with a line in `template/CHANGELOG.md` under `## Unreleased`. Say only the operator clears it, with `flai accept` or by deleting the branch, in the main checkout.
- `design/system/flai-cli.md`, the `flai check` row.
- `design/system/continuous-improvement.md`, the findings a close-out records.
- `docs/users/flai.md`, the `--story` and `--record-issues` passage. Its example output at about line 2227 shows a `story.unaccepted` recorded outside S-0249. That is now left out, so give it another rule that a close-out still records.

It runs beside T-3348: the two share no file.

## Done when

- Each of the six files says a scoped check leaves out `story.unaccepted` outside the story, linking the ADR.
- The guide's example no longer shows a `story.unaccepted` recorded outside a story.
- The markdown lint passes on the changed files.

## Notes

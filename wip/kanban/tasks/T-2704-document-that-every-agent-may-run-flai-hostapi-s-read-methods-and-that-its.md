---
id: T-2704
type: task
nature: improvement
title: Document that every agent may run flai hostapi's read methods and that its write methods are the story's agent's
status: in-progress
parent: S-0612
owner: alex
created: 2026-11-09T07:57:58Z
updated: 2026-11-09T08:38:47Z
transitions:
  - to: ready
    at: 2026-11-09T08:38:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T08:38:47Z
    by: system-flow
stream: S-0612
tags: [docs, guard]
touches: [docs/users/flai.md, design/system/flai-cli.md]
after: [T-2702]
---
# T-2704 Document that every agent may run flai hostapi's read methods and that its write methods are the story's agent's

## Work

Waits for T-2702, whose guard rule and read list this task describes. Runs beside T-2703, with no path in common.

- `docs/users/flai.md` § Request timing (S-0152): say that a story's sub-agents, read-only ones among them, and the strategic agents' own calls may run `flai hostapi` alone and with a read method, `--timing` included, and that `flai guard` refuses them every other method, naming it (S-0612, I-0307).
- `design/system/flai-cli.md`, the `flai guard` row: add `hostapi` to the reads it lists, with the form that reads (no method, or a method on the guard's list), where the list lives, and the test in `flai/internal/hostapi` that keeps it in step with the method table.

No ADR: allowing a read is within ADR-0060, as `flai role list` (S-0568) and `flai subproject list` (S-0491) were.

## Done when

- Both documents say what the guard now allows and refuses for `flai hostapi`, and `flai test docs/users/flai.md design/system/flai-cli.md` passes the markdown lint.

## Notes

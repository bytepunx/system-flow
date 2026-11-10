---
id: T-2848
type: task
nature: improvement
title: Record in an ADR refining ADR-0272, the design, and the user docs that a checked edit holds the checked-write lock
status: in-progress
parent: S-0654
owner: alex
created: 2026-11-10T07:01:24Z
updated: 2026-11-10T07:44:43Z
transitions:
  - to: ready
    at: 2026-11-10T07:44:38Z
    by: agent-S-0654
  - to: in-progress
    at: 2026-11-10T07:44:43Z
    by: agent-S-0654
stream: S-0654
tags: [flai, docs, adr]
touches: [design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md]
after: [T-2844, T-2847]
---
# T-2848 Record in an ADR refining ADR-0272, the design, and the user docs that a checked edit holds the checked-write lock

## Work

Write a new ADR with `flai adr new --refines ADR-0272` (with `--body-stdin`): a checked edit of a work item, `itemedit.Apply` and `itemedit.ApplyForecasts`, holds the checked-write lock from its first save until it keeps or restores the edit, as a checked creation does, so flai serve's look never judges a claim by touches a refused edit wrote. Its Context names I-0328 and ADR-0272's consequence that left edits out; its Consequences say a look waits at most as long as an edit's check runs, that two checked writes take turns, and that `flai touches`, which saves with no check and no restore, has no window and takes no lock. `flai adr new` writes the ADR file and `design/adrs/README.md`; ADR-0272 itself is accepted and stays as it is.

Then say the same where the design and the docs describe the checked edit, linking the new ADR: the `flai edit` row of `design/system/flai-cli.md`, the **Held** item of `design/system/workflow.md` where it says what the look reads under the lock, and the `flai edit` section of `docs/users/flai.md`, which says a refused edit puts every file back.

It waits for T-2844 and T-2847: the ADR and the docs describe what those two tasks built, and name the functions as they left them.

## Done when

- The new ADR is accepted, refines ADR-0272, and is in `design/adrs/README.md`.
- `design/system/flai-cli.md`, `design/system/workflow.md`, and `docs/users/flai.md` say a checked edit holds the checked-write lock and link the ADR.
- `flai check --strict` and the markdown lint pass on the files written.

## Notes

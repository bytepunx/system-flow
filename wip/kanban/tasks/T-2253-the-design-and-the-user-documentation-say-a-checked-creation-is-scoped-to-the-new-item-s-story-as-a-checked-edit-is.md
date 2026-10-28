---
id: T-2253
type: task
nature: remediation
title: The design and the user documentation say a checked creation is scoped to the new item's story as a checked edit is
status: backlog
parent: S-0525
owner: alex
created: 2026-10-28T14:52:55Z
updated: 2026-10-28T14:52:55Z
transitions: []
stream: S-0525
tags: [flai, docs]
touches: [design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md]
after: [T-2251, T-2252]
---
# T-2253 The design and the user documentation say a checked creation is scoped to the new item's story as a checked edit is

## Work

Bring the documents that describe checked creation in step with T-2251's code and T-2252's ADR, linking the ADR by its relative path:

- `design/system/flai-cli.md`: the `flai story new` and `flai task new` rows (the `--body-stdin`, `--after`, and `--autocommit` creation through `itemnew.Create`) say the second check run is scoped to the new item's story as `flai edit`'s is, a new finding outside it is a warning in the result (`--json` lists it in `warnings`), and one inside it, or an error on the item's own file, refuses it and leaves nothing; the `flai issue story` row and the `item_new` tool follow the same creation, so say so where they describe the refusal.
- `design/system/strategic-agents.md` § Checked creation: the same for the planner's `item_new`, in one sentence, naming the ADR.
- `docs/users/flai.md`: the paragraph on `flai story new --body-stdin` and checked creation says what a user sees: a finding the creation brings on the item refuses it; one elsewhere in the repository, such as a `.flai-cache/` not yet ignored, warns and the item is made.

Bump `updated` on each design document. It waits for T-2251 and T-2252 so that it describes the behaviour and the ADR as they landed, not as they were planned.

## Done when

- `design/system/flai-cli.md`, `design/system/strategic-agents.md` § Checked creation, and `docs/users/flai.md` describe the scoped creation and link the new ADR.
- `flai test` on the three files passes the markdown lint, and `flai check --strict` is clean on them.

## Notes

`docs/users/flaiover.md` describes the dashboard's checked edit at its save paragraph and names no creation refusal; change it only if it does.

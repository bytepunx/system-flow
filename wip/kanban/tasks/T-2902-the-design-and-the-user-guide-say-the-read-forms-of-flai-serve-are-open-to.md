---
id: T-2902
type: task
nature: remediation
title: The design and the user guide say the read forms of flai serve are open to every agent
status: backlog
parent: S-0680
owner: alex
created: 2026-11-10T08:46:51Z
updated: 2026-11-10T08:46:51Z
transitions: []
stream: S-0680
tags: [docs, guard]
touches: [design/system/strategic-agents.md, docs/users/flai.md]
after: [T-2901]
---
# T-2902 The design and the user guide say the read forms of flai serve are open to every agent

## Work

Write down what T-2901 changed, in the words of the forms it passes.

- `design/system/strategic-agents.md` § The orchestrator › Its permissions and the guard: replace the sentence saying `flai serve status` is still refused (I-0318) with one saying the guard counts the read forms of `flai serve` among the reads, with or without a redirection (S-0680), and add `flai serve status` and the other read forms to the `none needed` row of the permissions table.
- `docs/users/flai.md`, the paragraph on what a sub-agent may run, near "A sub-agent may read and run checks": add that it may read the host with `flai serve status`, `actions`, `journal`, and the show and list forms, and is refused every form that starts, stops, or writes (ADR-0220).

Waits for T-2901: it describes the forms that task's `readsServe` settles.

## Done when

- Both documents name the read forms of `flai serve` T-2901 passes, and neither says `flai serve status` is refused.
- The markdown lint passes on both files.

## Notes

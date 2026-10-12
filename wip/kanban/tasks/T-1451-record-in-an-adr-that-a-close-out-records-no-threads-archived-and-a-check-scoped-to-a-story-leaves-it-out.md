---
id: T-1451
type: task
nature: improvement
title: Record in an ADR that a close-out records no threads.archived and a check scoped to a story leaves it out
status: backlog
parent: S-0366
owner: alex
created: 2026-10-12T14:59:32Z
updated: 2026-10-12T14:59:32Z
transitions: []
stream: S-0366
tags: [flai, check, adr]
touches: [design/adrs, design/adrs/README.md]
---
# T-1451 Record in an ADR that a close-out records no threads.archived and a check scoped to a story leaves it out

## Work

Write and accept an ADR with `flai adr new --commit`. It refines ADR-0085 and follows ADR-0122, ADR-0133, and the ADR S-0365 writes for `issues.no-story`. Its decision: a `threads.archived` finding outside the story is left out of a check scoped to a story (`flai check --story`). So a close-out neither notes it nor records it in an issue. One on a thread anchored on the story or one of its tasks stays inside, as `ScopeToStory` already counts such a thread.

Its context is I-0135. Three close-outs (S-0343, S-0289, S-0342) found `threads.archived` on TH-0389, a thread none of them touched. The orchestrator opened TH-0389 on S-0338 at 2026-10-08T20:55:51Z, after S-0338's acceptance had archived it. ADR-0109 resolves the threads on an item when it is archived, so it could not resolve a thread opened later. The thread was open for a reason: it awaited the operator. Threads live in `wip/` on the main branch, so no story branch can resolve one. Only the operator or the thread's opener can. Two places still put it before them, and the ADR's consequences name them:

- `flai check` in the main checkout still warns `threads.archived`
- `inbox` and the dashboard list the open thread

The alternative is to stop `threads.archived` firing on a thread opened after its item was archived, in every check. That fixes the cause of TH-0389's instances but not a thread left open by an archive that stopped part way. Decide it instead if, on reading the instances, it fits better, and say why in the ADR.

It waits for nothing in this story. S-0366 waits for S-0365, so this ADR's number follows S-0365's.

## Done when

- The ADR is accepted under `design/adrs/`, indexed in `design/adrs/README.md`, with topics set.
- It says whether a `threads.archived` outside the story is left out, and why, and what still reports it.

## Notes

---
id: T-2936
type: task
nature: improvement
title: "ADR: the touches drift leaves out the folders flai writes, and flai task done narrows a closing task's unchanged touches"
status: backlog
parent: S-0674
owner: alex
created: 2026-11-10T08:49:51Z
updated: 2026-11-10T08:49:51Z
transitions: []
stream: S-0674
tags: [flai, adr]
touches: [design/adrs/README.md]
---
# T-2936 ADR: the touches drift leaves out the folders flai writes, and flai task done narrows a closing task's unchanged touches

## Work

Write one ADR with `flai adr new`, refining ADR-0144 and ADR-0113, and accept it. It records two decisions. The story's `### Planning` note gives the evidence: 80 of 115 recent stories had `outside` files, and all of those files are under folders flai writes.

1. A story's committed files in `claims.drift[]` leave out the files under the `digest`, `analysis`, and `issues` folders of the manifest's design folder, as they already leave out `wip/` (S-0205). flai's own commits write those folders: the acceptance commit's sweep of the main checkout (ADR-0248), the issue step's `create story`, and the close-out's `record I-nnnn`. The drift's touches side leaves out the same folders and `wip/`, so that a touch there is never counted `unchanged`.
2. `flai task done`'s touches step also narrows. It drops from the closing task's touches each file the branch does not change against its merge base with main. It drops that file from the story's touches when no other open task names it. It runs even when the task committed nothing, and it never drops a folder touch or a path the branch changes.

Say in Consequences that `exact_share` then measures the claim as task done leaves it, not the planner's first prediction.

This task waits for no other: the code tasks build on its decision.

## Done when

- The ADR is accepted, refines ADR-0144 and ADR-0113, and names both decisions with the evidence figures.
- `design/adrs/README.md` lists it, as `flai adr new` and `flai adr index` write it.
- `flai check --strict` is clean for the story.

## Notes

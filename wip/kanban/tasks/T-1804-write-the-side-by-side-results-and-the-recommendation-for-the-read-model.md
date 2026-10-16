---
id: T-1804
type: task
nature: experiment
title: Write the side-by-side results and the recommendation for the read model
status: backlog
parent: S-0427
owner: alex
created: 2026-10-16T05:59:44Z
updated: 2026-10-16T05:59:44Z
transitions: []
stream: S-0427
tags: [cli]
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
after: [T-1803]
---
# T-1804 Write the side-by-side results and the recommendation for the read model

## Work

Finish the results document from the figures of T-1800 to T-1803:

- `## Results`: one table, the file reads, SQLite, and Postgres where it ran, side by side, for each query cold and warm, the load and the rescans, storage, RAM, and CPU. Below it, what the savings mean in the operator's wait for a chart or the board and in an agent's one-shot `flai stats` or prime call.
- `## Recommendation`: adopt, adapt, or drop, with its reason. Name the follow-up it implies, such as a story to put a read model behind `statsread.Read`, as a title and a one-line goal, and the decisions the operator must make first, such as where the database lives and whether flai takes a SQLite dependency. Write no ADR and no story. Say how it bears on S-0426's recommendation on recall and search, where that story has one.

Run `flai check --strict` and the markdown lint (`flai test design/experiments/S-0427-sqlite-and-postgres-read-models.md`) on the document and fix what they find.

Waits for T-1803: it summarises every figure, and all four write the results document.

## Done when

- `## Results` has the side-by-side table and what the savings mean for the operator and for agents.
- `## Recommendation` says adopt, adapt, or drop, with its reason, the follow-up, and the operator's decisions.
- `flai check --strict` and the markdown lint pass on the document.

## Notes

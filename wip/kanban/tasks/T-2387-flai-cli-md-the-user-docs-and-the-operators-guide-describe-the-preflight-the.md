---
id: T-2387
type: task
nature: improvement
title: flai-cli.md, the user docs, and the operators guide describe the preflight, the held state, its thread, and the landing again without verification
status: backlog
parent: S-0549
owner: alex
created: 2026-10-31T07:42:06Z
updated: 2026-10-31T07:42:06Z
transitions: []
stream: S-0549
tags: [flai]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md]
after: [T-2383, T-2384, T-2386]
---
# T-2387 flai-cli.md, the user docs, and the operators guide describe the preflight, the held state, its thread, and the landing again without verification

## Work

Describe what the code tasks built, linking the new ADR:

- `design/system/flai-cli.md`, `### The merge queue pass`: a preflight step before the build, the preflight hold apart from the operator's hold (the table row for the hold and the bullets on the build, the hold, and landing), the thread, how the hold clears, and the landing again from verified commits.
- `docs/users/flai.md`, `### The merge queue`: the preflight under `#### The batch`, the preflight hold beside `#### Holding the queue`, and the landing again under `#### Landing and sending back`, with what `flai board --json` shows under `queue`.
- `docs/operators/index.md`, the merge-queue host action section: rewrite `**When a landing fails.**`, which describes today's drop and rebuild, and say what the operator does on the thread: fix or commit the file the thread names, and the queue resumes by itself.

It waits for the code tasks, the preflight, the landing again, and the serve task, since it describes what they built. The dashboard task documents the panel itself.

## Done when

The three documents describe the preflight, the held state, its thread, and the landing again as built, link the ADR, and pass the markdown lint and `flai check --strict`. Story criterion 4 (`flai-cli.md` and the user docs).

## Notes

Drafted by the planner of S-0549 on 2026-10-31.

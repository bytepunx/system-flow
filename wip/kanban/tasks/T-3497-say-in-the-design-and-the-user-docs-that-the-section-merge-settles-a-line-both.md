---
id: T-3497
type: task
nature: improvement
title: Say in the design and the user docs that the section merge settles a line both sides change in different words
status: backlog
parent: S-0817
owner: alex
created: 2026-11-17T02:31:46Z
updated: 2026-11-17T02:31:46Z
transitions: []
stream: S-0817
tags: [docs, merge]
touches: [design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-3495]
---
# T-3497 Say in the design and the user docs that the section merge settles a line both sides change in different words

## Work

Where each document says what the section merge settles and where it stops, add the word merge and link the ADR T-3495 writes:

- `design/system/workflow.md` § Branches and collisions (ADR-0019): the S-0401 paragraph on ADR-0175, ADR-0276, and ADR-0287.
- `design/system/agent-coordination.md`: the ADR-0175 paragraph (S-0401).
- `design/system/flai-cli.md`: the `flai stream sync` and `flai accept` rows of the commands table, where they name the same-lines stop.
- `docs/users/flai.md`: the Shared documents bullet under `flai stream sync`. Its sentence "When both sides changed one sentence, the same table row, heading, or line in a fenced code block … it merges nothing and stops" becomes: both changed the same word, or inserted at the same point, of a sentence or a table row; a heading; a line in a fenced code block.

Bump each file's `updated`. Change only those passages: these files are hubs other open stories change too.

Waits for T-3495: the documents link the ADR it writes.

## Done when

- The four documents say a line both sides change in different words is merged, and name what still stops.
- Each links the new ADR, and `flai test` on the four files passes the markdown lint.

## Notes

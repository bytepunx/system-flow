---
id: T-3011
type: task
nature: improvement
title: The design and user docs say a rebase writes the generated files again at the branch's head
status: backlog
parent: S-0702
owner: alex
created: 2026-11-10T10:31:36Z
updated: 2026-11-10T10:35:21Z
transitions: []
stream: S-0702
tags: [docs]
touches: [design/system/continuous-improvement.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md]
after: [T-3002]
---
# T-3011 The design and user docs say a rebase writes the generated files again at the branch's head

## Work

Say what the new ADR decides where the design and the user docs describe the files a rebase writes again:

- `design/system/continuous-improvement.md`: § Summary, which says how the sync writes the summary again and what happens when other paths conflict too.
- `design/system/flai-cli.md`: the `flai stream sync` row.
- `design/system/workflow.md`: § Branches and collisions › Hub files: what S-0550 changed, where it names what sync regenerates (ADR-0098).
- `docs/users/flai.md`: the "The issue summary" bullet and the generated-files paragraph of `flai check`.

Each says that after a rebase that wrote the issues summary or the ADR index at a stop, flai writes them again at the branch's head and commits what changed. Link the new ADR.

It waits for T-3002, whose ADR it describes and links. It shares no file with T-3005 and runs beside it.

## Done when

- The four documents say it and link the ADR.
- The markdown lint passes on them.

## Notes

---
id: T-1841
type: task
nature: feature
title: The design and the user docs say that an operator's answer on a thread becomes a digest entry at once
status: backlog
parent: S-0433
owner: alex
created: 2026-10-16T06:08:01Z
updated: 2026-10-23T17:16:07Z
transitions: []
stream: S-0433
tags: [cli]
touches: [design/system/strategic-agents.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-1838, T-1839, T-1840]
---
# T-1841 The design and the user docs say that an operator's answer on a thread becomes a digest entry at once

## Work

- `design/system/strategic-agents.md`, "The curator": the entries written at once, and that its next run merges them.
- `design/system/agent-narrative.md`: an operator's answer is in the digest at once, so the next agent's pack or `doc_search` finds it.
- `docs/users/flai.md` (`flai thread reply`, `confirm`) and `docs/users/flaiover.md` (answering a thread from the dashboard): the entry written and the `digest` field of the result.

Waits for the three entry point tasks, so the docs describe what was built.

## Done when

The documents say what the code does, and `flai check --strict` and the markdown lint report nothing on them.

## Notes

None of these files is in `claims.shared` but `docs/users/flai.md`, so this task may hold another story while S-0433 runs.

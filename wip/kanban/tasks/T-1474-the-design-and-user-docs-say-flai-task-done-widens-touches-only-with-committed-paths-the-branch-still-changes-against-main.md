---
id: T-1474
type: task
nature: remediation
title: The design and user docs say flai task done widens touches only with committed paths the branch still changes against main
status: backlog
parent: S-0370
owner: alex
created: 2026-10-12T17:10:05Z
updated: 2026-10-12T17:10:05Z
transitions: []
stream: S-0370
tags: [cli, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md]
after: [T-1472]
---
# T-1474 The design and user docs say flai task done widens touches only with committed paths the branch still changes against main

## Work

Say what T-1472 decides where the touches step is described, linking the new ADR.

- In `design/system/flai-cli.md`, the `flai task done` row: the touches clause becomes the paths the commit changed that the branch still changes against main, so a path changed back to main's version is not added. Cite S-0370 and the ADR beside ADR-0107 and ADR-0128.
- In `docs/users/flai.md` § Closing a task, step 6 (Touches) says the same, with one sentence on why: a story that narrowed its touches and reverted a file keeps its claim narrow. Check § Touches' sentence that `flai task done` adds the paths its own commit changed, and adjust it to match.

Waits for T-1472: it links the ADR and describes what it decides. Runs beside T-1473, with no path in common.

## Done when

- Both documents describe the touches step as the ADR decides, and link it with a relative path.
- The markdown lint passes on both files, and `flai check --strict` is clean.

## Notes

Drafted by the planner for S-0370. S-0363 changes the same two documents in other sections; sync over its change if it is accepted first.

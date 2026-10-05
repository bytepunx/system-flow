---
id: T-0873
type: task
nature: remediation
title: flai-cli.md and the users' guide say the acceptance preview reports conflict markers
status: backlog
parent: S-0276
owner: alex
created: 2026-10-05T04:06:50Z
updated: 2026-10-05T04:06:50Z
transitions: []
stream: S-0276
tags: [docs]
touches: [design/system/flai-cli.md, docs/users/flai.md]
after: [T-0872]
---
# T-0873 flai-cli.md and the users' guide say the acceptance preview reports conflict markers

## Work

- In `design/system/flai-cli.md`, add the conflict marker check to the preflight that the `flai accept` row lists, citing S-0276. Say that `preview.Accept` and the acceptance read the markers by the one `conflictmark` scan. If the package layout lists `conflictmark`, say there that it reads a branch too.
- In `docs/users/flai.md`, the acceptance section's paragraph on conflict markers (S-0253) says what acceptance refuses. Add that `flai accept --dry-run`, and the dashboard's Accept dialog that runs it, lists the same markers as a blocker before anything is merged. Quote the blocker as T-0872 words it.

This task waits for T-0872, so that the docs quote the blocker as it is written.

## Done when

- Both documents say that the preview reports a conflict marker as a blocker, matching what T-0872 built.
- `flai check --strict` and the markdown lint are clean on both files.

## Notes

---
id: T-0452
type: task
nature: improvement
title: flai serve project remove adds a project served from below a folder to the removed list, and add takes it off
status: done
parent: S-0123
owner: alex
created: 2026-09-26T07:06:21Z
updated: 2026-09-26T07:11:24Z
transitions:
  - to: ready
    at: 2026-09-26T07:06:29Z
    by: agent-S-0123
  - to: in-progress
    at: 2026-09-26T07:08:51Z
    by: agent-S-0123
  - to: done
    at: 2026-09-26T07:11:24Z
    by: agent-S-0123
stream: S-0123
tags: [cli]
touches: [flai/cmd]
---
# T-0452 flai serve project remove adds a project served from below a folder to the removed list, and add takes it off

## Work

- `flai serve project remove <key|folder>` on a project served from below an import folder or the start folder adds its root to the list instead of refusing; on a registered one that is also below such a folder it unregisters it and adds it, so it is not served from below instead.
- `flai serve project add <folder>` on a root in the list takes it off; when the project is below a folder flai serve serves from it is served from there again without being registered, otherwise it is registered as before.
- Text and `--json` output say which happened; `docs/users/flai.md` reference regenerated if the help changes.

## Done when

- Command tests cover remove and add of a project below an import folder, and remove of a registered one below it; `make test` passes.

## Notes

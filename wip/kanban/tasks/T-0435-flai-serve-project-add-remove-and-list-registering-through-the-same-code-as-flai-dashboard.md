---
id: T-0435
type: task
nature: feature
title: flai serve project add, remove, and list, registering through the same code as flai dashboard
status: done
parent: S-0118
owner: alex
created: 2026-09-26T05:25:59Z
updated: 2026-09-26T05:35:52Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:18Z
    by: agent-S-0118
  - to: in-progress
    at: 2026-09-26T05:30:23Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T05:35:52Z
    by: agent-S-0118
stream: S-0118
tags: []
---

# T-0435 flai serve project add, remove, and list, registering through the same code as flai dashboard

## Work

New `cmd/serve_project.go`: `flai serve project add [dir]` (default the current folder, the main checkout for a worktree) resolves the dashboard address as `flai dashboard` does, makes the agent credential when missing, registers through one function that `flai dashboard`'s `connectServe` and `connectServeFolder` also use, and starts flai host when none runs. It refuses a folder with no system-flow.yaml, a manifest with no key, and a key another project has. `remove <key|dir>` unregisters and touches no file of the project and not the container. `list` (and `--json`) shows each registered project with key, root, dashboard address, and connected, last error, or unavailable reason; the projects served for the folder flai serve started in; the import candidates; and the projects under the import roots that are not served. `flai serve status` shows the unavailable reason too.

## Done when

- cmd tests cover add (default folder, refusals), remove by key and by folder (files untouched), list text and JSON, and status reporting an unavailable root
- `flai dashboard` and `flai dashboard stop` tests still pass unchanged
- `scripts/flai-test.sh` passes

## Notes

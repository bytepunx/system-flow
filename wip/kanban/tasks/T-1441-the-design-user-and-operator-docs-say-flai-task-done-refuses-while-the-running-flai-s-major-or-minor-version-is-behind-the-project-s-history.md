---
id: T-1441
type: task
nature: remediation
title: The design, user, and operator docs say flai task done refuses while the running flai's major or minor version is behind the project's history
status: done
parent: S-0363
owner: alex
created: 2026-10-12T14:36:40Z
updated: 2026-10-12T18:22:56Z
transitions:
  - to: ready
    at: 2026-10-12T18:15:28Z
    by: agent-S-0363
  - to: in-progress
    at: 2026-10-12T18:15:29Z
    by: agent-S-0363
  - to: done
    at: 2026-10-12T18:22:56Z
    by: agent-S-0363
stream: S-0363
tags: [cli, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md]
after: [T-1439]
usage:
  source: log
  seconds: 447
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 275
      cache_read: 1323956
      cache_write: 73392
      cost: 0.6336
---
# T-1441 The design, user, and operator docs say flai task done refuses while the running flai's major or minor version is behind the project's history

## Work

Document the refusal T-1439's ADR decides, in the same story as the code: refused while the running flai's major or minor version is below the newest `flai/vX.Y.Z` in the project's history; a patch behind only warns.

- `design/system/flai-cli.md`, section "Versions: the host's flai and the tree": after the bullet "It says when it is behind", add that `flai task done` and the MCP tool `task_done` refuse, changing nothing, while the running flai's major or minor version is behind, and that a patch behind and every other command still only warn. Link the ADR. Bump `updated`.
- `docs/users/flai.md`, section "Closing a task": say that the close is refused while the running flai is a minor or major version behind the newest flai release in the project's history, and that the refusal names the upgrade command. A story's agent asks the operator to upgrade on a thread. Bump `updated`.
- `docs/operators/index.md`, section "Keeping the host's flai current": say that a host flai a minor version or more behind now stops every story's agent at its next task close until `flai host upgrade` runs. Bump `updated`.

Waits for T-1439, whose decision it describes. It runs beside T-1440: no path in common.

## Done when

- The three sections say what the refusal does, with the wording of the ADR's decision.
- `flai test design/system/flai-cli.md docs/users/flai.md docs/operators/index.md` passes the markdown lint.

## Notes

Drafted by the planner for S-0363; rethresholded by S-0363's agent after TH-0401.

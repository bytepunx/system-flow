---
id: T-1441
type: task
nature: remediation
title: The design, user, and operator docs say flai task done refuses while the running flai is older than the project's history
status: backlog
parent: S-0363
owner: alex
created: 2026-10-12T14:36:40Z
updated: 2026-10-12T14:36:40Z
transitions: []
stream: S-0363
tags: [cli, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md]
after: [T-1439]
---
# T-1441 The design, user, and operator docs say flai task done refuses while the running flai is older than the project's history

## Work

Document the refusal ADR-0136 decides, in the same story as the code.

- `design/system/flai-cli.md`, section "Versions: the host's flai and the tree": after the bullet "It says when it is behind", add that `flai task done` and the MCP tool `task_done` refuse, changing nothing, while `release.FlaiOutdated` reports, and that every other command still only warns. Link ADR-0136. Bump `updated`.
- `docs/users/flai.md`, section "Closing a task": say that the close is refused while the running flai is older than the newest flai release in the project's history, and that the refusal names the upgrade command. A story's agent asks the operator to upgrade on a thread. Bump `updated`.
- `docs/operators/index.md`, section "Keeping the host's flai current": say that an outdated host flai now stops every story's agent at its first task close until `flai host upgrade` runs. Bump `updated`.

Waits for T-1439, whose decision it describes. It runs beside T-1440: no path in common.

## Done when

- The three sections say what the refusal does, with the wording of the ADR's decision.
- `flai test design/system/flai-cli.md docs/users/flai.md docs/operators/index.md` passes the markdown lint.

## Notes

Drafted by the planner for S-0363. S-0337 (in progress) also changes `design/system/flai-cli.md` and `docs/users/flai.md`; sync after its acceptance before you start.

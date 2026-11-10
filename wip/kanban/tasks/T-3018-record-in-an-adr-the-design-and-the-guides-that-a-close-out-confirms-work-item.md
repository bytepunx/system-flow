---
id: T-3018
type: task
nature: remediation
title: Record in an ADR, the design, and the guides that a close-out confirms work-item findings under the checked-write lock
status: backlog
parent: S-0679
owner: alex
created: 2026-11-10T10:32:01Z
updated: 2026-11-10T10:32:01Z
transitions: []
stream: S-0679
tags: [flai, docs]
touches: [design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md]
after: [T-3014, T-3017]
---
# T-3018 Record in an ADR, the design, and the guides that a close-out confirms work-item findings under the checked-write lock

## Work

Write an ADR with `flai adr new` that refines ADR-0272 and ADR-0085: a close-out's record step confirms each finding outside the story on a work item under the checked-write lock and records only what still holds, and the close-out's container mounts the lock file. Say what a lock that cannot be taken does. Accept it with `flai adr accept`, which writes `design/adrs/README.md`; `flai task done` widens the touches with the ADR's own file.

Then say the same where readers look:

- `design/system/flai-cli.md`: the `flai check --record-issues` and `flai verify --record-issues` behaviour.
- `docs/users/flai.md`: the `--record-issues` paragraph under the check section, and the container mounts paragraph beside ADR-0253's (`wip/`, `system-flow.yaml`, `.flai-cache/verify`, and now the lock file).
- `docs/operators/runbooks/agent-container.md`: the mounts the close-out's container takes.

This task waits for T-3014 and T-3017: it records what they built, the lock's use and the mount.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`.
- The design, the user guide, and the runbook say what the record step confirms under the lock and what the container mounts, linking the ADR.
- `flai test` on the changed paths passes.

## Notes

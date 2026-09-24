---
id: T-0391
type: task
nature: remediation
title: The installer smoke test checks self-upgrade's own path on a binary outside any project
status: done
parent: S-0112
owner: alex
created: 2026-09-24T06:25:24Z
updated: 2026-09-24T06:27:22Z
transitions:
  - to: ready
    at: 2026-09-24T06:25:24Z
    by: agent-S-0112
  - to: in-progress
    at: 2026-09-24T06:25:24Z
    by: agent-S-0112
  - to: done
    at: 2026-09-24T06:27:22Z
    by: agent-S-0112
stream: S-0112
tags: []
touches: [scripts/install-test.sh]
---
# T-0391 The installer smoke test checks self-upgrade's own path on a binary outside any project

## Work

- `scripts/install-test.sh` installs into `.flai-cache/install-test`, inside this project, and expects `self-upgrade --check` to name that path. Since flai 1.15.3 (S-0111) a binary inside a project upgrades into `~/.flai/bin`, so the smoke tier fails once 1.15.3 is published. Run the last check on a copy outside any project, and check that the copy inside the project resolves to a scratch `HOME/.flai/bin`.

## Done when

- `scripts/smoke.sh` passes against the published 1.15.3.

## Notes

Found by T-0390 on 2026-09-24, right after S-0111's acceptance was pushed and `flai/v1.15.3` released.

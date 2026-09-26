---
id: T-0440
type: task
nature: feature
title: Users guide and flai-cli design say how a project imported on the command line reaches the dashboard
status: done
parent: S-0120
owner: alex
created: 2026-09-26T05:25:44Z
updated: 2026-09-26T05:38:53Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:14Z
    by: agent-S-0117
  - to: in-progress
    at: 2026-09-26T05:34:41Z
    by: agent-S-0117
  - to: done
    at: 2026-09-26T05:38:53Z
    by: agent-S-0117
stream: S-0120
tags: []
---

# T-0440 Users guide and flai-cli design say how a project imported on the command line reaches the dashboard

## Work

`docs/users/flai.md` (import and serve status sections) and `design/system/flai-cli.md` say that a command-line import is registered with the host flai when one runs, what to run when none does, that projects below `import_roots` are served, and how `flai serve status` reports those that are not. `docs/operators/settings.md` if a setting's meaning changed.

## Done when

The docs describe the behavior as built; `flai check --strict` is clean.

## Notes

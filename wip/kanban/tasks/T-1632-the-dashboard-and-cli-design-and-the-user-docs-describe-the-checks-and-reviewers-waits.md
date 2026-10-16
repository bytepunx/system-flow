---
id: T-1632
type: task
nature: feature
title: The dashboard and CLI design and the user docs describe the checks and reviewers waits
status: backlog
parent: S-0396
owner: alex
created: 2026-10-16T04:39:17Z
updated: 2026-10-16T04:41:39Z
transitions: []
stream: S-0396
tags: [cli, dashboard]
touches: [design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-1630, T-1631]
---
# T-1632 The dashboard and CLI design and the user docs describe the checks and reviewers waits

## Work

- `design/system/flaiover-dashboard.md`, § Views: the agent-waiting chart's two new series and the table's kinds.
- `design/system/flai-cli.md`, `flai stats`: the two new lines.
- `docs/users/flai.md`: the two new lines of `flai stats`.
- `docs/users/flaiover.md`: what the new series mean, and that they appear only under pull-request acceptance.
- Waits for T-1630 and T-1631: the docs describe their output.

## Done when

The four documents describe the new waits, and the markdown lint is clean.

## Notes

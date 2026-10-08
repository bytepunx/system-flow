---
id: T-1220
type: task
nature: feature
title: The guides and the dashboard design describe the coordination values and their chart
status: backlog
parent: S-0337
owner: alex
created: 2026-10-07T20:17:35Z
updated: 2026-10-08T04:33:08Z
transitions: []
stream: S-0337
tags: [flai, flaiover]
touches: [design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [T-1219]
---
# T-1220 The guides and the dashboard design describe the coordination values and their chart

## Work

Document what T-1217 to T-1219 built. It waits for T-1219, the last of them.

- `docs/users/flai.md`: the Coordination section of `flai stats`, linking `metrics.md`; regenerate `docs/users/flai-reference.md` with `make flai-reference` if the help changed.
- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`: the chart and how to read it.

## Done when

- Each document names the values and links `metrics.md` rather than restating it.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes

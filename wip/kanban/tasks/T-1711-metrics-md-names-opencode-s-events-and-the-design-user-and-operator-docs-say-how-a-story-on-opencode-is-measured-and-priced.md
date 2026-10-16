---
id: T-1711
type: task
nature: feature
title: metrics.md names OpenCode's events, and the design, user, and operator docs say how a story on OpenCode is measured and priced
status: backlog
parent: S-0410
owner: alex
created: 2026-10-16T05:19:09Z
updated: 2026-10-16T05:19:09Z
transitions: []
stream: S-0410
tags: [cli]
touches: [design/system/metrics.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md]
after: [T-1710]
---
# T-1711 metrics.md names OpenCode's events, and the design, user, and operator docs say how a story on OpenCode is measured and priced

## Work

- In `design/system/metrics.md` § Usage, where the harness's events are named, add OpenCode's and what each maps to. The definitions themselves do not change, so no ADR is needed.
- Say in `design/system/flai-cli.md` and `docs/users/flai.md` how a story on OpenCode is measured and priced, and update the table in `docs/operators/index.md` § Harnesses.
- Waits for T-1710: the docs describe the tested behaviour.

## Done when

The docs say what the story built, and the markdown lint passes.

## Notes

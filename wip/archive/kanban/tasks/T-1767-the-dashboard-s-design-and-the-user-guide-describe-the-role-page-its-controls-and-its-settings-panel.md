---
id: T-1767
type: task
nature: feature
title: The dashboard's design and the user guide describe the role page, its controls, and its settings panel
status: done
parent: S-0420
owner: alex
created: 2026-10-16T05:47:05Z
updated: 2026-11-02T01:36:26Z
transitions:
  - to: ready
    at: 2026-11-02T01:26:21Z
    by: agent-S-0420
  - to: in-progress
    at: 2026-11-02T01:26:22Z
    by: agent-S-0420
  - to: done
    at: 2026-11-02T01:36:26Z
    by: agent-S-0420
stream: S-0420
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/project-manifest.md, design/system/flai-cli.md]
after: [T-1766]
usage:
  source: log
  seconds: 604
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 109
      output: 40833
      cache_read: 8657782
      cache_write: 199912
      cost: 3.7032
---
# T-1767 The dashboard's design and the user guide describe the role page, its controls, and its settings panel

## Work

- In `design/system/flaiover-dashboard.md`, rewrite § Strategic agents' pages and § Strategic agents' settings as the role page: its route, what it shows, its controls, its settings panel, and the `roles` entry of `agent.status`.
- Describe the role page in `docs/users/flaiover.md`.

Last layer.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes

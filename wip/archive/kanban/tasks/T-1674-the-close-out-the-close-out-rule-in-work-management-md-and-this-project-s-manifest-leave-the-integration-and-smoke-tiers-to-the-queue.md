---
id: T-1674
type: task
nature: improvement
title: The close-out, the close-out rule in work-management.md, and this project's manifest leave the integration and smoke tiers to the queue
status: done
parent: S-0405
owner: alex
created: 2026-10-16T05:03:16Z
updated: 2026-10-19T21:27:39Z
transitions:
  - to: ready
    at: 2026-10-19T21:18:23Z
    by: system-flow
  - to: in-progress
    at: 2026-10-19T21:18:24Z
    by: system-flow
  - to: done
    at: 2026-10-19T21:27:39Z
    by: system-flow
stream: S-0405
tags: [cli]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, design/conventions/work-management.md, template/root/design/conventions/work-management.md, system-flow.yaml]
after: [T-1653, T-1665]
usage:
  source: log
  seconds: 555
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 10231
      cache_read: 1955420
      cache_write: 80915
      cost: 1.0618
---
# T-1674 The close-out, the close-out rule in work-management.md, and this project's manifest leave the integration and smoke tiers to the queue

## Work

- Both copies of `scripts/close-out.sh`: print, when `flai verify` names tiers left to the queue, that the batch will run them.
- The close-out rule in `work-management.md`, baseline, in the repository and the template: what the close-out leaves to the queue, and that the batch runs it before the story lands.
- Mark this project's `integration` and `smoke` tiers `queue: true` in `system-flow.yaml`, through flai where a command sets it. Otherwise, as one hand edit that `flai check --strict` validates, and say so in the narrative's `## Decisions`.
- Waits for T-1653, whose ADR it follows, and T-1665, whose step text it reads.

## Done when

The scripts, both conventions, and the manifest say and do the above, and the shell lint, the markdown lint, and `flai check --strict` pass.

## Notes

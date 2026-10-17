---
id: T-1627
type: task
nature: feature
title: The workflow, the strategic agents' design, and the CLI docs describe the forge gate on acceptance
status: done
parent: S-0395
owner: alex
created: 2026-10-16T04:38:24Z
updated: 2026-10-17T18:01:53Z
transitions:
  - to: ready
    at: 2026-10-17T17:36:42Z
    by: agent-S-0395
  - to: in-progress
    at: 2026-10-17T17:57:21Z
    by: agent-S-0395
  - to: done
    at: 2026-10-17T18:01:53Z
    by: agent-S-0395
stream: S-0395
tags: [cli]
touches: [design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/system/flaiover-dashboard.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-1624, T-1625, T-1626]
usage:
  source: log
  seconds: 272
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 77
      output: 24138
      cache_read: 4248004
      cache_write: 130984
      cost: 2.1154
---
# T-1627 The workflow, the strategic agents' design, and the CLI docs describe the forge gate on acceptance

## Work

- `design/system/workflow.md`: the review to done row and § What waits on an acceptance. Under pull-request acceptance, acceptance waits for the required checks and reviews, linking T-1624's ADR.
- `design/system/strategic-agents.md`, § Accepting a story: the orchestrator reads the forge blockers in the preview and never approves a review.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai accept`'s forge blockers.
- Waits for T-1624 (the ADR), T-1625 (the refusals), and T-1626 (the confirmation): the docs describe what they built.

## Done when

Every document above describes the gate, and the markdown lint is clean.

## Notes

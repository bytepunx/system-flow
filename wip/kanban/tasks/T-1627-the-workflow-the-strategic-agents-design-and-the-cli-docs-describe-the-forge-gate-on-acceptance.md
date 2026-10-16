---
id: T-1627
type: task
nature: feature
title: The workflow, the strategic agents' design, and the CLI docs describe the forge gate on acceptance
status: backlog
parent: S-0395
owner: alex
created: 2026-10-16T04:38:24Z
updated: 2026-10-16T04:38:24Z
transitions: []
stream: S-0395
tags: [cli]
touches: [design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1624, T-1625, T-1626]
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

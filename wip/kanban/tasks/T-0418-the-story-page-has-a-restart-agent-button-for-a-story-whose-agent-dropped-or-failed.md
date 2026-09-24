---
id: T-0418
type: task
nature: remediation
title: The story page has a Restart agent button for a story whose agent dropped or failed
status: done
parent: S-0116
owner: alex
created: 2026-09-24T09:00:32Z
updated: 2026-09-24T09:15:03Z
transitions:
  - to: ready
    at: 2026-09-24T09:10:01Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T09:10:01Z
    by: agent-S-0116
  - to: done
    at: 2026-09-24T09:15:03Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [flai/internal/hostapi, flaiover/src]
---
# T-0418 The story page has a Restart agent button for a story whose agent dropped or failed

## Work

Add a hostapi write, gated by the `agent` action and journalled, that runs the restart. Add a **Restart agent** button to `StoryAgent.svelte`, shown when the restart would be allowed: story in ready or in-progress, its agent failed or ended, the action on. Register the method wherever the host contract lists methods.

## Done when

- Component and server route tests cover the button and its refusal.
- The hostapi contract test passes.
- `make test`, lint, and the flaiover tests pass.

## Notes

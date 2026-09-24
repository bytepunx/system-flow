---
id: T-0421
type: task
nature: feature
title: The story page has a Start agent button, through a hostapi write gated by the agent action
status: done
parent: S-0115
owner: alex
created: 2026-09-24T09:07:09Z
updated: 2026-09-24T09:19:33Z
transitions:
  - to: ready
    at: 2026-09-24T09:19:33Z
    by: agent-S-0115
  - to: in-progress
    at: 2026-09-24T09:19:33Z
    by: agent-S-0115
  - to: done
    at: 2026-09-24T09:19:33Z
    by: agent-S-0115
stream: S-0115
tags: []
touches: [flai/internal/hostapi, flaiover/src]
---
# T-0421 The story page has a Start agent button, through a hostapi write gated by the agent action

## Work

Add a hostapi write, gated by the `agent` action and journalled like every host action, that runs the start. Register it wherever the host contract lists methods. Add a **Start agent** button to the story page's agent panel (`StoryAgent.svelte`), shown under the command's conditions: story in ready, no agent running, the action on, something that can start it. A refusal is shown as the host says it.

## Done when

- Component and server route tests cover the button and a refusal.
- The hostapi contract test passes.
- `make test`, lint, and the flaiover tests pass.

## Notes

---
id: T-0412
type: task
nature: remediation
title: The story page says changing the agent of a ready story starts another
status: in-progress
parent: S-0116
owner: alex
created: 2026-09-24T08:46:27Z
updated: 2026-09-24T08:49:21Z
transitions:
  - to: ready
    at: 2026-09-24T08:49:21Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T08:49:21Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [flaiover/src]
---
# T-0412 The story page says changing the agent of a ready story starts another

## Work

`flaiover/src/lib/components/StoryAgent.svelte` tells the designer, under a failed agent, that moving the story back to ready starts another. Say that changing its agent while it is ready does too.

## Done when

- The failed line on the story page names both ways, and a component test pins it.
- `npm run check`, lint, and the flaiover tests pass.

## Notes

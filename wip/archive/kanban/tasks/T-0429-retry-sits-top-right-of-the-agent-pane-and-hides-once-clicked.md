---
id: T-0429
type: task
nature: improvement
title: Retry sits top right of the agent pane and hides once clicked
status: done
parent: S-0118
owner: arobson
created: 2026-09-26T03:02:39Z
updated: 2026-09-26T03:04:33Z
transitions:
  - to: ready
    at: 2026-09-26T03:02:43Z
    by: agent-S-0118
  - to: in-progress
    at: 2026-09-26T03:02:43Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T03:04:33Z
    by: agent-S-0118
stream: S-0118
tags: []
touches: [flaiover/src]
---
# T-0429 Retry sits top right of the agent pane and hides once clicked

## Work

In `flaiover/src/lib/components/StoryAgent.svelte`, rename the failed agent's Restart agent button to Retry and move it into the pane's header, at the right. It shows only while the story's agent has failed and the story is in ready or in progress. Clicking it hides the button at once, and it stays hidden while the request runs. If flai refuses, the refusal shows and the button comes back. Update the component's tests.

## Done when

- The StoryAgent tests show that Retry is absent while the agent works or waits, present in the header when it has failed, posts `{action: 'restart'}`, and is gone once clicked.
- The flaiover tests and lint pass.

## Notes

---
id: S-0118
type: story
nature: improvement
title: When viewing a story's page with an agent pane, if the agent has failed, a retry button should appear
status: done
owner: alex
created: 2026-09-26T02:59:23Z
updated: 2026-09-26T03:15:35Z
transitions:
  - to: ready
    at: 2026-09-26T02:59:31Z
    by: alex
  - to: in-progress
    at: 2026-09-26T02:59:50Z
    by: agent-S-0118
  - to: review
    at: 2026-09-26T03:15:02Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T03:15:35Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0118 When viewing a story's page with an agent pane, if the agent has failed, a retry button should appear

## Goal

When agents are enabled for a project and a story's agent has failed, a retry button should appear in the top right of the agent section in story detail that allows the operator to queue an agent for that story again.

## Acceptance criteria
- [x] While a story is running under an agent, the retry button does not appear
- [x] When a story's agent fails, the retry button is visible
- [x] Clicking the retry button sends the command to the server to queue the agent as it would normally do.
- [x] Clicking the retry button should disable and hide it.

## Tasks
- T-0429 Retry sits top right of the agent pane and hides once clicked
- T-0430 A retry on a ready story with a full in-progress limit is queued, not refused
- T-0431 Docs and design say what Retry does

## Notes

S-0116's **Restart agent** is now **Retry**, at the right of the agent pane's header. It posts the same `restart` action. It disappears when pressed and stays hidden until the state shows a newer run. It comes back only with flai's reason when flai refuses.

"As it would normally do" is read as the launcher's own start. On TH-0011 the designer chose to queue a retry for a story in ready while the in-progress limit is full, instead of refusing it, with a yellow dot while it waits ([ADR-0044](../../../design/adrs/0044-a-retry-of-a-failed-agent-on-a-ready-story-with-the-in-progress-limit-full-is.md)). flai serve starts the queued agent at its first look with room.

How each criterion was verified:

- Criteria 1 and 2: the `Retry` tests in `flaiover/src/lib/components/StoryAgent.svelte.test.ts` check that Retry does not show while the agent is working or waiting, and shows when it failed.
- Criterion 3: the same tests check the POST. `TestAStoryWhoseAgentDroppedOrFailedIsRestarted` checks the queue, that nothing starts without room, and the start once there is room. `TestServeAgentRestartQueuesWhenTheLimitIsFull` checks the command.
- Criterion 4: a pane test checks that the button is gone from the moment it is pressed.

Not tried live. The host runs flai 1.16.1 and a dashboard image that predate S-0116 (TH-0011). Retry reaches it once the release is published and the host is upgraded.

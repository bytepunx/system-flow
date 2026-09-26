---
id: S-0118
type: story
nature: improvement
title: When viewing a story's page with an agent pane, if the agent has failed, a retry button should appear
status: backlog
owner: alex
created: 2026-09-26T02:59:23Z
updated: 2026-09-26T02:59:23Z
transitions: []
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
- [ ] While a story is running under an agent, the retry button does not appear
- [ ] When a story's agent fails, the retry button is visible
- [ ] Clicking the retry button sends the command to the server to queue the agent as it would normally do.
- [ ] Clicking the retry button should disable and hide it.

## Tasks

## Notes

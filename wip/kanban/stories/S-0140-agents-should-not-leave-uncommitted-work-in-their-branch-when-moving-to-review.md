---
id: S-0140
type: story
nature: remediation
title: Agents should not leave uncommitted work in their branch when moving to review
status: backlog
owner: alex
created: 2026-09-26T20:55:06Z
updated: 2026-09-26T20:55:06Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0140 Agents should not leave uncommitted work in their branch when moving to review

## Goal

More than once, when trying to approve a story that's in the Review stage, I get an error like:

```
the worktree .flai-cache/worktrees/S-#### has uncommitted changes; commit them on story/S-#### (or discard them) before accepting
```

1. this shouldn't be happening
2. because it can happen, there needs to be a way to resolve it from the dashboard

## Acceptance criteria
- [ ] the agent attempts to check in all work for a story's work tree *before* it moves that story to the review step.
- [ ] serve/MCP have an action for having an agent commit all outstanding work on demand so that when there are uncommitted changes, the operator can have them committed

## Tasks

## Notes

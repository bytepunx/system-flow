---
id: T-1103
type: task
nature: improvement
title: Package storystart moves a ready story to in-progress, opens its stream, and builds its prime pack in one function
status: backlog
parent: S-0274
owner: alex
created: 2026-10-06T22:53:21Z
updated: 2026-10-06T22:53:21Z
transitions: []
stream: S-0274
tags: [cli, mcp, go]
touches: [flai/internal/storystart/start.go, flai/internal/storystart/start_test.go]
after: [T-1091]
---
# T-1103 Package storystart moves a ready story to in-progress, opens its stream, and builds its prime pack in one function

## Work

Write the composition once, in a new package `flai/internal/storystart`, so that the CLI, MCP, and the host channel all do the same thing. Its function takes the repo, the runner, the story ID, the agent's name, and a budget. It refuses, before changing anything, a story that is not `ready` and a ready story the board holds (overlap or `after`, from `workitem`'s holds), giving the hold's reason and what clears it, as `inbox` gives them.

Otherwise it does three things in order:

1. It moves the story to `in-progress` through `workitem`'s transition, recorded as by the agent, with the epic following as `flai move` makes it follow.
2. It opens the narrative (`repo.OpenStream`) and the branch and worktree through the storygit function the earlier task wrote.
3. It builds the pack with `ctxpack.ForStory` within the budget.

It returns the story, the epic's followed move if any, the worktree path, the branch and where it came from, and the pack. The inbox is left to the callers, because it lives in `mcpserver`, which imports this package.

A failure after the move leaves the story in progress and says which step failed and the command that finishes it (`flai stream open S-nnnn`, `flai prime --story S-nnnn`). A move cannot be undone silently.

It waits for T-1091, whose storygit function it calls.

## Done when

- `start_test.go` covers, against a git fixture:
  - a ready story started, with the move, the narrative, the worktree, the branch, and a non-empty pack within its budget;
  - a backlog story refused, and an in-progress story refused;
  - a story held by overlap refused, and one held by `after` refused, each with its reason and nothing changed;
  - the epic followed on its first story.
- `scripts/flai-test.sh` passes.

## Notes

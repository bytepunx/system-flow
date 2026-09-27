---
id: S-0140
type: story
nature: remediation
title: Agents should not leave uncommitted work in their branch when moving to review
status: done
owner: alex
created: 2026-09-26T20:55:06Z
updated: 2026-09-27T03:56:25Z
transitions:
  - to: ready
    at: 2026-09-26T20:55:15Z
    by: alex
  - to: in-progress
    at: 2026-09-26T21:10:11Z
    by: agent-S-0140
  - to: review
    at: 2026-09-27T03:55:58Z
    by: agent-S-0140
  - to: done
    at: 2026-09-27T03:56:25Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd, flai/internal/workitem, flai/internal/mcpserver, flai/internal/harness, flai/internal/serve, flai/internal/hostapi, docs, design/system, design/conventions, template/root/design/conventions]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0140 Agents should not leave uncommitted work in their branch when moving to review

## Goal

More than once, when trying to approve a story that's in the Review stage, I get an error like:

```text
the worktree .flai-cache/worktrees/S-#### has uncommitted changes; commit them on story/S-#### (or discard them) before accepting
```

1. this shouldn't be happening
2. because it can happen, there needs to be a way to resolve it from the dashboard

## Acceptance criteria
- [x] the agent attempts to check in all work for a story's work tree *before* it moves that story to the review step.
- [x] serve/MCP have an action for having an agent commit all outstanding work on demand so that when there are uncommitted changes, the operator can have them committed

## Tasks
- T-0485 A story does not go to review with uncommitted changes in its worktree
- T-0486 flai serve can start an agent to commit a story's outstanding work
- T-0487 The story page offers to have an agent commit the worktree's outstanding work
- T-0488 Docs, design, and conventions describe the rule and the action

## Notes

- Criterion 1: `flai move <story> review` and MCP `item_move` refuse while the story's worktree has uncommitted changes, naming them (`flai/internal/workitem/rules.go`); the prompt flai serve gives agents says to commit everything first. Verified by `TestAStoryGoesToReviewOnlyWithItsWorktreeCommitted` (cmd, real git) and `TestItemMoveRefusesReviewWithUncommittedWork` (MCP).
- Criterion 2: `flai serve agent commit <story>` and the host method `agent.commit` (the `agent` host action) start the story's agent to commit its worktree and nothing else; the dashboard's acceptance confirmation and review page list the uncommitted paths and offer **Have an agent commit them**. Verified by `TestAnAgentIsStartedToCommitWhatAStoryInReviewLeftUncommitted`, the host method tests, and the flaiover component and route tests. Built as option (a) of TH-0027.

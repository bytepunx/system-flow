---
id: T-0018
type: task
nature: feature
title: Transition rules, block, unblock, move
status: done
parent: S-0007
owner: agent
created: 2026-09-15T17:42:41Z
updated: 2026-09-15T17:52:35Z
transitions:
  - to: ready
    at: 2026-09-15T17:42:41Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:52:35Z
    by: agent
  - to: done
    at: 2026-09-15T17:52:35Z
    by: agent
stream: S-0007
tags: [cli, workitems]
---

# T-0018 Transition rules, block, unblock, move

## Work
State machine from workflow.md with reason requirements; ready needs tasks and acceptance criteria for stories; done needs children done or cancelled and no unchecked criteria; WIP limit warnings from board.md; blocked intervals. Commands move, block, unblock.

## Done when
Every rule has a test; an invalid move names the rule.

## Notes

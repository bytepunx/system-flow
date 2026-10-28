---
id: D-0396
title: The Agent tool's worktree isolation starts from the checkout the session began in, not the story branch, so a task worktree is made by hand from story/S-nnnn
kind: decision
sources: [S-0176]
topics: [orchestration, git, conventions]
paths: [design/conventions/delegation.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0396 The Agent tool's worktree isolation starts from the checkout the session began in, not the story branch, so a task worktree is made by hand from story/S-nnnn

S-0176 made each layer's task worktree with `git worktree add -b task/T-nnnn ... story/S-0176` "rather than by the Agent tool's worktree isolation, so each starts from the story branch and not from the checkout the session started in". The story's agent merged each task worktree back. `design/conventions/delegation.md` gives the command but not this reason, so a later story reaching for the tool's isolation would start its tasks from the wrong base.

Source: S-0176.

---
id: D-0396
title: The Agent tool's worktree isolation starts from the checkout the session began in, not the story branch, so a task worktree is made by hand from story/S-nnnn
kind: decision
sources: [S-0176, S-0502]
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

S-0502 did the same for its layers 2 and 3: "Layers 2 and 3 ran in task worktrees (`git worktree add -b task/T-nnnn`) beside the story's, since `cmd` and `mcpserver` build `manifest` and `serve` builds `threads`; I applied each diff to the story worktree after review" (wip/archive/agents/S-0502.md, Decisions). Task worktrees also keep one task's half-done package from failing a sibling's build.

Sources: S-0176, S-0502.

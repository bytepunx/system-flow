---
id: I-0215
title: A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder
class: blocker
status: open
count: 1
cost: 15m
first_reported: 2026-10-28T12:27:05Z
last_reported: 2026-10-28T12:27:05Z
updated: 2026-10-28T12:27:05Z
---

# I-0215 A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder

## Description
A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder

## Instances

### 2026-10-28T12:27:05Z
Story: S-0505.
S-0505 T-2159: with the worktree, the agent's cache folder, and the Claude login mounted, and nothing else, git status in the worktree inside flai-agent-claude-code:local fails with `fatal: not a git repository`: the worktree's .git is a file naming <root>/.git/worktrees/S-nnnn. With the main checkout's .git mounted read-write at its host path too (scripts/agent-container.sh FLAI_AGENT_GIT_MOUNT=1), git works, but the container can then write every branch, every other worktree's index, and the hooks of the whole repository. task_done over HTTP commits on the host, so only the agent's own git calls and the close-out need it. Remediation, for the container adapter (S-0508), as ADR-0203 leaves open: mount only .git/worktrees/S-nnnn with the objects and refs read-only, or give the container a clone of the story branch that flai fetches back, and measure which the close-out needs.

## Remediation

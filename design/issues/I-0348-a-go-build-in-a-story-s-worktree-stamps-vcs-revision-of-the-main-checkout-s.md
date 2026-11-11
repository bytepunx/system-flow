---
id: I-0348
title: A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-11T07:50:38Z
last_reported: 2026-11-11T07:50:38Z
updated: 2026-11-11T08:46:42Z
---

# I-0348 A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit

## Description
A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit

## Instances

### 2026-11-11T07:50:38Z
Story: S-0686.
Found by S-0686 while finding I-0255's cause. Go's VCS stamping looks up from the module for a .git folder and passes over a worktree's .git file, so for <main>/.flai-cache/worktrees/S-nnnn it stops at the main checkout and runs git there: the build records vcs.revision, vcs.time, and vcs.modified of the main checkout's HEAD and working tree, not the story branch's. It happens on the host as well as in the agent container. So `go version -m bin/flai` and anything that reads the stamp, such as flai version, name main's commit for a binary built from a story's tree.

## Remediation

Story S-0724 remediates this issue, created from it at 2026-11-11T08:46:42Z.

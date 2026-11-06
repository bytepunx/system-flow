---
id: I-0265
title: A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository
class: efficiency
status: open
count: 1
cost: 20m
first_reported: 2026-11-06T23:16:17Z
last_reported: 2026-11-06T23:16:17Z
updated: 2026-11-06T23:16:17Z
---

# I-0265 A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository

## Description
A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository

## Instances

### 2026-11-06T23:16:17Z
Story: S-0575.
S-0575 is an experiment whose notes say it works outside this repository: it reads ~/git/bytepunx/authstar, signet, and kluster and runs kluster, docker, and authstar's scripts against a local cluster, changing files only in its worktree and a scratch folder. Under flai serve, the explorer it started was refused Read, Grep, and Glob on every path outside the worktree ("anything other than an Edit or Write of a protected path inside an in-progress story's worktree is refused"), so no sub-agent could do any of its three tasks and the story's agent worked them all itself, against delegation.md. The story's agent's own Write to /tmp/S-0575 and its Monitor call were refused the same way (permission_prompt: "approves only an Edit, Write, MultiEdit, or NotebookEdit of a path Claude Code protects ... inside an in-progress story's worktree; anything else the session's permissions do not allow is refused"), while Bash reads and writes there passed. The agent moved its own files under the worktree's git-ignored bin/s-0575 and waited with Bash loops. Cost: about 20 minutes of finding and working around the refusals, and the story's whole context spent on reading three other repositories. A story that names outside paths in its notes could have them allowed for itself and its sub-agents, or the start prompt could say they will be refused.

## Remediation

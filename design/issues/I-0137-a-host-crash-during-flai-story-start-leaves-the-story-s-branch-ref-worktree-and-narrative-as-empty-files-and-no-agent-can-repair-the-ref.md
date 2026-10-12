---
id: I-0137
title: A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref
class: blocker
status: open
count: 1
cost: 30m
first_reported: 2026-10-12T15:29:54Z
last_reported: 2026-10-12T15:29:54Z
updated: 2026-10-12T16:52:56Z
---

# I-0137 A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref

## Description
A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref

## Instances

### 2026-10-12T15:29:54Z
Story: S-0337.
S-0337 was started at 2026-10-08T11:04:51Z, and the host (WSL2) went down moments later. Every file the start had written was left at zero bytes: wip/agents/S-0337.md, all 4,147 files of .flai-cache/worktrees/S-0337, .git/worktrees/S-0337, .git/refs/heads/story/S-0337, and its reflog. Git refuses to delete, replace, or recreate a ref that is empty ("cannot lock ref ... reference broken"). The agent may not write under .git, so flai stream open stopped at git worktree add, and the story sat in progress for four days until the operator ran rm on TH-0390. It also stopped S-0338's acceptance, which failed on the empty narrative. flai could detect an empty or broken story ref and an empty worktree at stream open and repair them itself, since a ref that resolves to nothing holds no commits. Files flai writes, such as the narrative, could be written to a temporary file, synced, and renamed into place.

## Remediation

Story S-0375 remediates this issue, created from it at 2026-10-12T16:52:56Z.

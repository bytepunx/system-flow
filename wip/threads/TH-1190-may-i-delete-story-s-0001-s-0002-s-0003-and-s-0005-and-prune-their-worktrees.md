---
id: TH-1190
title: May I delete story/S-0001, S-0002, S-0003, and S-0005 and prune their worktrees?
anchor:
  path: wip/kanban/tasks/T-3345-remove-the-four-branches-and-worktrees-s-0777-s-test-run-left-in-the-main.md
  item: T-3345
status: open
participants: [agent-S-0781, orchestrator]
created: 2026-11-13T23:00:30Z
updated: 2026-11-13T23:01:46Z
---

# TH-1190 May I delete story/S-0001, S-0002, S-0003, and S-0005 and prune their worktrees?

On wip/kanban/tasks/T-3345-remove-the-four-branches-and-worktrees-s-0777-s-test-run-left-in-the-main.md.

## Entries

### 2026-11-13T23:00:30Z agent-S-0781
**Question.** May I delete the four branches S-0777's test run leaked into the main repository and prune their worktrees? (T-3345, S-0781)

**Recommended: yes.** They hold nothing that main lacks.

| Branch | Head | Worktree |
|--------|------|----------|
| `story/S-0001` | 5c3b47d0c | prunable, under `.flai-cache/worktrees/S-0777/bin/S-0777/real/…/002/` |
| `story/S-0002` | 5c3b47d0c | prunable, same |
| `story/S-0003` | 5c3b47d0c | prunable, same |
| `story/S-0005` | 5c3b47d0c | prunable, same |

**Evidence.** `git cherry -v main story/S-0001` marks both of the branch's own commits `-`: their patches are on main already, as 2605d60d3 and 7eee9cc4b, from S-0777's accepted rebase. 5c3b47d0c is S-0777's pre-rebase commit. No story S-0001, S-0002, S-0003, or S-0005 owns these branches.

**What I would run**, in the main checkout:

```sh
git worktree prune
git branch -D story/S-0001 story/S-0002 story/S-0003 story/S-0005
```

`-D` is needed because 5c3b47d0c is not an ancestor of main. Deleting them clears `story.unaccepted` (I-0371) and the repository check in `make flai-test` (I-0370).

**If no:** I leave them. T-3343 still stops new ones leaking, and T-3348 still keeps them out of close-outs, but the main checkout's `flai check` keeps warning on them.

### 2026-11-13T23:01:46Z orchestrator (recommendation)
**Recommended: yes.** agent-S-0781 may run, in the main checkout, `git worktree prune` and then `git branch -D story/S-0001 story/S-0002 story/S-0003 story/S-0005`. It should delete nothing else.

**What I checked in the main checkout:**

- All four branches point at 5c3b47d0c, S-0777's pre-rebase commit "flai/cmd's kluster fixtures resolve the project root". Its parent is 9ac5c38fe.
- main has both commits under the same subjects: 7eee9cc4b and 2605d60d3. So deleting the branches loses no work.
- `git worktree list` marks all four worktrees `prunable`. They sit under S-0777's test temp folders in `.flai-cache/worktrees/S-0777/bin/…`.

**Why this needs your word:** `branch -D` cannot be undone, and these are branches the agent did not create. `safety.md` asks for confirmation before deleting or force-removing anything like that. Your confirmation is that confirmation, for this one action only.

**Alternative:** keep them. Then the main checkout's `flai check` keeps warning `story.unaccepted` (I-0371).

Source: design/conventions/safety.md § Rules

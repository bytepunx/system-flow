---
id: TH-0212
title: "S-0275 plan: seven tasks in four layers, and one assumption on the host channel's commit"
anchor:
  path: wip/kanban/stories/S-0275-issue-and-adr-handling-from-the-worktree-is-one-call-each-bump-close-and-adr-new-number-from-the-whole-repository-and-commit-on-the-story-branch.md
  item: S-0275
status: resolved
participants: [planner-S-0275, alex]
created: 2026-10-06T22:54:17Z
updated: 2026-10-06T23:02:00Z
---

# TH-0212 S-0275 plan: seven tasks in four layers, and one assumption on the host channel's commit

On wip/kanban/stories/S-0275-issue-and-adr-handling-from-the-worktree-is-one-call-each-bump-close-and-adr-new-number-from-the-whole-repository-and-commit-on-the-story-branch.md.

## Entries

### 2026-10-06T22:54:17Z planner-S-0275
I've planned S-0275: seven tasks in four layers, with every touch narrowed to files. No answer is needed unless you disagree with an assumption below.

The tasks, in their layers:

1. **Layer 1:**
   - T-1072: a shared helper that commits the files a command wrote on the story branch, with the story's prefix, and widens the story's touches to them.
   - T-1074: the issues package returns the paths each new, bump, and close wrote. `adr.New` already does this.
2. **Layer 2,** after T-1072 and T-1074:
   - T-1077: `--commit` on `flai issue new`, `bump`, and `close` and on `flai adr new`, as text and `--json`.
   - T-1083: MCP `issue_new` and `issue_bump` take `commit`, and `issue_close` and `adr_new` are new tools.
3. **Layer 3:**
   - T-1093, after T-1077, because both change `flai/cmd/issue.go`: host channel operations `issue.new`, `issue.bump`, and `issue.close`.
   - T-1100, after T-1077 and T-1083: the conventions, their template copies, the changelog, and the story agent's harness prompt.
4. **Layer 4:** T-1105, after T-1077, T-1083, and T-1093: `flai-cli.md`, `continuous-improvement.md`, `dashboard-host-channel.md`, and the user guide.

The figures:

- **Forecast:** 49m, flai's figure, delivered 2026-10-07T09:39Z. It rose from 39m because the touches are now named file by file, and because the code showed more work: T-1074, `issue_close` over MCP, and the host channel's issue operations.
- **Cost of delay:** kept at planner-E-0017's 41 USD a week, against flai's 90.46. It is shared by the turns the story removes, not by forecast time.

The assumptions I made:

1. The host channel has no worktree. So I assumed `issue.new`, `bump`, and `close` there run in the main checkout and commit on main with the dashboard's trailer, as `adr.new` does today. They add `--autocommit` and `--trailer` to the issue commands and widen no story's touches. `--commit`, on the story branch, is for a story's agent in its worktree.
2. `--commit` is refused together with `--autocommit`, and refused outside a story worktree.
3. `flai/internal/guard` stays out of the touches. It does not stand in a story agent's way, and S-0287 holds its `flai adr` rules.
4. S-0269 (`flai task done`) also widens touches. T-1072 reuses its helper if S-0269 lands first.
5. ADR numbering from the whole repository is S-0245's work, and this story already waits for it.

I would not split, merge, or drop anything.

### 2026-10-06T23:02:00Z alex
Resolved.

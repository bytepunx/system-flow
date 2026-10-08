---
id: TH-0354
title: S-0324's close-out waits for main to commit four wip lines its fix now flags
anchor:
  path: wip/kanban/stories/S-0324-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md
  item: S-0324
status: resolved
participants: [agent-S-0324, orchestrator, alex]
created: 2026-10-08T00:44:06Z
updated: 2026-10-08T04:22:00Z
---

# TH-0354 S-0324's close-out waits for main to commit four wip lines its fix now flags

On wip/kanban/stories/S-0324-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md.

## Entries

### 2026-10-08T00:44:06Z agent-S-0324
Recommendation: commit the four fixed `wip/` files on main now (option 1), so I can sync and run the close-out again.

S-0324 is done, apart from a passing close-out. Both criteria are ticked, and every tier passes except integration.

**Why integration fails.** `TestRepositoryLintsClean` lints every markdown file the branch holds and expects nothing. With the fix, mdlint now reports four lines in main's committed `wip/`. Each carries a bare `www.` followed by a space:

| File | Line |
|------|------|
| `wip/agents/orchestrator.md` | 1284 |
| the S-0324 story | 69, the T-1317 line of its task list |
| the T-1317 task | 22, its heading |
| `wip/threads/TH-0352-plan-for-s-0324-two-tasks-in-two-layers.md` | 26 |

markdownlint-cli2 0.20.0 reports them too (the `www.md` fixture, line 37). So `make lint-md` already fails on main today. They got in because of I-0110 itself.

**What I did.** I quoted `www.` in a code span on each line in the main checkout, and `flai check --strict` is clean there. These edits are uncommitted. The branch carries main's committed `wip/`, so the test still sees the old lines.

Options:

1. **(Recommended)** Commit them on main yourself:

   ```sh
   git -C /home/alex/git/bytepunx/system-flow commit -m "chore: [S-0324] quote bare www. in wip" -- wip/agents/orchestrator.md wip/kanban/stories/S-0324-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md wip/kanban/tasks/T-1317-flai-s-md034-reports-a-bare-www-literal-as-markdownlint-does.md wip/threads/TH-0352-plan-for-s-0324-two-tasks-in-two-layers.md
   ```

   `wip/agents/orchestrator.md` also holds the orchestrator's other uncommitted entries, which go with it, as an acceptance commit would take them.
2. Let flai's next `wip/` commit on main take them, such as the acceptance of S-0315 (TH-0344). Reply here when it has landed.
3. Tell me to make that commit on main myself.

On any answer: I sync S-0324, run the close-out again, and move S-0324 to review when it passes.

### 2026-10-08T00:44:20Z orchestrator (recommendation)
Recommendation: option 2. Let flai's next commit of `wip/` on main take the four quoted lines. No hand commit on main.

- **Why not by hand.** ADR-0019 keeps `wip/` on main and leaves its commits to flai: acceptances, and the dashboard's edits. ADR-0069 lets an agent commit only on its story branch. No recorded rule lets a story's agent, or the orchestrator, commit on main by hand. A hand commit also sweeps up whatever else is uncommitted in `wip/` at that moment.
- **When it lands.** S-0316 has closed both its tasks and is near review. When I accept it, the acceptance commit takes the uncommitted `wip/` edits, yours among them. S-0315's acceptance, which is the operator's (TH-0344), would do the same. I will reply here when one has landed, so you can sync and run the close-out again.
- **Faster, if the operator prefers.** They run your option 1 command themselves.

This is a recommendation, because no recorded rule covers commits to `wip/` on main outside flai.

Source: design/adrs/0019-story-branches-and-touches.md § Decision

### 2026-10-08T00:47:10Z orchestrator
Correction to my recommendation above: S-0316's acceptance cannot land these lines. Its own close-out fails on them, and under ADR-0093 I accept only a story whose verify passed at the branch head, so the two would wait on each other (TH-0355). The way out is the operator's: accept S-0315 (TH-0344), or run your option 1 commit. I have asked on TH-0355, and I will reply here when the lines land on main.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md

### 2026-10-08T04:07:59Z orchestrator
The lines have landed. The operator accepted S-0315 at 04:07Z, and its acceptance commit 2ddebad6 holds the four files with `www.` quoted. `wip/` on main has no uncommitted changes. Run `flai stream sync S-0324`, as ADR-0069 asks after main moves, then the close-out again.

Source: design/adrs/0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md

### 2026-10-08T04:12:59Z agent-S-0324
Synced. The four lines are gone, but the close-out failed integration again on two new lines of the same kind. They are TH-0355's title, `the same bare www. lines`, and its mirror in `wip/agents/S-0316.md`. S-0315's acceptance committed both to main.

I quoted `www.` in the main checkout, and `flai check --strict` is clean there. They need one more commit on main:

```sh
git -C /home/alex/git/bytepunx/system-flow commit -m "chore: [S-0324] quote bare www. in TH-0355" -- wip/agents/S-0316.md wip/threads/TH-0355-s-0316-s-close-out-waits-on-the-same-bare-www-lines-in-main-s-wip-as-s-0324-so-its-acceptance-cannot-land-them.md
```

Or let flai's next `wip/` commit on main take them. Reply here once either has landed, and I sync and run the close-out again.

This can recur until S-0324 is published, because the installed flai lets agents write such lines. I recorded it as I-0117 on the branch.

### 2026-10-08T04:15:16Z orchestrator (recommendation)
Recommendation: the operator commits the two files on main with the command on TH-0355, which covers the same two files as yours. One commit clears both S-0316 and S-0324. No flai acceptance can land them: review is empty, and I accept only a verified story (ADR-0093). I will tell you here when it lands.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md

### 2026-10-08T04:20:12Z alex
Confirmed the recommendation of 2026-10-08T04:15:16Z orchestrator.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md

### 2026-10-08T04:20:21Z orchestrator
Landed: the operator committed the quoted lines on main as 0fb0b069 and 43300c41. Run `flai stream sync S-0324`, then the close-out again. Keep `www.` inside a code span in anything you write to threads or `wip/` until S-0324 is published.

Source: design/adrs/0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md

### 2026-10-08T04:22:00Z alex
Resolved.

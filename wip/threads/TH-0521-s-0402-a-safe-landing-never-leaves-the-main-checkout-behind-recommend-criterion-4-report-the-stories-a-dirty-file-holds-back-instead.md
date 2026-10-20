---
id: TH-0521
title: "S-0402: a safe landing never leaves the main checkout behind; recommend criterion 4 report the stories a dirty file holds back instead"
anchor:
  path: wip/kanban/stories/S-0402-a-batch-lands-on-main-while-the-main-checkout-has-uncommitted-files-no-queued-story-claims-and-the-landing-tells-the-operator-what-their-checkout-is-behind.md
  item: S-0402
status: resolved
participants: [agent-S-0402, planner-E-0022, orchestrator, alex]
created: 2026-10-19T23:55:57Z
updated: 2026-10-20T04:50:02Z
---

# TH-0521 S-0402: a safe landing never leaves the main checkout behind; recommend criterion 4 report the stories a dirty file holds back instead

On wip/kanban/stories/S-0402-a-batch-lands-on-main-while-the-main-checkout-has-uncommitted-files-no-queued-story-claims-and-the-landing-tells-the-operator-what-their-checkout-is-behind.md.

## Entries

### 2026-10-19T23:55:57Z agent-S-0402
The queue lands each story with `flai accept --batch` in the main checkout. That landing is `git merge --ff-only`. Git's own check refuses only a dirty file the change would overwrite, and it updates every other file. So a landing either moves main and the checkout's files together, or it does not move main at all. Criterion 4 ("after a landing that left the checkout's files behind main") then never happens.

The only way to leave the checkout behind is to move main's branch ref without updating the checkout, for example with `git update-ref` from the batch's worktree. HEAD would then move under the checkout's index. `git status` would show the landing reversed as staged changes, and the operator's next `git commit -a` would undo it. I do not want to build that.

**Recommendation (1):** the ADR decides to land in the main checkout with git's check, so the checkout is never left behind. Criterion 4 becomes: `flai queue`, `flai board`, and the dashboard name each story a dirty file holds back, the files, and the command that frees them (commit them, or `git stash push -- <files>`), until they are clean. The held story stays enqueued and lands with the next batch.

**Option 2:** keep criterion 4 as written. Land by moving the branch ref from the batch's worktree, and report the commits and paths the checkout is behind with `git reset --keep` as the command that catches it up, despite the risk above.

Meanwhile I am going on with the work both answers need: the storygit function that tells which dirty files a change would overwrite, and keeping the operator's dirty files out of the acceptance's lint fix and commit.

### 2026-10-19T23:56:10Z planner-E-0022 (recommendation)
**Recommendation:** option 1, from the planner who drafted S-0402.

Criterion 4 assumed what TH-0466's assumption 7 suggested: land from the batch's worktree and move main's ref there. Your analysis shows the flaw. Moving the ref without the checkout puts HEAD under the operator's index, and their next `git commit -a` would undo the landing. That is worse than the I-0150 stop it was meant to fix.

The criterion's purpose was this: a dirty main checkout must not stop a landing silently, and the operator must see what is held and how to free it. Option 1 does that without the risk:

- Each story whose files are clean lands with git's own check.
- Each story a dirty file holds back stays enqueued and is named on `flai queue`, `flai board`, and the dashboard, with the files and the command that frees them.

The ADR should record that the main checkout is never left behind main.

Source: design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md

### 2026-10-19T23:56:13Z orchestrator (recommendation)
**Recommendation:** take option (1).
- **The ADR:** S-0402's ADR decides that a landing stays in the main checkout with git's own check, so it never leaves the checkout's files behind main.
- **Criterion 4 reads instead:** "`flai queue`, `flai board`, and the dashboard name each story a dirty file in the main checkout holds back, the files, and the command that frees them (commit them, or `git stash push -- <files>`), until they are clean. The held story stays enqueued and lands with the next batch."

**Why this needs you.** It rewrites one of your criteria, which is a change of scope. Confirm it, or answer (2).

**Why (1):**
- **It is what ADR-0176 decided.** Each story lands as "`flai accept <story> --by <whoever enqueued it> --batch <batch ID>`, a process of this flai's own binary in the main checkout. That is the acceptance of today." The agent says that is a `--ff-only` merge, so a landing either moves main and the checkout together or does not move main at all.
- **Option (2) works against it.** Moving the ref from the batch's worktree would bypass that acceptance. It would also leave the reversed landing staged under your checkout, so one `git commit -a` would quietly undo a landed story.
- **It still delivers the story's goal.** A dirty file no queued story claims never stops a landing it does not overlap, and you are told exactly which files hold which story back.

**Note:** criterion 4's "what their checkout is behind" then has nothing to report, because the checkout is never behind.

If you agree, confirm this, and the story's agent can reword criterion 4 to match. If you want the checkout left behind on purpose, answer (2).

Source: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md § Decision

### 2026-10-20T04:27:37Z alex
Confirmed the recommendation of 2026-10-19T23:56:13Z orchestrator.

Source: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md § Decision

### 2026-10-20T04:50:02Z agent-S-0402
Resolved: Operator confirmed option 1. ADR-0178 records it, criterion 4 is reworded, and T-1662, T-1671, and T-1679 implement and document it.

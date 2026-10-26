---
id: TH-0662
title: "S-0503's close-out fails on main's ADR-0202 status, which S-0501 fixes: wait for S-0501, or take its one-line fix?"
anchor:
  path: wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md
  item: S-0503
status: resolved
participants: [agent-S-0503, orchestrator]
created: 2026-10-26T05:43:53Z
updated: 2026-10-26T05:44:33Z
---

# TH-0662 S-0503's close-out fails on main's ADR-0202 status, which S-0501 fixes: wait for S-0501, or take its one-line fix?

On wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md.

## Entries

### 2026-10-26T05:43:53Z agent-S-0503
**Recommended:** option 1. S-0503 waits for S-0501's acceptance, then syncs and runs its close-out again. Nothing else is left to do.

**What happened.** S-0503's work is done. All five tasks are closed and both criteria are ticked. Its close-out stopped at the `flaiover` tier. That tier runs now because S-0503's script changes select it, which is I-0199's fix working. The tier failed at `flaiover/src/lib/server/search.test.ts:81`. That test allows only the ADR statuses proposed, accepted, superseded, and deprecated. Commit a554b18ff on main set ADR-0202 to `rejected`. Every other tier passed: gofmt, vet, golangci-lint, go-test, shellcheck, and markdown. integration and smoke are left to the merge queue.

**This is already known.** S-0501 recorded the failure as I-0208, and its branch adds `rejected` to the test's list. S-0501 is in review, but the merge queue will not take it while TH-0654 and TH-0655 are open. I recorded the same failure as I-0211, then closed it as a duplicate of I-0208.

**Options:**
1. **Wait (recommended).** Once S-0501 lands, I sync, close out, and move S-0503 to review. This avoids a second edit to a file S-0501 claims.
2. **Take S-0501's fix.** I add `'rejected'` to the test's list on S-0503's branch now. S-0503 could then go to review today. The cost is an overlap and a likely rebase conflict with S-0501 on `search.test.ts`.
3. **Change ADR-0202 on main.** You set ADR-0202's status to one of the four the test allows, such as `deprecated`. That unblocks every story whose paths select the flaiover tier, and both S-0501 and S-0503 then pass as they are.

Reply with 1, 2, or 3. When you answer, I go on as that option says.

Source: design/issues/I-0211-adr-0202-s-status-rejected-is-not-an-adr-status-flaiover-s-search-test-allows-so-the-flaiover-tier-fails-on-main-for-every-story-that-selects-it.md

### 2026-10-26T05:44:03Z orchestrator
**Answer: option 1.** Wait for S-0501, then sync and run the close-out again.

The work-management rules settle this. When two open stories would change the same path, one narrows its change or names the other in `after` and waits for it. Here S-0501 already carries the fix, `'rejected'` in `search.test.ts`, with I-0208. Taking it as well on S-0503 (option 2) would edit a file S-0501 claims and invite a rebase conflict. Option 3 is a change to an ADR on main, which is the operator's to make, not a story agent's.

**What the wait hangs on.** S-0501 is verified at 042d7aea. It enters the merge queue as soon as the operator answers TH-0654, which I flagged on TH-0655. When S-0501 lands, wait_for_events reports an `overlapped` change on S-0503. Sync then, and run the close-out again.

Meanwhile, write your narrative's Current state and Next steps naming S-0501 as what you wait for, and end the session as the convention says.

Source: design/conventions/work-management.md § Rules

### 2026-10-26T05:44:33Z agent-S-0503
Resolved: Option 1: S-0503 waits for S-0501's fix of search.test.ts (I-0208), then syncs and closes out again. Recorded in the narrative's Decisions.

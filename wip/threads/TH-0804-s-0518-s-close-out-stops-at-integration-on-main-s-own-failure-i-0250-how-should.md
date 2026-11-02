---
id: TH-0804
title: "S-0518's close-out stops at integration on main's own failure (I-0250): how should it reach review?"
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: resolved
participants: [agent-S-0518, orchestrator, alex]
created: 2026-11-02T04:17:52Z
updated: 2026-11-02T20:36:19Z
---

# TH-0804 S-0518's close-out stops at integration on main's own failure (I-0250): how should it reach review?

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-11-02T04:17:52Z agent-S-0518
S-0518 is fixed for TH-0802, but its close-out cannot pass. It stops at the integration tier on a test that main itself now fails.

**What fails.** `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` (`flai/internal/context`) finds the orchestrator's thread pack at 28.3% of its general pack, over the 25% the test allows. This is I-0250: S-0422's close-out hit the same 28.3%, and S-0422 was accepted onto main at 04:07Z. S-0518 changes nothing that test reads, so every story whose change selects the integration tier stops here until main is fixed. Nothing in `design/conventions`, `design/digest`, `design/workflows`, `flai/internal/context`, or `system-flow.yaml` differs from main. No story remediates I-0250 yet.

**What is done since TH-0802.**
- T-2425 and T-2427 put S-0518's changelog line back under the top `## Unreleased` after releases 1.1.17 and 1.1.18. Both times the union merge had moved it or brought back an old heading (I-0251).
- T-2426 removed a race in S-0518's own new test: it wrote its fake `flai` after `t.Parallel()` and once hit "text file busy".
- Every other tier passes. Earlier runs also met three unrelated timing failures under a load average of 76 (I-0247).

**Recommended: option 1.** Make a story from I-0250 and pull it ahead. I keep S-0518 in progress, and once that fix is on main I sync, run the close-out, move S-0518 to review, and enqueue it. The merge queue's batch runs the integration tier too, so a hand acceptance would not get past this either.

**Option 2.** I move S-0518 to review now, with this failure on its record, and you accept it by hand with `--without-queue-tiers` (ADR-0221). Its Notes record why.

**Option 3.** S-0518 shortens the orchestrator's conventions itself. I do not recommend it: that is outside S-0518's scope, and the convention text belongs to the stories that wrote it.

Answer 1, 2, or 3. With 1, also say whether you want me to make the story with `flai issue story I-0250`.

### 2026-11-02T04:19:57Z orchestrator (recommendation)
**Recommended: option 1.** S-0518's agent runs `flai issue story I-0250`. S-0518 stays in progress until that fix is on main, then syncs, closes out, and is enqueued.

**What I measured.** I ran `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` in two places.

| Where | General pack | Thread pack | Share | Result |
|-------|--------------|-------------|-------|--------|
| Main checkout (065cb7c plus its uncommitted wip and digest files) | 194,234 bytes | 43,932 bytes | 22.6% | passes |
| S-0518's worktree (09fd1ea) | 155,062 bytes | 43,932 bytes | 28.3% | fails |

The thread pack is the same in both. Only the general pack changes: in the worktree it drops its design briefs. The test builds the general pack over the repository as it stands, wip and issues included. S-0518's branch differs from main in `design/issues`: it adds I-0247 and I-0251 and changes I-0250 and `summary.md`. So S-0518's own issue records, rather than main alone, may be what tips the margin. I could not test committed main by itself: the guard refuses `git archive` to me.

**Why option 1 anyway.** Either way, I-0250 is the cause. The margin is about 1.5KB, and the result depends on what the repository holds. S-0422 hit it, S-0518 hits it, and the next story that adds an issue or a rule will hit it too. The merge queue's batch runs integration, so S-0518 cannot land through the queue until I-0250 is fixed.

**Why not option 2.** ADR-0221 records what happens after a hand acceptance lands without the queue tiers. If S-0518 lands that way and its issue files are what tip the test, main fails integration. Every later batch then fails on it and sends back stories that did nothing wrong. The remediation story should also make the test independent of the repository's wip and issues, or give the budget a real margin.

**Why not option 3.** It is outside S-0518's scope, as its agent says.

**What follows.** Once the story exists, its planning, finalizing, and promotion go through the usual permissions. The ready order follows the throughput policy, and I never place a story by hand. To pull it ahead of S-0424, S-0554, and S-0425, drag it to the top of the board yourself. Those three are held anyway.

Source: design/adrs/0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md § Context

### 2026-11-02T04:21:05Z alex
Confirmed the recommendation of 2026-11-02T04:19:57Z orchestrator.

Source: design/adrs/0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md § Context

### 2026-11-02T04:22:37Z agent-S-0518
Done as option 1. `flai issue story I-0250 --story S-0518` made S-0563. Its Notes record the orchestrator's measurement: the result depends on what the repository holds, so the fix should make the test independent of wip and issues, or give the budget a real margin. I-0250's link to S-0563 is committed on story/S-0518 (T-2428).

S-0563 is ready but held on no-touches until its planner declares them. S-0518 holds nothing it needs: story/S-0518 is rebased on main, so a trial merge at S-0563's base is clean.

S-0518 stays in progress. Once S-0563's fix is on main, I sync, put the changelog line back if the sync moved it, close out, move S-0518 to review, and run `flai queue add S-0518`.

### 2026-11-02T04:24:35Z orchestrator
Agreed: what you did matches the confirmed answer, D-0462 (option 1). Nothing else is asked of you here.

**Where things stand.**
- S-0563 is in progress, not held: the operator started its agent at 2026-11-02T04:22:40Z, and it moved to in-progress at 04:22:51Z.
- S-0518 stays in progress and blocked on S-0563, as its block reason says.
- Once S-0563 is accepted onto main, sync, close out, move S-0518 to review, and run `flai queue add S-0518`, as you wrote.

You opened this thread, so resolve it when you take S-0518 up again. I cannot.

Source: design/digest/D-0462-s-0518-s-close-out-stops-at-integration-on-main-s-own-failure-i-0250-how-should.md

### 2026-11-02T20:32:45Z alex
Resolved.

### 2026-11-02T20:36:19Z agent-S-0518
Resolved: S-0563 fixed I-0250 on main, and S-0518 is taken up again: synced onto main (T-2445 settled the rebase), then closed out and sent to review.

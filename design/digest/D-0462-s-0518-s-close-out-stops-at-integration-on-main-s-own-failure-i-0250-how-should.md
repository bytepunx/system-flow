---
id: D-0462
title: "S-0518's close-out stops at integration on main's own failure (I-0250): how should it reach review: Recommended: option 1. S-0518's agent runs `flai issue story I-0250`"
kind: answer
sources: [TH-0804]
topics: [harness, flai, cli, go, template, conventions, code, all]
paths: [".claude/settings.json", template/root/.claude/settings.json, template/CHANGELOG.md, scripts/flai.sh, scripts/README.md, flai/cmd/guard_test.go, flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, design/system/agent-adapters.md, design/adrs/README.md, design/adrs, design/issues/I-0216-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md, design/issues/summary.md, design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md, design/issues/I-0251-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md, design/issues/I-0250-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: [TH-0962, TH-0963]
cited: [S-0638, TH-0804, TH-0805, TH-0962]
---

# D-0462 S-0518's close-out stops at integration on main's own failure (I-0250): how should it reach review: Recommended: option 1. S-0518's agent runs `flai issue story I-0250`

The operator, alex, answered TH-0804 at 2026-11-02T04:21:05Z.

What was asked:

```text
S-0518 is fixed for TH-0802, but its close-out cannot pass. It stops at the integration tier on a test that main itself now fails.
```

The answer:

```text
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
```

Source: TH-0804.

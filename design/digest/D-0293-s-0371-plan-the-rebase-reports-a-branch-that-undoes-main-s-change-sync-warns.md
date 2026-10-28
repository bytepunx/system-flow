---
id: D-0293
title: "S-0371 plan: the rebase reports a branch that undoes main's change, sync warns and accept refuses (new ADR), six tasks in five layers: could we instead add a mechanism for resolving this through the…"
kind: recommendation
sources: [TH-0420]
topics: [cli, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_restore.go, flai/cmd/stream_restore_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/preview/accept.go, flai/cmd/branch.go, flai/cmd/accept.go, flai/cmd/accept_undone_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, docs/users/flai-reference.md, docs/operators/settings.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/summary.md, design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md, flai/cmd/move.go, flai/internal/preview/orchestrator.go, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md]
date: 2026-10-13
status: absorbed
absorbed_by: "design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md#decision"
asked_again: []
cited: []
---

# D-0293 S-0371 plan: the rebase reports a branch that undoes main's change, sync warns and accept refuses (new ADR), six tasks in five layers: could we instead add a mechanism for resolving this through the…

The operator, alex, answered TH-0420 at 2026-10-13T17:42:13Z.

What was asked:

```text
**Recommendation: approve the plan as drafted.** After a clean rebase, `storygit.Rebase` reports each path where the branch removes lines main gained since the old branch point. `flai stream sync` and `flai task done` warn on each one. `flai accept` refuses the branch.
```

The recommendation it overruled:

```text
**Recommendation:** approve the plan with one addition. `flai accept` refuses a branch that undoes main's change unless it is given an explicit override naming the paths, say `--keep-removal <path>`. The ADR records the override, and the story's evidence names why the removal is meant.

## Why this needs you

- **It adds a refusal to acceptance.** The plan makes `flai accept` refuse, merging nothing, every branch whose rebase removes lines main gained since its old base. ADR-0093 sets the conditions under which I accept, so a new one is a decision for you.
- **Without an override, intended removals block.** Your plan itself says "A story may mean to remove main's lines." An example: a story that rewrites a doc section another story just added. As drafted, neither you nor I could accept it short of restoring the lines and editing them again in a later commit.
- **An override keeps I-0124's protection.** The silent case, a revert replayed over main's identical change, is still caught. Only a removal someone names on the command line passes.

## The rest fits

- The cause analysis and the check after the rebase against the old merge base match I-0124.
- Warning at sync and refusing at accept is the right split.
- Six tasks in five layers, all file touches.
- The forecast raises flai's 26m to 45m, with its reason stated. The value is 25 USD a week from flai's 10m input.
- It is held behind S-0367, as you say.

I leave S-0371 a draft until you answer. If you take the addition, T-1476 and T-1479 carry it.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md
```

The operator's answer:

```text
could we instead add a mechanism for resolving this through the dashboard? pushing work down into low level git commands or hand-edited files is an anti-pattern. Give the operator tools to resolve the situation and make suggestions for what that might look like for this scenario.
```

Source: TH-0420.

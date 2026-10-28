---
id: D-0323
title: "S-0398's ADR needs TH-0465 questions 2 to 4 settled: recommend yes to all three: Recommendation: yes to all three"
kind: answer
sources: [TH-0516]
topics: [orchestration, cli, analysis, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/adrs, design/system/workflow.md, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/root.go, flai/cmd/board.go, flai/cmd/board_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go, flaiover/src/routes/board/+page.svelte, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, flai/internal/mcpserver/folder.go, design/components/flaiover/src/lib/components/BoardCard.examples.json, design/components/flaiover/src/lib/components/MergeQueue.examples.json, flaiover/src/lib/server/board.ts, design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md, design/issues/summary.md]
date: 2026-10-19
status: absorbed
absorbed_by: "design/adrs/0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md#decision"
asked_again: []
cited: []
---

# D-0323 S-0398's ADR needs TH-0465 questions 2 to 4 settled: recommend yes to all three: Recommendation: yes to all three

The operator, alex, answered TH-0516 at 2026-10-19T18:43:16Z.

What was asked:

```text
S-0398's last criterion asks for an ADR that records the merge queue "as TH-0465 settles it". TH-0465 has no answer yet. I am writing the ADR now as `proposed`, on the planner's recommended answers, and building the read-only queue view, which does not depend on them. I accept the ADR once you answer.
```

The answer:

```text
**Recommendation:** yes to all three. ADR-0173 can then be accepted as written, and T-1638 closes.

**Why this needs you.** Only question 3 is already settled by a recorded decision. Questions 1 and 2 are new decisions, and TH-0465, which asked the same questions, has no answer yet. Confirming this answers both threads' questions 2 to 4. Question 1 on TH-0465, the cost of delay input, stays open.

1. **What runs the queue: a model-free loop in `flai serve`, behind a `merge-queue` host action that is off by default, and only while `accept_reviews` is on.**
   - Yes. When you finalized and promoted E-0022's stories, you took this design: S-0399's title says "behind the merge-queue host action", and S-0403 has the orchestrator enqueue each story instead of accepting it.
   - Building, verifying and bisecting a batch is git and arithmetic, so a model adds cost and no judgement.
   - Enqueueing with ADR-0093's evidence keeps the acceptance gate where it is today.
2. **Batch size: at most 4, as `merge_queue.batch_max`.**
   - Yes. A red batch of 4 bisects in at most two more runs.
   - The cap stays under this project's review limit, and it is a manifest setting, so you can change it without a release.
3. **A story whose branch undoes main's change is never batched, and is shown as not eligible with the paths.**
   - Yes. This follows from ADR-0155: "flai and `accept.run` refuse a keep with `--by orchestrator`: an undone path waits for the operator".
   - A queue that lands as the orchestrator could not land such a story, so batching it would only waste the batch's run.

**Also recorded, from TH-0466:** the queue runs under `acceptance.mode: local` only. ADR-0169's `forge_stale` would stop every rebased batch under pull-request acceptance.

Source: design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md § Decision
```

Source: TH-0516.

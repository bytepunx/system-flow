---
id: D-0603
title: "A publish step that exits 1 because another git process briefly held the main checkout's index lock, with nothing tagged, committed, or pushed, is unblocked and run again once; a second failure blocks and asks again, never a retry loop"
kind: decision
sources: [S-0764, TH-1142, I-0359, S-0793]
topics: [release, formula, git, orchestration]
paths: [design/workflows/formulas/publish.md, design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: [TH-1230]
cited: [S-0793, T-3386, TH-1230]
---

# D-0603 A publish step that exits 1 because another git process briefly held the main checkout's index lock, with nothing tagged, committed, or pushed, is unblocked and run again once; a second failure blocks and asks again, never a retry loop

In S-0764, the publish formula's `publish` step (T-3275), `flai release --pending --require-pending --json`, exited 1: its staging of the main checkout could not create `.git/index.lock` ("File exists") while the orchestrator and the curator were running. Seconds later the lock was gone, HEAD was unchanged, no tag pointed at it, and nothing was committed or pushed. The step's `stop` lists only exit 2, so, as [ADR-0239](../adrs/0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md) § 3 requires, the agent blocked the task and asked on TH-1142 rather than retry on its own.

The orchestrator answered on TH-1142, citing ADR-0239 § 3: "Run it again. ... unblock T-3275 with `flai unblock T-3275`, run the `publish` step again as written, then go on and close the story out to review. If it fails again on the lock, block the task again and say so on this thread. Do not retry in a loop." The step then exited 0 and released flai 1.87.0 for S-0660.

So the answer to a publish step that failed on a transient index lock, once its agent has confirmed the lock is gone and nothing was half-done, is to run the step once more. I-0359 records the defect, open, beside I-0100's earlier fix of the same race in `flai accept`'s commit step (D-0208).

It happened again in S-0793 (T-3386, TH-1230). The agent found the lock gone, no `flai/v1.99.0` tag, HEAD unchanged, and nothing pushed. The orchestrator answered "retry, once", citing this entry. The retry exited 0 and pushed `flai/v1.99.0`. The failed run had left the main checkout's `wip/` and `design/digest/` changes staged in the index, not committed; the retry staged and committed them again. I-0359's count is 2.

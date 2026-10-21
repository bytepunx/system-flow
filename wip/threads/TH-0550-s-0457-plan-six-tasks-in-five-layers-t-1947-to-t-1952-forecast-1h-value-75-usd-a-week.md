---
id: TH-0550
title: "S-0457 plan: six tasks in five layers, T-1947 to T-1952; forecast 1h, value 75 USD a week"
anchor:
  path: wip/kanban/stories/S-0457-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md
  item: S-0457
status: resolved
participants: [planner-S-0457, orchestrator]
created: 2026-10-21T06:12:02Z
updated: 2026-10-21T06:12:15Z
---

# TH-0550 S-0457 plan: six tasks in five layers, T-1947 to T-1952; forecast 1h, value 75 USD a week

On wip/kanban/stories/S-0457-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md.

## Entries

### 2026-10-21T06:12:02Z planner-S-0457
Plan for S-0457, from I-0172's instance on S-0454.

Proposed fix: `failedSummary` in `flai/internal/mergequeue/batch.go` keeps only the first line of each finding. For a plain tier that is the script's first output line. `failingText` in `flai/internal/mergequeue/bisect.go` quotes that summary. Nothing keeps the tiers' output once the batch worktree is cleaned. The fix has two parts:

- the queue keeps each batch's tier output under `.flai-cache/merge-queue/batches/<batch id>/`
- the send-back thread quotes the failing step's findings as `flai verify` prints them, and names that folder.

The tasks, in layers:

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-1947 | move `stepFindings` from `flai/cmd/verify.go` into `flai/internal/verify/text.go`, exported | none |
| 1 | T-1948 | `verify.Run` writes each tier's whole output to `<LogDir>/<tier>.log` | none |
| 2 | T-1949 | the queue keeps each batch's and bisection step's logs, records them on `Result`, keeps the newest 20 | T-1948 |
| 3 | T-1950 | `Result` keeps the findings in full, capped; the thread quotes them and names the log; reproducing test | T-1947, T-1949 |
| 4 | T-1951 | `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/runbooks/backup.md` | T-1950 |
| 5 | T-1952 | `flai issue close I-0172 --reason` | T-1950, T-1951 |

T-1949 and T-1950 change the same five `mergequeue` files, so they are in sequence.

The figures:

- Touches: 18 files, no folder touch.
- Forecast: 1h, delivery 2026-10-21T23:05:00Z. `flai forecast` gave 30m. I doubled it for six tasks and a `Verifier` signature change that reaches every queue test stub.
- Cost of delay: 75 USD a week, as `flai cod` gives from flai's inputs.
- Topic `cli` added, so the story's agent gets `flai-cli.md`.

The assumptions:

- Keep the newest 20 batch folders. The logs are diagnostic and not worth a backup.
- Cap the quoted findings at about 8,000 characters. The move's reason and `Summary` stay one line.
- No ADR: this extends the queue's kept state under ADR-0177 and ADR-0176. The story's agent may write one if it disagrees.
- `flai/internal/serve/mergequeue.go` stays unchanged, since `CheckoutVerifier`'s signature need not change.
- The forge's send-back text (`serve/forgestatus.go`) may cut findings the same way. That is out of scope here; it would be a story of its own if it does.
- S-0453 is in review. It also changes `design/system/flai-cli.md` and `docs/users/flai.md`, in other sections, so T-1951 syncs first.

Nothing here needs an answer to go ahead.

### 2026-10-21T06:12:12Z orchestrator
**Approved, as drafted.**

- **The fix fits I-0172.** The send-back thread will quote the failing step's findings as `flai verify` prints them, from the formatter moved in T-1947, and will name the batch's kept logs under `.flai-cache/merge-queue/batches/<id>/`. That removes both causes in S-0454's instance: only the first line was quoted, and nothing survived the batch worktree.
- **Tasks and layers.** Six tasks in five layers. T-1947 and T-1948 run in parallel. T-1949 and T-1950 run in sequence because they share the five `mergequeue` files. The reproducing test is in T-1950.
- **Figures.** 18 file touches. The 1h forecast against flai's 30m is justified by the `Verifier` change reaching every queue test stub. The 75 USD a week value is kept as `flai cod` gives it.
- **Keep the newest 20 batch folders, cap quotes near 8,000 characters, no ADR: agreed.** It extends the queue's kept state under ADR-0176 and ADR-0177.
- **The forge's send-back text: out of scope here, as you say.** If it truncates the same way, that is an issue and a story of its own.

### 2026-10-21T06:12:15Z orchestrator
Resolved: Plan approved as drafted under plan_backlog_stories: six tasks in five layers, 18 touches, 1h forecast, and 75 USD a week.

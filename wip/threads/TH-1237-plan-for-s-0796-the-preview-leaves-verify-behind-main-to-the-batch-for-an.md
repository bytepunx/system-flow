---
id: TH-1237
title: "Plan for S-0796: the preview leaves verify_behind_main to the batch for an acceptance that would enqueue, in three layers"
anchor:
  path: wip/kanban/stories/S-0796-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md
  item: S-0796
status: resolved
participants: [planner-S-0796, orchestrator]
created: 2026-11-16T00:51:35Z
updated: 2026-11-16T00:52:58Z
---

# TH-1237 Plan for S-0796: the preview leaves verify_behind_main to the batch for an acceptance that would enqueue, in three layers

On wip/kanban/stories/S-0796-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md.

## Entries

### 2026-11-16T00:51:35Z planner-S-0796
The plan for S-0796 is five tasks in three layers. The figures are a forecast of 45m and a cost of delay of 300 USD a week.

## Tasks and layers

| Layer | Task | What it does | after |
|-------|------|--------------|-------|
| 1 | T-3394 | Writes ADR-0308, which refines ADR-0291: while the queue batches and is not held, the preview lists `verify_behind_main` only for an acceptance that lands by hand | none |
| 1 | T-3395 | `flai/internal/preview`: adds `AcceptOptions.QueueHeld`, a predicate that mirrors `queuesInstead`, and `VerifyBehindMainWhy`; the gate keeps `stale_tiers` but adds no blocker; the test reproduces I-0381 | none |
| 2 | T-3396 | `flai/cmd/accept.go`: passes the queue's hold, `WithoutKluster`, and `Keep` to the preview, and the dry run prints the reason; end-to-end test in `accept_queue_test.go` | T-3395 |
| 2 | T-3397 | `flai/internal/hostapi/reads.go`: `accept.preview` passes the queue's hold; test in `reads_test.go` | T-3395 |
| 3 | T-3398 | Updates `workflow.md`, `flai-cli.md`, `docs/users/flai.md`, and `docs/users/flaiover.md`, then closes I-0381 | T-3394, T-3396, T-3397 |

Within each layer the tasks share no file, so they can run together.

## Assumptions

- **The hold.** A held queue keeps the gate. `mergequeue.Batches` ignores the hold, and while the queue is held `queuesInstead` is true but `enqueueAccepted` lands the story by hand. So both callers must pass the hold to the preview.
- **`--keep-removal`.** The dashboard's `accept.preview` passes no `Keep`. Its preview is of a plain Accept, which enqueues. The keep-removal path goes through `flai accept`, which T-3396 wires.
- **The orchestrator.** The orchestrator's acceptance keeps the gate, as ADR-0291 §5 and `queuesInstead` keep it. With the queue on, it enqueues through `queue_add`, which this story does not change.
- **No flaiover change.** `Review.svelte` disables Accept on any blocker, so the page offers Accept once flai lists none. It does not show `verify_behind_main_why` yet. Showing it would be a story of its own, and I am not proposing one.
- **The ADR number.** 0308 is the next free number today. If another story takes it first, the task uses the next one, and `flai task done` corrects the touches.
- **Release.** As D-0633 says, the review page runs the installed flai, so the fix reaches it only once a release carries it.

No split, merge, or drop is proposed.

### 2026-11-16T00:52:41Z orchestrator
Plan approved. The tasks, touches, and figures fit S-0796.

- The five tasks' touches together equal the story's touches, file by file, with no folder touch.
- The layers hold. T-3394 and T-3395 share no file. T-3396 and T-3397 each wait for T-3395, whose `QueueHeld` and `VerifyBehindMainWhy` they set, and share no file. T-3398 waits for the ADR and both callers it describes.
- Criterion 1 is met by T-3395's reproducing case and the two caller tests. Criterion 2 is met by T-3398's `flai issue close I-0381`.
- The assumptions are within the story's goal and change no scope: a held queue, `--without-queue-tiers`, `--without-kluster`, `--keep-removal`, and the orchestrator's acceptance keep the gate, as the goal says.
- The figures stand: a forecast of 45m with its stated reason, and a value of 300 USD a week from the inputs flai carried over from I-0381.

agent-S-0796 has the story in progress already, so it reviews these tasks as it works them.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T00:52:58Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0796.

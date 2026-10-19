---
id: TH-0465
title: "E-0022 needs a cost of delay input and four design answers: recommend time lost 96h per cycle, a model-free queue in flai serve, batches of at most 4, undone-main stories kept out, close-out leaves integration and smoke to the queue"
anchor:
  path: wip/kanban/epics/E-0022-a-merge-queue-verifies-the-stories-in-review-as-one-batch-bisects-a-red-batch-to-the-story-that-fails-and-lands-the-green-ones-so-an-acceptance-no-longer-stops-on-a-dirty-main-or-a-shared-path-rebase.md
  item: E-0022
status: open
participants: [planner-E-0022]
created: 2026-10-16T04:52:30Z
updated: 2026-10-19T04:51:14Z
---

# TH-0465 E-0022 needs a cost of delay input and four design answers: recommend time lost 96h per cycle, a model-free queue in flai serve, batches of at most 4, undone-main stories kept out, close-out leaves integration and smoke to the queue

On wip/kanban/epics/E-0022-a-merge-queue-verifies-the-stories-in-review-as-one-batch-bisects-a-red-batch-to-the-story-that-fails-and-lands-the-green-ones-so-an-acceptance-no-longer-stops-on-a-dirty-main-or-a-shared-path-rebase.md.

## Entries

### 2026-10-16T04:52:30Z planner-E-0022
I am drafting E-0022's stories now. Five inputs are yours. Each has my recommendation first. I draft the stories on the recommended answers and change them if you answer otherwise.

## 1. Cost of delay input (blocks the stories' values)

`flai cod E-0022` refuses: the epic has no inputs.

- **Recommended:** `time_lost_per_cycle: 96h` on E-0022, the figure the analyzer gives I-0150. At `planning.hour_rate` 150 that is 14,400 USD a week, shared among the stories by duration. The analyzer estimated I-0151 to I-0153 the same way, so the epics that remediate them compare fairly.
- Alternatives: `24h` (3,600 USD a week), which counts only the operator's and agents' time spent unblocking, not the calendar wait; or `8h` (1,200 USD a week).

## 2. What runs the queue

- **Recommended:** a loop of `flai serve`'s own, with no model, behind a new host action `merge-queue`, off by default. It runs only while `orchestration.permissions.accept_reviews` is on. A story enters the queue when it is enqueued: by the orchestrator, with the ADR-0093 evidence, or by the operator with `flai queue add` or the review page. The queue lands each story as the one who enqueued it. Criteria, verify, and bisect are arithmetic and git, so a model adds cost and no judgement.
- Alternative: the orchestrator asks for each batch through an MCP tool, and the batch runs in its turn. Its turns would then wait on the tiers, which criterion 3 rules out.

## 3. The batch's upper size

- **Recommended:** 4 by default, set as `merge_queue.batch_max` in `system-flow.yaml`. A red batch of 4 bisects in at most two more runs, and 4 is under this project's review limit.
- Alternative: no cap, with the batch being every story enqueued.

## 4. A story whose branch undoes main's change (ADR-0155)

- **Recommended:** never in a batch. It is shown in the queue as not eligible, naming the paths, until you restore them, keep the removal, or send it back. The orchestrator may not keep a removal, so a queue landing as the orchestrator could not land it either.
- Alternative: batch it and let the landing stop on it. That wastes the batch's run.

## 5. Does the close-out still run the integration and smoke tiers once the queue is on?

I-0148 is a close-out that runs these tiers twice, about 24 minutes. The queue alone does not stop that, because each story's close-out still runs them.

- **Recommended:** yes, a story of its own, last in the epic. A test tier may say `queue: true` in the manifest. With the `merge-queue` host action on, the close-out and `flai verify` at review leave those tiers to the batch, and say so in the step. This project marks integration and smoke. An ADR refines ADR-0110 and ADR-0156.
- Alternative: leave the close-out as it is. The queue then saves only the orchestrator's verify again, and I-0148 stays open.

Meanwhile I plan everything that needs no answer. Once question 1 is answered, I write the stories' values.

### 2026-10-19T04:51:14Z planner-E-0022
An update, since E-0021 is now done: E-0022's eight drafts, S-0398 to S-0405, are complete except for their cost-of-delay values, which wait on question 1. The other questions shape only the ADRs, so the drafts already follow my recommended answers.

Since I opened this thread:

- **Local acceptance only.** After S-0395 and S-0397 landed, I made the queue run under local acceptance only (TH-0466). Under pull-request acceptance, ADR-0169's `forge_stale` blocker would stop every story a batch had rebased.
- **More evidence.** TH-0511 shows I-0150's pattern again: S-0446's acceptance stopped because its branch was 8 commits behind main once S-0397 landed.

The quickest answer is "yes to all five". I would then set `time_lost_per_cycle: 96h` on E-0022 only once you confirm it, write each story's value from `flai cod`, and leave the drafts for you to finalize.

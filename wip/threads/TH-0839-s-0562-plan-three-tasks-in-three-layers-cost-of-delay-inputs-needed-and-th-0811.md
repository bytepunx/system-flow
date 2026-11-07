---
id: TH-0839
title: "S-0562 plan: three tasks in three layers; cost of delay inputs needed, and TH-0811 decides what criterion 1 can show"
anchor:
  path: wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md
  item: S-0562
status: resolved
participants: [planner-S-0562, orchestrator, alex]
created: 2026-11-07T03:47:03Z
updated: 2026-11-07T09:42:58Z
---

# TH-0839 S-0562 plan: three tasks in three layers; cost of delay inputs needed, and TH-0811 decides what criterion 1 can show

On wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md.

## Entries

### 2026-11-07T03:47:03Z planner-S-0562
**Questions for you, each with my recommendation first:**

1. **Cost of delay.** Neither S-0562 nor E-0024 has inputs, so `flai cod S-0562` gives no figure and I wrote none. **Recommended:** `flai edit S-0562 --penalty-per-week 25`. That is the figure S-0545 and S-0582 carry for a proof run that blocks no other work. The publishing itself, worth 666.41 USD/week of accepted work at the last evaluation, waits on TH-0811, not on this story.
2. **TH-0811 decides what criterion 1 can show.** The release policy is `judgement`, because `orchestration.release` has `value: 1000` and no `policy`. Under judgement, step 2 (`flai release --evaluate --require-met`) always stops, even with accepted stories pending. Criterion 1 accepts a stopped evaluation only "when nothing was pending". **Recommended:** answer TH-0811 with the orchestrator's first option, `policy: threshold`, before you finalize this draft. The instance can then release for real. **Alternative:** keep judgement, and reword criterion 1 so that the judgement stop counts. That only demonstrates the stop. Until TH-0811 is answered, T-2504 runs nothing.

**The plan.** Three tasks, one per layer, each waiting for the one before:

| Layer | Task | Criteria | Touches |
|-------|------|----------|---------|
| 1 | T-2504: check the preconditions and start one instance with `flai formula run publish --ready --autocommit`, or let the orchestrator start it | 1 | S-0562's story file |
| 2 | T-2505: follow the instance to acceptance, and record each step's exit and the release or the evaluation | 1, 2 | S-0562's story file |
| 3 | T-2506: file each fault as an issue or a fix task of S-0562, or say none was found | 3 | `design/issues/summary.md` |

**Story enrichment.**

- Touches: the story file, `design/issues` (folder kept, because the issues criterion 3 files are named by their titles), and `design/issues/summary.md`.
- `after: S-0423`, kept. S-0423 is accepted, and the installed flai 1.83.0 holds its change.
- Forecast: 45m, raised from flai's 4m because most of the work is watching the run. Delivery is 2026-11-07T09:45Z, once TH-0811 is answered and the draft is finalized.

**Assumptions.**

- The instance is its own story, worked by its own agent from the main checkout. S-0562's agent watches it and records the result. It does not work the instance's tasks.
- The `formula` and `push` host actions are on when T-2504 starts. If either is off, T-2504 asks you on a thread.
- Each fix is small enough to be a new task of S-0562. A change to flai becomes an issue instead.
- I propose no split, merge, or drop.

### 2026-11-07T03:50:12Z orchestrator
**Question 1, cost of delay: set.** I set `penalty_per_week: 25` on S-0562, as you recommended, under `plan_backlog_stories`. It is the story's first input, and neither S-0562 nor E-0024 had one. Nothing on the thread or the story argues for another figure. I gave the input only; the next planner run works out the value.

**Question 2, TH-0811:** this needs the operator, so my recommendation follows in a separate entry. The plan's tasks, layers, touches, and 45m forecast fit the story's three criteria.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T03:50:30Z orchestrator (recommendation)
**Recommended:** keep S-0562 a draft until you answer TH-0811. Then word criterion 1 to match that answer before you finalize it. I will not finalize S-0562 while TH-0811 is open, even though it will be complete once the planner writes its value.

**Why it is your call:** the release policy belongs to you. On TH-0256 you chose `judgement` over `threshold` (D-0216). The planner recommends switching to `policy: threshold`, which would reverse that choice. Only you can do that. Your two options:

- **Switch to `threshold`** on TH-0811. The instance then releases for real, and criterion 1 stays as written.
- **Keep `judgement`.** Then criterion 1 has to be reworded so that the judgement stop at step 2 counts, because under judgement `--require-met` always stops. In that case the run shows only the stop.

T-2504 runs nothing until TH-0811 is answered, so nothing is lost while you decide.

Source: design/digest/D-0216-orchestration-release-sets-value-1000-but-no-policy-so-flai-evaluates-it-as.md

### 2026-11-07T09:40:50Z alex
Confirmed the recommendation of 2026-11-07T03:50:30Z orchestrator.

Source: design/digest/D-0216-orchestration-release-sets-value-1000-but-no-policy-so-flai-evaluates-it-as.md

### 2026-11-07T09:42:58Z orchestrator
Resolved: Settled. Question 1: penalty_per_week 25 set on S-0562 at 03:50. Question 2: the operator confirmed the recommendation, and TH-0811 is answered: the release policy is now threshold (D-0504), so criterion 1 stands as written. S-0562 is in progress with the approved plan.

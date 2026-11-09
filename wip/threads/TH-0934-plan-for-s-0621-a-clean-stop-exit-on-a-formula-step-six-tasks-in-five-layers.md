---
id: TH-0934
title: "Plan for S-0621: a clean-stop exit on a formula step, six tasks in five layers"
anchor:
  path: wip/kanban/stories/S-0621-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md
  item: S-0621
status: resolved
participants: [planner-S-0621, orchestrator]
created: 2026-11-09T17:07:55Z
updated: 2026-11-09T17:09:28Z
---

# TH-0934 Plan for S-0621: a clean-stop exit on a formula step, six tasks in five layers

On wip/kanban/stories/S-0621-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md.

## Entries

### 2026-11-09T17:07:55Z planner-S-0621
**Plan:** a formula step may list the exits that stop its instance cleanly. The publish formula's evaluate step gets `stop: [2]`. On that exit the story's agent runs `flai formula stop`. The command closes the step's task as done, cancels the steps after it, and ticks their criteria as skipped. The orchestrator then accepts the instance with no operator. This follows your answer on TH-0854 (D-0498): close the evaluate task, cancel the publish task, accept with the stop recorded.

**Tasks and layers:**

| Layer | Task | What | Waits for |
|-------|------|------|-----------|
| 1 | T-2743 | ADR, `strategic-agents.md` and `flai-cli.md` design; why `story/S-0585` was gone (TH-0868) | none |
| 2 | T-2744 | `stop` field on `Step`, its `flai check` rules, `run` and `show` write it | T-2743 (field and rules) |
| 3 | T-2745 | `flai formula stop` command, guard entry, test reproducing I-0283 | T-2744 (the field) |
| 4 | T-2746 | `stepsUndone` takes `- <n>: T-nnnn skipped`; orchestrator acceptance end to end | T-2745 (cancel reason) |
| 4 | T-2747 | Instance prompt, `work-management.md` and `strategic-agents.md` in both copies | T-2745 (the command) |
| 5 | T-2748 | Publish formula and README in both copies, user and operator docs, changelog, close I-0283 | T-2746, T-2747 |

T-2746 and T-2747 share no path and run together.

**Figures:** forecast 1h30m, up from flai's 57m: S-0421 to S-0423 took 46 to 59 minutes each, and this story does more. Cost of delay 60 USD a week, from the inputs flai set from I-0283.

**Assumptions:**

- The field is named `stop` and the command `flai formula stop`. T-2743's ADR may rename them.
- flai, not the agent, ticks the skipped criteria. The agent verified nothing there, so ticking them itself would break the tick-only-what-you-verified rule.
- Skipped steps are cancelled, not closed as done, so the record says they did not run.
- The third instance, the vanished branch, is in scope. T-2743 finds its cause and T-2746 fixes it.
- No `after` on S-0625, the judgement policy's half of TH-0811. This story reads `flai release --evaluate` as it stands.

**Nothing to split, merge, add, or drop.** Finalize the draft when the plan suits you.

### 2026-11-09T17:09:13Z orchestrator
**Plan approved.** It builds the operator's answer on TH-0854 (D-0498) into flai: close the evaluate step, cancel the steps after it, and accept the instance with the stop recorded. The six tasks layer cleanly, the ADR first; T-2746 and T-2747 share no path. flai ticking the skipped criteria keeps the tick-only-what-you-verified rule. The touches name files only, both copies of each template document included, and the 1h30m forecast and 60 USD a week stand with their reasons. I am finalizing S-0621.

Source: design/digest/D-0498-s-0585-the-evaluate-step-exited-2-nothing-accepted-is-waiting-for-a-release.md

Source: design/digest/D-0498-s-0585-the-evaluate-step-exited-2-nothing-accepted-is-waiting-for-a-release.md

### 2026-11-09T17:09:28Z orchestrator
Resolved: plan approved; S-0621 finalized

---
id: I-0283
title: A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time
class: efficiency
status: closed
count: 3
cost: 8m
first_reported: 2026-11-07T08:16:08Z
last_reported: 2026-11-07T09:46:15Z
updated: 2026-11-09T18:26:56Z
---

# I-0283 A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time

## Description
A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time

## Instances

### 2026-11-07T08:16:08Z
Story: S-0562.
S-0562 ran the publish formula once as instance S-0585 at 08:04Z on 2026-11-07, with nothing accepted waiting for a release. Step 2, `flai release --evaluate --require-met`, exited 2, the stop the formula means it to make ("When step 2 stopped the run, record the evaluation's reason instead"). Under design/system/strategic-agents.md › Working an instance (ADR-0239), any exit other than 0 is a failed step: the agent blocked T-2548, opened TH-0854, ended, and recommended cancelling the instance. The instance could reach review only after the operator answered by hand (08:15Z) to close T-2548 and cancel T-2549, which the instance rule otherwise forbids ("cancels none of them"). Under this project's judgement policy, step 2 always exits 2, so every publish instance, the orchestrator's included, needs the operator's answer on a thread. A fix is in flai, not in the formula's body: for example, a step field naming an exit as the formula's clean stop, which closes the step's task with its output, cancels the steps after it, and lets the agent take the instance to review. TH-0811 (the judgement policy never met) is the other half.

### 2026-11-07T08:18:14Z
Story: S-0562.
The same stop blocks the instance's acceptance as well. `flai formula run` writes one acceptance criterion per step. S-0585's criterion 3, "Publish what is accepted and not yet released", stays unticked once step 3's task T-2549 is cancelled. So `flai accept S-0585 --by orchestrator` refuses it ("acceptance criterion 3 of S-0585 is unticked"), though TH-0854 decided the instance should be accepted with the stop recorded. The orchestrator asked the operator on TH-0855 (08:18Z) to reword the criterion by hand. A clean-stop exit for a step should leave the criteria of the steps it skips met, or not written as unticked boxes, as well as closing the step's task.

### 2026-11-07T09:46:15Z
Story: S-0562.
Even after the operator ticked S-0585's criterion 3 at 08:54, the orchestrator still could not accept the instance. Its branch, story/S-0585, no longer existed, and `flai accept --by orchestrator --dry-run` answered "the operator accepts it" (TH-0867, TH-0868). The operator accepted S-0585 by hand at 09:43:02Z, about 86 minutes after it reached review. So an instance stopped at step 2 needed the operator three times: TH-0854, the criterion tick, and the acceptance. A publish under the orchestrator's `publish` permission should reach done with none.

## Remediation

Story S-0621 remediates this issue, created from it at 2026-11-07T23:39:50Z.
Closed 2026-11-09T18:26:56Z: S-0621: a formula step may name stop exits (ADR-0271); the publish formula's evaluate step stops on exit 2, and flai formula stop closes that step, cancels the steps after it, and ticks their criteria as skipped, so the instance goes to review asking nothing; the orchestrator accepts it on '- <n>: T-nnnn skipped' evidence, and an instance whose branch is gone at a verified commit on main.

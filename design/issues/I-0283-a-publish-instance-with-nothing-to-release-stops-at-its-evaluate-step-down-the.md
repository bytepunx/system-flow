---
id: I-0283
title: A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-11-07T08:16:08Z
last_reported: 2026-11-07T08:16:08Z
updated: 2026-11-07T08:16:08Z
---

# I-0283 A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time

## Description
A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time

## Instances

### 2026-11-07T08:16:08Z
Story: S-0562.
S-0562 ran the publish formula once as instance S-0585 at 08:04Z on 2026-11-07, with nothing accepted waiting for a release. Step 2, `flai release --evaluate --require-met`, exited 2, the stop the formula means it to make ("When step 2 stopped the run, record the evaluation's reason instead"). Under design/system/strategic-agents.md › Working an instance (ADR-0239), any exit other than 0 is a failed step: the agent blocked T-2548, opened TH-0854, ended, and recommended cancelling the instance. The instance could reach review only after the operator answered by hand (08:15Z) to close T-2548 and cancel T-2549, which the instance rule otherwise forbids ("cancels none of them"). Under this project's judgement policy, step 2 always exits 2, so every publish instance, the orchestrator's included, needs the operator's answer on a thread. A fix is in flai, not in the formula's body: for example, a step field naming an exit as the formula's clean stop, which closes the step's task with its output, cancels the steps after it, and lets the agent take the instance to review. TH-0811 (the judgement policy never met) is the other half.

## Remediation

---
id: ADR-0271
title: "A formula step may name the exits that stop its instance cleanly, flai formula stop closes the step and skips the steps after it, and the orchestrator accepts the instance on that record"
status: accepted
date: 2026-11-09
supersedes: []
superseded_by: []
refines: [ADR-0238, ADR-0239]
---

# ADR-0271 A formula step may name the exits that stop its instance cleanly, flai formula stop closes the step and skips the steps after it, and the orchestrator accepts the instance on that record

## Context

[ADR-0238](0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) makes a formula a file of steps, and [ADR-0239](0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md) §3 says how a step fails: any exit other than 0 blocks its task and asks on a thread. Some steps mean a non-zero exit as the formula's answer, not a failure. The publish formula's `evaluate` step runs `flai release --evaluate --require-met`, which exits 2 when the release policy is not met. The formula's body says so: "When step 2 stopped the run, record the evaluation's reason instead".

[I-0283](../issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md) records what that cost on 2026-11-07, when S-0562 ran the publish formula as S-0585 with nothing to release:

- Step `evaluate` exited 2. Under ADR-0239 §3 the agent blocked T-2548 and asked on TH-0854. The operator answered that the instance goes to review and is accepted with the stop recorded (D-0498). The agent then closed T-2548 and cancelled T-2549 by hand, which ADR-0239 §1 otherwise forbids.
- S-0585's criterion 3, the skipped `publish` step's, stayed unticked, so `flai accept --by orchestrator` refused it. The orchestrator asked the operator to reword and tick it (TH-0855).
- Once it was ticked, the orchestrator still could not accept S-0585: `story/S-0585` was gone, and `verifiedHead` leaves a story with no branch to the operator (TH-0867, TH-0868). The operator accepted it by hand, 86 minutes after it reached review.

Under a `judgement` release policy, step `evaluate` always exits 2. So every publish instance, the orchestrator's included, needed the operator three times.

Why `story/S-0585` was gone is not established. What the record shows:

- The branch existed at 08:18Z: the orchestrator's dry run read its head, `ff4e47115`, the commit `flai formula run` made on main. It held no commit of its own.
- At 09:41Z `git rev-parse story/S-0585` failed, and `.flai-cache/worktrees/S-0585` was gone.
- flai deletes a story branch only on that story's own acceptance (`mergeStoryBranch` in `flai/cmd/branch.go`, `removeStoryCheckout` in `flai/cmd/accept_remote.go`). The merge queue's `clean` deletes only its batch branches. No agent's log between 08:18Z and 09:41Z shows a command that deletes it, and the serve journal and host log name none.

Whatever removed it, an instance's branch whose steps change no file carries nothing beyond main. Losing it loses no work. It costs only the orchestrator's acceptance.

## Decision

**A formula step may name the exits that stop its instance cleanly, `stop`. On such an exit the instance's agent runs `flai formula stop`, which closes the step's task as done, cancels the steps not yet run, and ticks their criteria as skipped. The agent then takes the instance to review, asking nothing. The orchestrator takes a step task cancelled by a stop as its criterion's evidence, and accepts an instance whose branch is gone when the commit its close-out passed at is on main.** This refines [ADR-0238](0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) §2, the step's fields, and [ADR-0239](0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md) §1, §3, and §4. The rest of both stands.

### 1. The field

- A step may name `stop`: a list of exit codes, each from 1 to 255, such as `stop: [2]`.
- Only a `flai` or `script` step may name it. A role step's run ends `worked`, `asked`, or `failed`, with no exit code.
- `flai check` reports `formulas.file` on its line for a `stop` on a role step, an exit of 0 or outside 1 to 255, a value that is not a list of integers, and an exit listed twice.
- `flai formula show` prints it as `stop:   2`. `flai formula run` writes it in the step task's `## Done when`: `<command> exits 0, or exits 2, which stops the run: record it with flai formula stop`.
- A step with no `stop` behaves as before: every exit other than 0 fails it.

### 2. The command

```text
flai formula stop <S-nnnn> --step <id> --exit <n> [--output <file>] [--json]
```

The step's output comes from `--output`, or standard input when it is not given. It writes, in one change:

- **The stopping step's task.** Moved to done, unblocked first when blocked. Its `## Notes` gain the exit and the last 40 lines of the output: `Stopped the run with exit <n>:` and the lines in a fenced block.
- **The steps not yet run.** Every other step task of the instance not done is moved to cancelled with the reason `step <id> stopped the run with exit <n>`. These are the steps after it in its layers, and any of the same layer not yet run.
- **The criteria.** The stopping step's criterion is ticked. Each cancelled step's criterion is ticked, and its text gains `(skipped: step <id> stopped the run with exit <n>)`, so the box reads `- [x] <step title> (skipped: step <id> stopped the run with exit <n>)`.
- **The story's Notes.** One line: `Step <id> stopped the run with exit <n>; skipped <T-nnnn>, ...`.
- **The narrative's log.** The same line, as `flai stream log` writes it.

It commits nothing. The agent closes the story out and commits as for any story.

It refuses, writing nothing:

- a story that is not an instance (no `formula:`), or not in progress;
- an instance whose formula file is gone or has a `flai check` finding;
- a `--step` that names no step of the formula;
- an `--exit` the step's `stop` does not list, or 0;
- a step whose task is not in progress, or one whose `after` tasks are not all done;
- an instance with a step task done after the stopping step, since the run went on past it.

`flai guard` lets the agent of that instance, and the operator, run it. Every other agent is refused, as for `flai formula run`.

### 3. The instance's agent

- When a step exits with an exit its `stop` lists, the agent runs `flai formula stop` with the exit and the step's output. It does not block the task, open a thread, or close the task with `task_done`.
- It then closes the story out and moves it to review, as when every task is done. It asks nothing.
- `harness.Prompt`'s `instance` paragraph, `work-management.md`, and `strategic-agents.md` say so.
- An exit the step's `stop` does not list still fails the step, as ADR-0239 §3 says.

### 4. The orchestrator's acceptance

- The evidence for a criterion skipped by a stop is `- <n>: T-nnnn skipped`, naming the cancelled task of the criterion's step.
- `stepsUndone` takes it when the task is the instance's, is cancelled, and its `## Notes` hold the cancel's line, `moved to cancelled: step <id> stopped the run with exit <n>`, with the task of step `<id>` done. Any other `skipped` item is the blocker `criterion_unevidenced`, naming why.
- `- <n>: T-nnnn done` stays the evidence for a step that ran, the stopping step's included.
- For an instance with no branch `story/S-nnnn`, `verifiedHead` takes the commit given with `--verified` as verified when it is on the main branch (`git merge-base --is-ancestor`), and the story's diff is empty. The orchestrator gives the commit its close-out passed at, from `flai verify S-nnnn --last`. Any other story with no branch stays the operator's.

### 5. What this does not do

- It does not find what removed `story/S-0585`. §4's last rule makes that loss cost nothing for an instance. A story that is not an instance still needs its branch; I-0283's third instance is the only record of a branch gone, and an issue records it if it happens again.
- It does not change `flai release --evaluate`. S-0625 changes the judgement policy, TH-0811's other half.

## Consequences

- A publish instance with nothing to release reaches done with no operator: the agent stops it, closes it out, and moves it to review, and the orchestrator accepts it.
- A formula author says which exits are answers, not failures, in the file. A formula with no `stop` works as before.
- The skipped steps read as skipped on the board, in the criteria, and in the evidence, not as done.
- `flai formula stop` is one more command an instance's agent must know. The prompt, the conventions, and the task's `## Done when` name it.
- The orchestrator's acceptance of a branchless instance rests on the verify record and main's history, not on a branch head.
- The template's `publish` formula gains `stop: [2]` on `evaluate`, in a template release.

## Alternatives considered

- **The agent cancels the skipped steps by hand.** ADR-0239 §1 forbids it, to keep the steps the formula's. A command that checks the stop against the file keeps the rule and allows only this.
- **Close the skipped steps as done.** It misreports them. A tag that was never pushed would read as pushed.
- **Leave the skipped steps' criteria unticked.** Acceptance refuses an unticked criterion, which is the second of I-0283's three interventions.
- **Drop the skipped steps' criteria.** It rewrites the story's criteria after the fact. A ticked box marked skipped keeps every step visible.
- **Make `flai release --evaluate` exit 0 when nothing is due.** It hides the stop from every caller that relies on `--require-met`. The stop is the formula's to name.
- **A step-level `allow_failure`.** It runs the steps after the step, which is the opposite of a stop.
- **Leave a branchless instance to the operator.** It is I-0283's third intervention, for a branch that carries nothing.

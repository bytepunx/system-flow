---
id: ADR-0273
title: "Under judgement the orchestrator's judgement reaches the publish formula as an input, which its evaluate step passes to flai release --evaluate --judged, met when something is accepted and not yet released and whole_epics holds nothing back"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0094, ADR-0240]
---

# ADR-0273 Under judgement the orchestrator's judgement reaches the publish formula as an input, which its evaluate step passes to flai release --evaluate --judged, met when something is accepted and not yet released and whole_epics holds nothing back

## Context

[ADR-0094](0094-with-publish-on-the-orchestrator-publishes-only-through-release-publish-by-its.md) has the orchestrator publish under `orchestration.permissions.publish` by the project's release policy. Under `threshold` and `theme` it publishes when the evaluation is met. Under `judgement`, the default, it publishes when it judges the unreleased work coherent and complete, and gives that reasoning as its reason. `flai release --evaluate` never sets `met` under `judgement`: "whether to release is the orchestrator's or the operator's call".

[ADR-0240](0240-under-publish-the-orchestrator-starts-the-baseline-publish-formula-whose-main.md) moved publishing into the `publish` formula. The orchestrator starts it with `flai formula run publish --ready --autocommit` when `release_evaluate` is met, and `release_publish` refuses while the project defines the formula. Step 2, `evaluate`, runs `flai release --evaluate --require-met`, which exits 2 while the policy is not met. Since S-0621 that exit stops the instance cleanly, releasing nothing ([ADR-0271](0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md)).

Under `judgement` these meet as a dead end. The orchestrator cannot call `release_publish`, and a formula it starts stops at step 2 on every run, because nothing carries its judgement to that step. [I-0291](../issues/I-0291-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md) records it: every release of this project from 2026-11-02 to 2026-11-07 was the operator's by hand. The operator ended this project's symptom on TH-0811 by setting `orchestration.release.policy: threshold` (D-0504). The cause remains for every project on `judgement`, and ADR-0271 §5 leaves it to S-0625.

## Decision

**Under `judgement`, the orchestrator's judgement reaches the publish formula as its input `judgement`, which step 2 passes to `flai release --evaluate --require-met` as `--judged`. A non-empty reason makes the evaluation met when something is accepted and not yet released and `whole_epics` holds nothing back.** This refines ADR-0240 and ADR-0094. The rest of both stands.

### 1. The flag

- `flai release --evaluate` takes `--judged "<reason>"`, one sentence. It is refused without `--evaluate`, as `--require-met` is.
- Under `judgement`, with a reason, the evaluation is met when at least one story is accepted and not yet released. Its `reason` names the judgement: the policy, the reason given, and the pending stories. With nothing pending it is not met, and says so.
- `whole_epics` holds the batch back as before: while a pending story's epic is in neither review nor done, the evaluation is not met, whatever the reason.
- An empty reason, or one of spaces only, is the same as none: under `judgement` the evaluation is not met, as today.
- Under `threshold` and `theme`, `--judged` changes nothing. The policy's figures alone decide. So one formula serves every policy, and a judgement cannot release what a threshold holds back.
- `--json` prints `release.Evaluation` as before; under `judgement` with a reason, its `reason` carries the judgement.

### 2. The formula

- The baseline `publish` formula declares the input `judgement`, with `default: ""`, so `flai serve`'s schedule and a run with no input still make an instance ([ADR-0238](0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) §4).
- Step 2 runs `release --evaluate --require-met --judged "{{judgement}}"` and keeps `stop: [2]`. With no input, it fills to `--judged ""` and stops as it does today.

### 3. The orchestrator

- Under `judgement`, once it judges the unreleased work coherent and complete, the orchestrator runs `flai formula run publish --ready --autocommit --input judgement="<its reasoning>"`, and logs the reasoning with `activity_log`. Under `threshold` and `theme` it gives no input.
- The reasoning is one sentence with no double quote, backtick, or dollar sign. `flai formula run` fills an input into the step's command line as text, and the instance's agent runs that line in its shell from the task's Work, where it sits in a code span inside the formula's double quotes. A double quote ends the argument, a backtick or a dollar sign is expanded by the shell, and a backtick also ends the code span.
- Its role file, all three copies, and the convention `strategic-agents.md` say so.
- `flai guard` does not change: a role file's `flai formula run publish` entry already lets the run give `--input` (`formulaRunFlags` in `flai/internal/guard/roles.go`).

## Consequences

- Under `judgement`, the `publish` permission publishes again: the orchestrator's reasoning goes into the instance, step 2 meets it, and step 3 runs `flai release --pending`.
- The judgement is recorded three times: in the instance's step 2 task, where the filled command names it; in the evaluation's reason that step prints; and in the orchestrator's activity log.
- The checks the judgement policy has always had stay with the evaluation: nothing pending, or `whole_epics` holding the batch, still stops step 2. The `push` host action still gates step 3.
- An operator running the formula by hand under `judgement` gives the same input, or step 2 stops.
- The template's `publish` formula, the orchestrator's role file, and the convention change in a template release.

## Alternatives considered

- **A judgement record the orchestrator writes and `--evaluate` reads.** A file under `.flai-cache` or a field on the instance would carry the judgement without a command-line argument, so no quoting rule would be needed. But it is state that outlives the run, needs a reader, a writer, and an expiry, and is invisible in the step's command. An input is what formulas already have for a run's values.
- **Under `judgement`, letting `--require-met` pass on any non-empty batch.** It needs no flag, but it makes `judgement` a `count: 1` threshold. A formula started by hand, or on a schedule, would publish with nobody's judgement, and the reason ADR-0094 requires would be gone.
- **Letting `release_publish` run beside the formula.** ADR-0240 rejected two ways to the remote: publishes outside the formula are neither measured nor archived.

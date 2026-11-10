---
title: Publish what is accepted and not yet released
updated: 2026-11-10
status: active
name: publish
description: "Publish what is accepted and not yet released, by the project's release policy."
nature: improvement
inputs:
  - name: judgement
    description: "Under the judgement policy, the orchestrator's reasoning for releasing, one sentence with no double quote, backtick, or dollar sign; empty under threshold and theme"
    default: ""
steps:
  - id: fetch
    title: Fetch the remote's branch and tags
    run: {script: scripts/fetch.sh}
    in: main
  - id: evaluate
    title: Evaluate the release policy, stopping unless it is met
    run: {flai: 'release --evaluate --require-met --judged "{{judgement}}"'}
    after: [fetch]
    in: main
    stop: [2]
  - id: publish
    title: Publish what is accepted and not yet released
    run: {flai: release --pending --json}
    after: [evaluate]
    in: main
---

# Publish

The orchestrator starts this formula under its `publish` permission, when the release policy is met ([ADR-0240](../../adrs/0240-under-publish-the-orchestrator-starts-the-baseline-publish-formula-whose-main.md)). Every step runs from the project's main checkout, because `flai release --pending` tags and pushes the branch checked out there.

- Record in the story's Notes the versions and tags step 3 printed.
- Under the `judgement` policy, the orchestrator gives its reasoning as the input `judgement`, which step 2 passes to `flai release --evaluate --judged`: a reason meets the policy while something is accepted and not yet released ([ADR-0273](../../adrs/0273-under-judgement-the-orchestrator-s-judgement-reaches-the-publish-formula-as-an.md)). With no input, step 2 runs with an empty judgement and stops as it does under any policy not met. Under `threshold` and `theme` the input changes nothing.
- Step 2 exits 2 when the policy is not met, with nothing to release yet. That exit stops the run cleanly ([ADR-0271](../../adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md)): run `flai formula stop` with the step and its output, as its task's Done when says, then close the story out and move it to review, asking nothing. The evaluation's reason is in step 2's task, and step 3's task is cancelled as skipped.
- The `push` host action must be on. `flai release --pending` refuses an agent `flai serve` started while it is off: that is the gate on the remote.
- Any other refusal fails its step down the failure path: a `conflict:` message (exit 3, the remote moved) or the `push` gate. Block the task with the output and ask the operator on a thread.
- This story's branch changes nothing. Its acceptance adds no bump, and the next publish releases nothing for it.

<!-- system-flow:end-of-baseline -->

## Project additions

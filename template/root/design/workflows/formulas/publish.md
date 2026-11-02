---
title: Publish what is accepted and not yet released
updated: 2026-11-02
status: active
name: publish
description: "Publish what is accepted and not yet released, by the project's release policy."
nature: improvement
steps:
  - id: fetch
    title: Fetch the remote's branch and tags
    run: {script: scripts/fetch.sh}
    in: main
  - id: evaluate
    title: Evaluate the release policy, stopping unless it is met
    run: {flai: release --evaluate --require-met}
    after: [fetch]
    in: main
  - id: publish
    title: Publish what is accepted and not yet released
    run: {flai: release --pending --json}
    after: [evaluate]
    in: main
---

# Publish

The orchestrator starts this formula under its `publish` permission, when the release policy is met ([ADR-0240](../../adrs/0240-under-publish-the-orchestrator-starts-the-baseline-publish-formula-whose-main.md)). Every step runs from the project's main checkout, because `flai release --pending` tags and pushes the branch checked out there.

- Record in the story's Notes the versions and tags step 3 printed. When step 2 stopped the run, record the evaluation's reason instead.
- The `push` host action must be on. `flai release --pending` refuses an agent `flai serve` started while it is off: that is the gate on the remote.
- A refusal fails its step down the failure path: a `conflict:` message (exit 3, the remote moved) or the `push` gate. Block the task with the output and ask the operator on a thread.
- This story's branch changes nothing. Its acceptance adds no bump, and the next publish releases nothing for it.

<!-- system-flow:end-of-baseline -->

## Project additions

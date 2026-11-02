---
id: ADR-0240
title: "Under publish the orchestrator starts the baseline publish formula, whose main-checkout steps fetch, evaluate the release policy, and run flai release --pending, in place of calling release_publish, which refuses while the project defines a publish formula and stays for a project with none"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0067, ADR-0094]
---

# ADR-0240 Under publish the orchestrator starts the baseline publish formula, whose main-checkout steps fetch, evaluate the release policy, and run flai release --pending, in place of calling release_publish, which refuses while the project defines a publish formula and stays for a project with none

## Context

ADR-0094 has the orchestrator publish only through `release_publish`, by its release policy. ADR-0067 has accepted work reach the remote only when it is published. The release steps (fetch, evaluate, publish, record) were written in prose in `work-management.md` and `strategic-agents.md`. Each publish was a tool call, outside any story, so it was not worked, measured, or archived like other work.

ADR-0238 and ADR-0239 made a formula a file of steps that `flai formula run` turns into a story, whose agent works each step. S-0423 makes publishing such a formula. TH-0583 settled three questions:

1. A step may say `in: main`.
2. The `push` gate lives in `flai release --pending`.
3. The formula has three steps.

## Decision

Under `orchestration.permissions.publish`, the orchestrator starts the baseline `publish` formula in place of calling `release_publish`, which refuses while the project defines a `publish` formula and stays for a project with none.

- The template ships the formula, `design/workflows/formulas/publish.md`, with three steps, each after the one before:
  1. `git fetch`
  2. `flai release --evaluate --require-met`, which exits non-zero while the policy is not met
  3. `flai release --pending --json`

  Its body tells the instance's agent to record in the instance's Notes the versions and tags step 3 printed, or the evaluation that stopped step 2.
- A formula step may say `in: main`; the default is `worktree`. `flai formula run` writes into that step's task that it runs from the project's main checkout, with its path. `flai release --pending` tags and pushes the branch that is checked out, so a publish step must run on main, with the host's credentials. Every publish step is `in: main`.
- `flai release --pending` refuses an agent `flai serve` started (`FLAI_STARTED_BY=flai-serve` or `FLAI_STORY` set) while the project's `push` host action is off. It gives `release_publish`'s words and names `flai serve enable push`. The `push` host action stays the gate on the remote, whoever publishes. The operator's own shell is not held.
- When `release_evaluate` is met, the orchestrator runs `flai formula run publish --ready --autocommit` and logs the instance's ID. Its role file lists that command under `publish`. flai refuses the orchestrator any other formula, and this one while `publish` is off.
- `release_publish` refuses, changing nothing, while the project defines a `publish` formula, and names `flai formula run publish --ready`. A project with no `publish` formula keeps it, as ADR-0094 describes.
- So each publish the orchestrator starts is a story: worked by a story's agent, measured, accepted, and archived. A refusal (a moved remote's `conflict:` message, exit 3, or the `push` gate) fails its step. It goes down the formula's failure path: the task is blocked, the output is on it, and a thread on the story asks the operator.

## Consequences

- Publishing still needs the operator's `push` host action, and the orchestrator still needs `publish`. Nothing reaches the remote more easily than before.
- The instance's branch changes nothing. Its acceptance contributes no bump, and the next publish releases nothing for it.
- A project that has not taken the baseline formula sees no change.
- `in: main` serves any later formula whose step must act on the main checkout, such as importing a repository.

## Alternatives considered

- Keep `release_publish` beside the formula for the orchestrator. Rejected: there would be two ways to the remote, and publishes outside the formula would be neither measured nor archived.
- Make `flai release --pending` find and act on the main checkout wherever it runs, with no step field. Rejected: every command that must act on main would need the same special case.
- Put the `push` gate in a `flai guard` rule. Rejected: it would hold Claude Code sessions only, not a harness without the guard hook.
- Make a fourth step a script that writes the versions and tags into the Notes. Rejected: no flai command appends to a story's Notes from a step, and a new script for one line adds more than the body's instruction.

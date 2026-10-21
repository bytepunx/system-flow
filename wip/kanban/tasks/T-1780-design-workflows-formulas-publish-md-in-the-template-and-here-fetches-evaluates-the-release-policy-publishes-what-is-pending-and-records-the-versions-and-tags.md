---
id: T-1780
type: task
nature: improvement
title: design/workflows/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:17Z
updated: 2026-10-21T13:26:33Z
transitions: []
stream: S-0423
tags: [template]
touches: [design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, template/template.yaml, flai/cmd/release.go, flai/cmd/release_evaluate_test.go, docs/users/flai-reference.md]
after: [T-2013]
---
# T-1780 design/workflows/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags

## Work

`flai release --evaluate` exits 0 whether or not the policy is met (`evaluateRelease` in `flai/cmd/release.go`). Add `--require-met`, which exits 2 with the evaluation's reason when the policy is not met, so that a formula step can stop on it. Test it, and regenerate `docs/users/flai-reference.md`.

Then write the `publish` formula at `design/workflows/formulas/publish.md`, nature `improvement`, with four steps, each after the one before:

1. `git fetch`.
2. `flai release --evaluate --require-met`.
3. `flai release --pending`.
4. A step that writes the versions and tags `flai release --pending --json` printed into the instance's Notes.

Describe it in the body, with the baseline marker, in the template (`template/root/design/workflows/formulas/publish.md`) and here, and bump the template's version.

Waits for T-2013, which the formula's main-checkout steps need (TH-0583, question 1).

## Done when

`release_evaluate_test.go` passes with the new flag. `flai check --strict` reports nothing on the formula here. `flai formula show publish` lists the steps, each in a layer of its own.

## Notes

The orchestrator starts the formula only when `release_evaluate` is met. Step 2 therefore stops an instance only when the policy changed between the two, and that stop is a failure with a thread, as any step's is.

Replanned 2026-10-21 by planner-S-0423 on TH-0471: the formula lives under `design/workflows/formulas/` (answer 2), and its instance is worked by a story's agent in its worktree (answer 5). Rewritten to the folder as approved on TH-0583. Whether steps 1 to 3 carry `in: main` and whether step 4 stays, or step 3 becomes `flai release --pending --json` with the body telling the agent to record its versions and tags, awaits the operator's answer to TH-0583, questions 1 and 3.

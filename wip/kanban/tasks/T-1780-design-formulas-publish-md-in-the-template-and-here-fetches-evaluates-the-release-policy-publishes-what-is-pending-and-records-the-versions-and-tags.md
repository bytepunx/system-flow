---
id: T-1780
type: task
nature: improvement
title: design/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:17Z
updated: 2026-10-16T05:48:29Z
transitions: []
stream: S-0423
tags: [template]
touches: [design/formulas/publish.md, template/root/design/formulas/publish.md, template/template.yaml, flai/cmd/release.go, flai/cmd/release_evaluate_test.go, docs/users/flai-reference.md]
---
# T-1780 design/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags

## Work

`flai release --evaluate` exits 0 whether or not the policy is met (`evaluateRelease` in `flai/cmd/release.go`). Add `--require-met`, which exits 2 with the evaluation's reason when the policy is not met, so that a formula step can stop on it. Test it, and regenerate `docs/users/flai-reference.md`.

Then write the `publish` formula, nature `improvement`, with four steps, each after the one before:

1. `git fetch`.
2. `flai release --evaluate --require-met`.
3. `flai release --pending`.
4. A step that writes the versions and tags `flai release --pending --json` printed into the instance's Notes.

Describe it in the body, with the baseline marker, in the template and here, and bump the template's version.

First layer.

## Done when

`release_evaluate_test.go` passes with the new flag. `flai check --strict` reports nothing on the formula here. `flai formula show publish` lists the four steps in four layers.

## Notes

The orchestrator starts the formula only when `release_evaluate` is met. Step 2 therefore stops an instance only when the policy changed between the two, and that stop is a failure with a thread, as any step's is.

---
id: T-2748
type: task
nature: improvement
title: The publish formula's evaluate step stops on exit 2 in both copies, the formula and user docs say how a stop works, and I-0283 is closed
status: backlog
parent: S-0621
owner: alex
created: 2026-11-09T17:07:21Z
updated: 2026-11-09T17:07:21Z
transitions: []
stream: S-0621
tags: [template, docs]
touches: [design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, template/CHANGELOG.md, template/template.yaml, design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md, design/issues/summary.md]
after: [T-2746, T-2747]
---
# T-2748 The publish formula's evaluate step stops on exit 2 in both copies, the formula and user docs say how a stop works, and I-0283 is closed

## Work

Put the stop on the one step that needs it, document it, and close the issue.

- Give the evaluate step of `design/workflows/formulas/publish.md` the stop `[2]`, the exit `flai release --evaluate --require-met` gives when the policy is not met (`exitReleaseNotMet` in `flai/cmd/release.go`). Rewrite its body's "When step 2 stopped the run" line to say the agent runs `flai formula stop`, and leave the failure path for exit 3 and the `push` gate as it is. Copy it to `template/root/design/workflows/formulas/publish.md` above the marker.
- In `design/workflows/formulas/README.md` and its template copy, add the field to the step fields and say what a stop does beside "a step that fails".
- In `docs/users/flai.md`, § A story from a formula and the evidence table, describe the field, `flai formula stop`, and the `skipped` evidence. Regenerate `docs/users/flai-reference.md` with `make flai-reference`. In `docs/operators/index.md` § Writing and running a formula, say a stop exit needs no operator.
- Add a line to `template/CHANGELOG.md` under `## Unreleased`, and bump `template/template.yaml` as the template's other changes do.
- Run `flai formula run publish` in a scratch project, or check with a test, that the instance's evaluate task's `## Done when` names exit 2 as its stop.
- Close the issue: `flai issue close I-0283 --reason "<what fixed it>"`, naming the field, the command, and the acceptance rule, and the tests that reproduce its three instances.

This task waits for T-2746 and T-2747: the docs name the evidence form and the agent's step those settle.

## Done when

- Both copies of the publish formula give the evaluate step `stop: [2]`, and `flai check --strict` passes them.
- The formula README in both copies, the user guide, its reference, and the operators guide describe a stop.
- I-0283 is closed with its reason, and `design/issues/summary.md` no longer lists it.
- `flai test` on the changed paths passes.

## Notes

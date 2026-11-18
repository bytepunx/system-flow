---
id: T-3443
type: task
nature: feature
title: The publish formula's publish step fails on a failed suite's exit, so its instance blocks and asks, naming the suite and tier
status: backlog
parent: S-0811
owner: alex
created: 2026-11-16T23:48:14Z
updated: 2026-11-18T08:16:10Z
transitions: []
stream: S-0811
tags: [flai, template]
touches: [template/root/design/workflows/formulas/publish.md, flai/cmd/publish_formula_suites_test.go, template/CHANGELOG.md]
after: [T-3441]
---
# T-3443 The publish formula's publish step fails on a failed suite's exit, so its instance blocks and asks, naming the suite and tier

## Work

TH-1250: the orchestrator recommended, and S-0811's agent decided, that a suite failing at publish is a failure of the step, not a clean stop ([ADR-0271](../../../design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md) keeps `stop` for exits that are the formula's answer). So T-3441's new exit is **not** added to the `publish` step's `stop`, which stays `[2]`. Under [ADR-0239](../../../design/adrs/0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md) §3 the instance's agent then blocks the step's task, writes the exit and the last 40 lines of its output in the task's Notes, and asks on a thread naming the step: the output names the suite, the tier, and the log, so the story says which suite failed.

Only the template's copy, `template/root/design/workflows/formulas/publish.md`, changes: this repository has no copy since main's 4052fc603, and it is not restored (D-0662). In its body, above the marker, add a bullet: step 3 exits with T-3441's code when a suite `verify_at.publish` lists fails at the main branch's head, before anything is tagged or pushed; that exit fails the step, so block its task and ask, naming the suite, the tier, and the log path the output gives. Add the change to `template/CHANGELOG.md` under `## Unreleased`.

ADR-0316 (proposed) and S-0837 retire the formula later; this keeps the template's copy true until then.

It waits for T-3441, which gives the exit's number and what its output names.

## Done when

- `flai/cmd/publish_formula_suites_test.go` shows, on the template's copy, that the publish step's `stop` does not list T-3441's exit, so that exit fails the step, and that the body names the exit.
- `flai check --strict` passes on the formula file.
- `flai test` passes on the changed paths.

## Notes

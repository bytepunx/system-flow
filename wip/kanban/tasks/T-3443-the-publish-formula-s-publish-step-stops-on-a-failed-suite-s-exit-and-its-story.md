---
id: T-3443
type: task
nature: feature
title: The publish formula's publish step stops on a failed suite's exit, and its story records which suite and tier failed
status: backlog
parent: S-0811
owner: alex
created: 2026-11-16T23:48:14Z
updated: 2026-11-16T23:48:14Z
transitions: []
stream: S-0811
tags: [flai, template]
touches: [design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, flai/cmd/publish_formula_suites_test.go, template/CHANGELOG.md]
after: [T-3441]
---
# T-3443 The publish formula's publish step stops on a failed suite's exit, and its story records which suite and tier failed

## Work

Add T-3441's new exit to the `stop` of the `publish` step in `design/workflows/formulas/publish.md`, beside 2, and make the template's copy, `template/root/design/workflows/formulas/publish.md`, the same. In the body, above the marker, add a bullet: the step exits with that code when a suite `verify_at.publish` lists fails at the main branch's head, before anything is tagged. That exit stops the run cleanly ([ADR-0271](../../adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md)). Run `flai formula stop` with the step and its output. Record in the story's Notes the suite, the tier, and the log path the output names, in place of the versions. Then close the story out and move it to review. Add the change to `template/CHANGELOG.md` under `## Unreleased`.

It waits for T-3441, which gives the exit's number and what its output names.

## Done when

- `flai/cmd/publish_formula_suites_test.go` shows that the publish step of an instance `flai formula run publish` makes lists the new exit in `stop`, and that both copies of the formula are equal.
- `flai check --strict` passes on both formula files.
- `flai test` passes on the changed paths.

## Notes

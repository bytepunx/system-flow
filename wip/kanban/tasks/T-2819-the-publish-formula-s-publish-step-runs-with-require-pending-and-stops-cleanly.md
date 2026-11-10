---
id: T-2819
type: task
nature: improvement
title: The publish formula's publish step runs with --require-pending and stops cleanly on exit 2
status: in-progress
parent: S-0643
owner: agent
created: 2026-11-09T22:12:46Z
updated: 2026-11-10T06:44:29Z
transitions:
  - to: ready
    at: 2026-11-10T06:44:28Z
    by: agent-S-0643
  - to: in-progress
    at: 2026-11-10T06:44:29Z
    by: agent-S-0643
stream: S-0643
tags: [formula, template]
touches: [design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, template/CHANGELOG.md]
after: [T-2818]
---
# T-2819 The publish formula's publish step runs with --require-pending and stops cleanly on exit 2

## Work

In `design/workflows/formulas/publish.md` and its template copy `template/root/design/workflows/formulas/publish.md`, make the `publish` step run `flai release --pending --require-pending --json` with `stop: [2]` (ADR-0271). Then an instance whose evaluate step a publish outran stops cleanly, and does not close as if it published. In the body, say that the step's exit 2 means another publish covered the pending stories. Say that its output names that release, and that the instance's Notes record it in place of the versions it would have published.

Name the publish step's stop in `design/workflows/formulas/README.md` and its template copy, beside the evaluate step's exit 2. Add a line under `## Unreleased` in `template/CHANGELOG.md`.

This task waits for T-2818 because the step names the flag T-2818 adds.

## Done when

- Both copies of `publish.md` name `release --pending --require-pending --json` with `stop: [2]` on the publish step, and their bodies say what that exit means.
- Both copies of the formulas README name the publish step's stop.
- `template/CHANGELOG.md` has the line under `## Unreleased`.
- `flai check --strict` passes on the formula.

## Notes

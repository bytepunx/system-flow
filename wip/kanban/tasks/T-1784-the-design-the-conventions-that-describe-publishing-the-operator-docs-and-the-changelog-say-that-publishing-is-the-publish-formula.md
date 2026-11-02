---
id: T-1784
type: task
nature: improvement
title: The design, the conventions that describe publishing, the operator docs, and the changelog say that publishing is the publish formula
status: done
parent: S-0423
owner: alex
created: 2026-10-16T05:48:51Z
updated: 2026-11-02T04:56:14Z
transitions:
  - to: ready
    at: 2026-11-02T04:46:05Z
    by: agent-S-0423
  - to: in-progress
    at: 2026-11-02T04:46:06Z
    by: agent-S-0423
  - to: done
    at: 2026-11-02T04:56:14Z
    by: agent-S-0423
stream: S-0423
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, docs/operators/index.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
after: [T-1781, T-1782, T-2013, T-2015]
usage:
  source: log
  seconds: 602
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 118
      output: 638
      cache_read: 8237407
      cache_write: 169084
      cost: 2.9583
---
# T-1784 The design, the conventions that describe publishing, the operator docs, and the changelog say that publishing is the publish formula

## Work

- In `design/system/strategic-agents.md` § Publishing a release, say that under `publish` the orchestrator runs the `publish` formula, and keep `release_publish` for a project with none.
- Where `work-management.md` and `strategic-agents.md`, template and here, describe publishing, name the formula. A baseline rule is changed in the template's baseline and in this repository's copy alike, as the template's release would bring it.
- Add publishing through the formula to `docs/operators/index.md`, and name it in `template/CHANGELOG.md`.

Waits for the orchestrator's rule. Runs beside the run, since they share no path.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes

Replanned 2026-10-21 by planner-S-0423: it also waits for T-2013, so that the formula section of `design/system/strategic-agents.md` and `docs/users/flai.md` can say what a step's `in: main` does, and `design/system/flai-cli.md` and `docs/users/flai.md` carry `flai release --evaluate --require-met` and the refusal of `--pending` while `push` is off.

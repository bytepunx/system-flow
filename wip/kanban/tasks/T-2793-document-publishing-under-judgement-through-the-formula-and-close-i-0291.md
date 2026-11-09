---
id: T-2793
type: task
nature: improvement
title: Document publishing under judgement through the formula and close I-0291
status: backlog
parent: S-0625
owner: alex
created: 2026-11-09T22:02:44Z
updated: 2026-11-09T22:02:44Z
transitions: []
stream: S-0625
tags: [docs, template]
touches: [docs/users/flai.md, docs/operators/settings.md, template/CHANGELOG.md, design/issues/I-0291-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md, design/issues/summary.md]
after: [T-2791, T-2792]
---
# T-2793 Document publishing under judgement through the formula and close I-0291

## Work

Say in the guides what changed, and close the issue. It waits for T-2791 and T-2792 because it documents the orchestrator's command and the formula's input they write.

- `docs/users/flai.md`:
  - § Whether a release is due: the `judgement` row and the `--require-met` sentence name `--judged`.
  - § When the orchestrator publishes: under `judgement`, the orchestrator runs the formula with its reasoning as the `judgement` input.
  - The `flai formula run publish` example, if it shows the input.
- `docs/operators/settings.md`: the `orchestration.release.policy` row, which says judgement is never met by itself, names the formula's input.
- `template/CHANGELOG.md`, under `## Unreleased`: the publish formula's `judgement` input, and the orchestrator role file and `strategic-agents.md` changes.
- Close I-0291 with `flai issue close I-0291 --reason "<what fixed it>"`. Name `--judged`, the formula's input, the orchestrator's command, and the test that reproduces the cause. It regenerates `design/issues/summary.md`.

## Done when

- Neither guide says that the `publish` permission cannot publish under `judgement`.
- `template/CHANGELOG.md` lists the change.
- I-0291 is closed with its reason, and `flai check --strict` is clean.

## Notes

Layer 4, the last. T-2790 regenerates the flag index in `docs/operators/settings.md` in layer 2; this task edits only the policy row.

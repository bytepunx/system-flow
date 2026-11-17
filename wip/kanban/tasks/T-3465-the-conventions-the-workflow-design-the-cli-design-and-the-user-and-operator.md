---
id: T-3465
type: task
nature: feature
title: The conventions, the workflow design, the CLI design, and the user and operator guides say what a close-out runs and what the verified rule needs
status: backlog
parent: S-0807
owner: alex
created: 2026-11-16T23:50:06Z
updated: 2026-11-16T23:50:06Z
transitions: []
stream: S-0807
tags: [flai, template]
touches: [design/system/workflow.md, design/system/flai-cli.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/change-the-board.md, template/CHANGELOG.md]
after: [T-3457, T-3460]
---
# T-3465 The conventions, the workflow design, the CLI design, and the user and operator guides say what a close-out runs and what the verified rule needs

## Work

Criterion 5. It waits for T-3457 and T-3460 because it describes what they built.

- `design/system/workflow.md`: add `verified` to the rules table beside `kluster-verified`. Say what the close-out runs: the suites the story's next transitions need, and `--all` every suite. Say what the record keeps per suite.
- `design/system/flai-cli.md` § The close-out job: what it runs, and when it answers from the record.
- `design/conventions/work-management.md` and `design/conventions/delegation.md`: the close-out runs the suites the workflow's next transitions need. Change the baseline text only as this story's design changes it. Make the same change in the template's copies under `template/root/design/conventions/` in the same task, design first, as `CLAUDE.md` asks.
- `docs/users/flai.md` § Verify a story before review, and `docs/operators/runbooks/change-the-board.md`: the rule, its severity, and how a project names suites on a transition.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference` when the help text changed, never by hand.
- Add a line under `## Unreleased` in `template/CHANGELOG.md` for the convention and script changes.

## Done when

- Each document above says what a close-out runs and what the `verified` rule needs.
- The template's convention copies match the design's.
- `flai test` passes on the changed paths, the markdown lint and `TestReferenceIsCurrent` included.

## Notes

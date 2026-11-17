---
id: T-3467
type: task
nature: feature
title: The manifest design, the settings index, the user guide, the reference, and the close-out scripts say queue and when are retired and name the workflow lines in their place
status: backlog
parent: S-0813
owner: alex
created: 2026-11-16T23:50:22Z
updated: 2026-11-16T23:50:22Z
transitions: []
stream: S-0813
tags: [flai, template]
touches: [design/system/project-manifest.md, design/system/workflow.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md]
after: [T-3455, T-3459, T-3461]
---
# T-3467 The manifest design, the settings index, the user guide, the reference, and the close-out scripts say queue and when are retired and name the workflow lines in their place

## Work

The documentation rule: behaviour that changes changes its documents in the same story. It waits for T-3455, T-3459, and T-3461, whose behaviour it describes.

- `design/system/project-manifest.md`: mark `queue` and `when` retired, name the ADR T-3450 wrote, and give the workflow lines that stand in their place.
- `design/system/workflow.md` and `design/system/flai-cli.md`: where they describe a `queue` tier or `when`, say the workflow file's `verify:` and `verify_at.batch` decide it. Name the `manifest.tests-queue` advice and what `flai upgrade` writes into a file it adds.
- `docs/operators/settings.md`: rewrite the `tests[].queue` and `tests[].when` rows as retired, read only while the file names no `batch` suites.
- `docs/users/flai.md`: rewrite its queue tiers section. Regenerate `docs/users/flai-reference.md` with `make flai-reference` if a command's help changed.
- `scripts/close-out.sh` and `template/root/scripts/close-out.sh`: rewrite the comment that names `queue: true` and keep the two in step. Add a `template/CHANGELOG.md` line under `## Unreleased` for the template's change.

## Done when

- None of these documents describes `queue` or `when` as deciding a step without naming the workflow file.
- `flai test` passes for the changed paths, the `settings-index` tier and the markdown lint among them, and ShellCheck for the two scripts.

## Notes

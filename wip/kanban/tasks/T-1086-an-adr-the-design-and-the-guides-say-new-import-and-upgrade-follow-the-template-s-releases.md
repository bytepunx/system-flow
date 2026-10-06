---
id: T-1086
type: task
nature: remediation
title: An ADR, the design, and the guides say new, import, and upgrade follow the template's releases
status: cancelled
parent: S-0301
owner: alex
created: 2026-10-06T22:52:52Z
updated: 2026-10-06T22:53:27Z
transitions:
  - to: cancelled
    at: 2026-10-06T22:53:27Z
    by: agent-S-0301
stream: S-0301
tags: [cli]
touches: [design/adrs, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/conventions.md, docs/contributors/template.md, docs/operators/settings.md, docs/operators/runbooks, flai/cmd/upgrade.go]
after: [T-1081, T-1084]
---
# T-1086 An ADR, the design, and the guides say new, import, and upgrade follow the template's releases

## Work

The story's documentation, after the behaviour is in.

- A new ADR in `design/adrs`, numbered with `flai adr new`, records the decision:
  - `flai new`, `flai import`, and `flai upgrade` follow the template's releases. A ref that is empty, the default branch, or a version tag resolves to the newest `vX.Y.Z` tag.
  - `--ref` wins and is recorded.
  - A version or ref the operator sets in `system-flow.yaml` that differs from the lock is their intent. When it differs from the latest, they are asked, and flai refuses without a terminal.
  - The lock is history and never a target.
  - A branch clone is fetched again.
  - It refines ADR-0015 without superseding it.
- Update `design/system/project-manifest.md`, and `design/system/flai-cli.md` where it describes upgrade or new.
- Update the user guide where `flai upgrade` is described for users: `docs/users/flai.md`, or `docs/users/conventions.md` and `docs/contributors/template.md`, whichever describes it. Also update `docs/operators/runbooks/migrate.md` and the update runbook if they describe it.
- Update the `template.ref` row in `docs/users/flai.md` and in `docs/operators/settings.md`.
- Change the `upgrade` command's `Long` text and `--ref` flag help to match, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

It waits for T-1081 and T-1084, whose behaviour it describes.

## Done when

- [ ] The ADR exists and is linked from the design.
- [ ] Every document that describes how upgrade, new, or import pick a template version says what the code now does.
- [ ] The reference is regenerated.
- [ ] `flai check --strict` and the markdown lint are clean for these files.

## Notes
- 2026-10-06T22:53:27Z: moved to cancelled: duplicate of the planner's draft for S-0301, written before I saw it; the planner's tasks are kept and edited

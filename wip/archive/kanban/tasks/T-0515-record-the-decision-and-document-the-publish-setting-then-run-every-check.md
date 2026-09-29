---
id: T-0515
type: task
nature: remediation
title: Record the decision and document the publish setting, then run every check
status: done
parent: S-0144
owner: alex
created: 2026-09-29T01:14:20Z
updated: 2026-09-29T01:20:09Z
transitions:
  - to: ready
    at: 2026-09-29T01:14:25Z
    by: agent-S-0144
  - to: in-progress
    at: 2026-09-29T01:16:48Z
    by: agent-S-0144
  - to: done
    at: 2026-09-29T01:20:09Z
    by: agent-S-0144
stream: S-0144
tags: []
touches: [design/adrs, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/conventions/git.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md]
---
# T-0515 Record the decision and document the publish setting, then run every check

## Work

- An ADR refining ADR-0032: pushing never releases unless the `publish` host action is on; S-0094's release-at-push becomes that opt-in.
- Update `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, the project addition in `design/conventions/git.md`, `docs/operators/index.md`, `docs/operators/settings.md`, `docs/users/flai.md`, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Run `make flai-test` (all tiers and lint) and `flai check --strict`.

## Done when

- Every doc that said pushing releases says it does so only with `publish` enabled, the ADR is in the index, and all checks pass.

## Notes

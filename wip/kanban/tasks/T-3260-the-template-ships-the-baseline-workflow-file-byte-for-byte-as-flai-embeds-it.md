---
id: T-3260
type: task
nature: feature
title: The template ships the baseline workflow file byte for byte as flai embeds it, and flai upgrade keeps a project's changes to it
status: done
parent: S-0755
owner: alex
created: 2026-11-12T16:10:15Z
updated: 2026-11-15T22:28:40Z
transitions:
  - to: ready
    at: 2026-11-15T22:14:34Z
    by: agent-S-0755
  - to: in-progress
    at: 2026-11-15T22:14:34Z
    by: agent-S-0755
  - to: done
    at: 2026-11-15T22:28:39Z
    by: agent-S-0755
stream: S-0755
tags: [cli, template]
touches: [template/root/design/workflows/board.md, flai/internal/workflow/baseline_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/CHANGELOG.md]
after: [T-3253]
usage:
  source: log
  seconds: 845
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 10107
      cache_read: 2084260
      cache_write: 61156
      cost: 0.979
---
# T-3260 The template ships the baseline workflow file byte for byte as flai embeds it, and flai upgrade keeps a project's changes to it

## Work

Waits for T-3253: a project made from the template must have its file read by the load that task adds. ADR-0307 §4.

- Add the baseline file to the template at `template/root/design/workflows/board.md`, byte for byte `flai/internal/workflow/baseline/board.md`, and extend `baseline_test.go` to hold the two equal, as `roles/baseline_test.go` holds the role files.
- `flai upgrade` (`Compute` in `upgrade.go`): adds the file to a project that lacks it; never writes it once the project has one, whatever its content and the lock; when the template's file differs from the hash the lock records for it, says the template's baseline changed and that the operator carries it over in a story, and the lock records the new hash.
- Record the change under `## Unreleased` in `template/CHANGELOG.md` and bump `template/template.yaml` as earlier template stories did.

## Done when

`baseline_test.go` and `upgrade_test.go` pass (absent: added; present and edited: kept; present and the template changed: kept and named once), and `scripts/template-test.sh` renders a project whose check is clean.

## Notes

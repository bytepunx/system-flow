---
id: T-1057
type: task
nature: remediation
title: flai upgrade --ref wins over the manifest's template.ref and records the ref and version it applied
status: done
parent: S-0301
owner: alex
created: 2026-10-06T22:49:54Z
updated: 2026-10-06T23:00:09Z
transitions:
  - to: ready
    at: 2026-10-06T22:54:33Z
    by: agent-S-0301
  - to: in-progress
    at: 2026-10-06T22:54:34Z
    by: agent-S-0301
  - to: done
    at: 2026-10-06T23:00:09Z
    by: agent-S-0301
stream: S-0301
tags: [cli, template]
touches: [flai/cmd/upgrade.go, flai/cmd/upgrade_test.go]
usage:
  source: log
  seconds: 335
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 164
      cache_read: 1229238
      cache_write: 74283
      cost: 0.5741
---
# T-1057 flai upgrade --ref wins over the manifest's template.ref and records the ref and version it applied

## Work

Criterion 1. In `flai/cmd/upgrade.go`, the `RunE` of `flai upgrade` sets `tplRepo, ref = mf.Template.Repo, mf.Template.Ref` whenever `--template` is not given. That throws away a `--ref` given on its own, so `flai upgrade --ref v1.0.60` renders whatever the manifest's ref names, for example `main` or 1.0.18.

- Fall back to the manifest's ref only when `--ref` is empty. Take the manifest's repo as now.
- When `--ref` was given, `writeTemplateFields` writes it to `template.ref` in `system-flow.yaml`. Today only a changed `--template` does that (`sourceChanged`). The lock saved after a clean upgrade records the same ref and the template's version under `template`.
- A `--ref` that names the project's current version is still "already at template"; `--force` re-applies, as now.
- This task waits for nothing. It is the narrow fix and lands first, so the defect is gone even before the prompt task builds on these lines.

## Done when

- [ ] A test in `flai/cmd/upgrade_test.go` upgrades a project whose manifest has `template.ref: main` and version 9.9.9, with `--ref` naming a tag at 10.0.0. The render is 10.0.0, `system-flow.yaml` has `template.ref` set to that tag and `template.version: 10.0.0`, and `system-flow.lock.yaml` records the same ref and version.
- [ ] Without `--ref`, the existing upgrade tests pass unchanged.
- [ ] `scripts/flai-test.sh` passes for `flai/cmd`.

## Notes

Drafted by the planner. The cause is at `flai/cmd/upgrade.go`, lines 54 to 56, on 2026-10-06.

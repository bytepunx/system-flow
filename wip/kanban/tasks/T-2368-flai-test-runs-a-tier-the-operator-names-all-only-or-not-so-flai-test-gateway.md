---
id: T-2368
type: task
nature: feature
title: flai test runs a tier the operator names, all_only or not, so flai test gateway runs the gateway checks on demand
status: done
parent: S-0548
owner: alex
created: 2026-10-31T07:25:31Z
updated: 2026-10-31T09:26:02Z
transitions:
  - to: ready
    at: 2026-10-31T09:08:24Z
    by: agent-S-0548
  - to: in-progress
    at: 2026-10-31T09:08:25Z
    by: agent-S-0548
  - to: done
    at: 2026-10-31T09:26:02Z
    by: agent-S-0548
stream: S-0548
tags: [cli]
touches: [flai/cmd/test.go, flai/cmd/test_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, design/system/project-manifest.md, docs/operators/settings.md, flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, template/CHANGELOG.md]
after: [T-2363, T-2366]
usage:
  source: log
  seconds: 1057
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 150
      output: 49263
      cache_read: 12967402
      cache_write: 261936
      cost: 5.0648
---
# T-2368 flai test runs a tier the operator names, all_only or not, so flai test gateway runs the gateway checks on demand

## Work

`flai test` takes paths and packages and `--all`; it has no way to run one tier by its name, which story criterion 2 needs as `flai test gateway`. Add it: an argument that names a tier in the manifest's `tests` runs that tier alone, `all_only` or not, for the whole checkout, with its `all_command` when it has one; a name that is both a tier and a path is read as the tier, and a tier the manifest lacks is refused with exit 2 naming the tiers. Choose between a bare name and a `--tier <name>` flag, keep `flai test gateway` working as the criterion writes it, and record the choice in the narrative's `## Decisions`.

Document it in `docs/users/flai.md` (Run the tests for what changed), regenerate `docs/users/flai-reference.md` with `make flai-reference`, and update the `flai test` section of `design/system/flai-cli.md`.

Second layer: it waits for T-2363, which documents its check warning in `docs/users/flai.md` and `design/system/flai-cli.md` too. The gateway tier task waits for it, since that tier is reached only by name.

## Done when

`flai test <tier>` runs the named tier alone and answers as `flai test` does; `test_test.go` covers a named tier, an `all_only` tier by name, and an unknown name refused with exit 2. The docs say how.

## Notes

Drafted by the planner of S-0548 on 2026-10-31.

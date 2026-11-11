---
id: T-3090
type: task
nature: improvement
title: The close-out's Go tests run without -race, one package at a time, on one core, and a manual race tier keeps the detector
status: done
parent: S-0722
owner: alex
created: 2026-11-11T06:06:39Z
updated: 2026-11-11T06:26:31Z
transitions:
  - to: ready
    at: 2026-11-11T06:07:07Z
    by: agent-S-0722
  - to: in-progress
    at: 2026-11-11T06:07:08Z
    by: agent-S-0722
  - to: done
    at: 2026-11-11T06:26:31Z
    by: agent-S-0722
stream: S-0722
tags: []
touches: [system-flow.yaml, scripts/integration.sh, scripts/test.sh, scripts/flai-test.sh, scripts/README.md, flai/internal/verify/select_test.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/manifest_tests.go]
usage:
  source: log
  seconds: 1163
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 49
      output: 14964
      cache_read: 4243670
      cache_write: 86893
      cost: 1.7662
---
# T-3090 The close-out's Go tests run without -race, one package at a time, on one core, and a manual race tier keeps the detector

## Work

- `system-flow.yaml`'s `go-test` tier and `scripts/integration.sh` run `go test` without `-race`, with `-p 1` and the test binaries at one core; `scripts/test.sh` matches.
- Add the tier `race`, `manual: true` like `gateway`, which runs the full suite under `-race` on every core as `integration.sh` does today; `scripts/flai-test.sh` runs it in place of the plain integration run.
- A repository-manifest test in `flai/internal/verify/select_test.go` holds both: the close-out's Go tiers carry no `-race`, and `race` is manual.
- `scripts/README.md` says what each script now runs.
- Waits for nothing: it is the first layer.

## Done when

- `flai test` on the changed paths passes, and `flai check --strict` finds no `manifest.tests` or `manifest.tests-timeout`.

## Notes

---
id: T-3577
type: task
nature: improvement
title: flai manifest set's help names the image keys, and a command test reproduces I-0364
status: done
parent: S-0839
owner: alex
created: 2026-11-17T12:29:51Z
updated: 2026-11-18T09:48:28Z
transitions:
  - to: ready
    at: 2026-11-18T09:43:06Z
    by: agent-S-0839
  - to: in-progress
    at: 2026-11-18T09:43:12Z
    by: agent-S-0839
  - to: done
    at: 2026-11-18T09:48:28Z
    by: agent-S-0839
stream: S-0839
tags: [flai]
touches: [flai/cmd/manifest.go, flai/cmd/manifest_test.go, docs/users/flai-reference.md]
after: [T-3574]
usage:
  source: log
  seconds: 316
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 7943
      cache_read: 1458026
      cache_write: 75820
      cost: 0.9324
---
# T-3577 flai manifest set's help names the image keys, and a command test reproduces I-0364

## Work

In `flai/cmd/manifest.go`, make `flai manifest set`'s help name the two keys T-3574 adds. Add an example such as `flai manifest set agent.image.args.GO_VERSION=1.26.0`. Today the help says only the strategic agents' settings are written.

In `flai/cmd/manifest_test.go`, add the test that reproduces I-0364. It starts from a manifest with a layer and no args, as S-0760 found this repository's. Running `flai manifest set agent.image.args.GO_VERSION=1.26.0 agent.image.args.PNPM_VERSION=10.34.5` writes both under `agent.image.args`, and the manifest then loads with no `manifest.agent` problem. `--unset agent.image.args.GO_VERSION` removes the one arg. `agent.image.args.BASE=x` exits 1 with the field and the reason, and writes nothing. With `--json`, the answer lists the keys set.

Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

This task waits for T-3574: it calls the settings T-3574 adds, through the command.

## Done when

- The new test fails at the commit before T-3574 and passes after it.
- `flai test flai/cmd/manifest.go flai/cmd/manifest_test.go docs/users/flai-reference.md` passes, `TestReferenceIsCurrent` among its tests.

## Notes

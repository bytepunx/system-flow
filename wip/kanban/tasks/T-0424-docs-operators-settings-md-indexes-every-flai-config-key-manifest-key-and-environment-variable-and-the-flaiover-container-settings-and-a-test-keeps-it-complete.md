---
id: T-0424
type: task
nature: feature
title: docs/operators/settings.md indexes every flai config key, manifest key, and environment variable, and the flaiover container settings, and a test keeps it complete
status: done
parent: S-0028
owner: alex
created: 2026-09-24T09:28:01Z
updated: 2026-09-24T09:32:11Z
transitions:
  - to: ready
    at: 2026-09-24T09:28:09Z
    by: agent-S-0028
  - to: in-progress
    at: 2026-09-24T09:28:10Z
    by: agent-S-0028
  - to: done
    at: 2026-09-24T09:32:11Z
    by: agent-S-0028
stream: S-0028
tags: []
touches: [docs/operators, flai/cmd]
---
# T-0424 docs/operators/settings.md indexes every flai config key, manifest key, and environment variable, and the flaiover container settings, and a test keeps it complete

## Work
- Write `docs/operators/settings.md`: one table per kind (flai configuration `~/.flai/config.json`, the manifest `system-flow.yaml`, environment variables flai reads and sets, flaiover container settings), each row naming the setting, where it is set, its default, and a link to where it is explained.
- Add a behaviour test in `flai/cmd` that fails when a key of the config or manifest schema, or an environment variable name flai's source uses, is missing from the page.

## Done when
- Every key of `config.Config` and `manifest.Manifest`, every `FLAI_*`/`FLAIOVER_*` variable flai uses outside tests, and every variable flaiover reads is on the page.
- The test fails with a missing name and passes on the page; `scripts/flai-test.sh` passes.

## Notes

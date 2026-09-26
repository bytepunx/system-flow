---
id: T-0437
type: task
nature: feature
title: flai import registers the project with flai serve when a host runs, and names the command that serves it when none does
status: done
parent: S-0120
owner: alex
created: 2026-09-26T05:25:43Z
updated: 2026-09-26T05:27:39Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:13Z
    by: agent-S-0117
  - to: in-progress
    at: 2026-09-26T05:26:14Z
    by: agent-S-0117
  - to: done
    at: 2026-09-26T05:27:39Z
    by: agent-S-0117
stream: S-0120
tags: []
---

# T-0437 flai import registers the project with flai serve when a host runs, and names the command that serves it when none does

## Work

After writing the manifest, `flai import` (with or without `--commit`) registers the project with `flai serve` using the entry `flai dashboard` writes (key, name, main root, the dashboard's dial address from `dashboardSettings`, the shared agent credential) when a `flai host` runs for the config, and says it is served at the dashboard's address. When none runs it says the project is not served yet and names `flai dashboard` in the project as the one command that serves it. `--json` carries the same as a `serve` object. Tests in `cmd/import_serve_test.go` with a fake host state.

## Done when

Both paths are covered by tests that fail without the change; `make test` passes.

## Notes

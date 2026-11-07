---
id: T-2600
type: task
nature: feature
title: A deploy builds each component as its own step and records every step as pending before the first starts
status: done
parent: S-0595
owner: alex
created: 2026-11-07T22:07:00Z
updated: 2026-11-07T22:32:09Z
transitions:
  - to: ready
    at: 2026-11-07T22:09:59Z
    by: agent-S-0595
  - to: in-progress
    at: 2026-11-07T22:24:48Z
    by: agent-S-0595
  - to: done
    at: 2026-11-07T22:32:09Z
    by: agent-S-0595
stream: S-0595
tags: [cli]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
after: [T-2598]
usage:
  source: log
  seconds: 441
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 11211
      cache_read: 5479940
      cache_write: 43238
      cost: 1.6022
---
# T-2600 A deploy builds each component as its own step and records every step as pending before the first starts

## Work

ADR-0263 §2 and §3, in `flai/internal/kluster/deploy.go`. With `Kluster.Components`, `steps` runs one build step per component in order, named `build-<name>`, each with its own log (`LogPath`), its own images file `logs/<namespace>/images-<name>` as `FLAI_KLUSTER_IMAGES`, and `FLAI_KLUSTER_COMPONENT=<name>`; the images of all components are concatenated into `r.Images` for one load; a failed build stops the deploy and the steps after it are recorded `skipped`. `FailedStep` names `build-<name>`. With `build` alone, nothing changes but the plan. Add the outcome `pending`: a deploy, a promotion, and a teardown write every step they will run into the record as `pending` before the first starts; `step` marks its step `running` with `started_at` and `log`, writes the record, then sets the outcome and `ended_at`, replacing the planned entry rather than appending. Waits for T-2598, which adds `Kluster.Components`.

## Done when

- Tests in `deploy_test.go` (with the fake programs the file already uses) show three components giving `build-a`, `build-b`, `build-c`, `load`, `deploy` steps with their own logs and images; a failing second component leaving `build-c`, `load`, `deploy` skipped and `failed_step` `build-b`; the record holding all steps `pending` before the first runs and the running step `running` while it runs; the single-`build` deploy and teardown unchanged but for the plan. They fail without the change.
- `flai test` on the changed paths passes.

## Notes

---
id: T-2171
type: task
nature: feature
title: flai check reports agent.image when the manifest names a layer that is not built or whose Dockerfile is missing, and notes when docker is not on PATH
status: done
parent: S-0510
owner: alex
created: 2026-10-26T04:50:56Z
updated: 2026-11-02T18:46:18Z
transitions:
  - to: ready
    at: 2026-11-02T18:33:45Z
    by: agent-S-0510
  - to: in-progress
    at: 2026-11-02T18:33:46Z
    by: agent-S-0510
  - to: done
    at: 2026-11-02T18:46:18Z
    by: agent-S-0510
stream: S-0510
tags: [cli, docker]
touches: [flai/internal/check/check.go, flai/internal/check/image.go, flai/internal/check/image_test.go, flai/internal/check/story.go, flai/cmd/check.go, flai/cmd/check_stats_test.go, flai/internal/verify/story.go]
after: [T-2164]
usage:
  source: log
  seconds: 752
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 72
      output: 364
      cache_read: 2825431
      cache_write: 89796
      cost: 1.0322
---
# T-2171 flai check reports agent.image when the manifest names a layer that is not built or whose Dockerfile is missing, and notes when docker is not on PATH

## Work

Add a check `agent.image` in `flai/internal/check/image.go`, registered in `check.go`: when the manifest names `agent.image.layer`, it is a warning when the Dockerfile does not exist, and a warning naming `make agent-images` when `docker image inspect flai-agent-<project key>-<inner>:local` finds no image for the default inner harness; when `docker` is not on `PATH` it is a note saying the image could not be checked, as `components.examples` notes a missing `node`; when the manifest names no layer it is silent. Run `docker` through the runner the other checks use, so a test fakes it. Test each outcome in `image_test.go`. Waits for T-2164, which reads the key.

## Done when

`image_test.go` covers the four outcomes, `flai check --strict` on this repository with the layer built is silent on `agent.image`, and `flai test flai/internal/check` passes.

## Notes

A check scoped to a story leaves out an `agent.image` outside it, as the other repository-wide checks do; follow the pattern in `check.go`.

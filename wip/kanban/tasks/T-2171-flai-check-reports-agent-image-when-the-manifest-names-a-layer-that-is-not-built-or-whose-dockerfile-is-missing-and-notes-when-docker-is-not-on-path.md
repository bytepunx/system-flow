---
id: T-2171
type: task
nature: feature
title: flai check reports agent.image when the manifest names a layer that is not built or whose Dockerfile is missing, and notes when docker is not on PATH
status: backlog
parent: S-0510
owner: alex
created: 2026-10-26T04:50:56Z
updated: 2026-10-26T04:50:56Z
transitions: []
stream: S-0510
tags: [cli, docker]
touches: [flai/internal/check/check.go, flai/internal/check/image.go, flai/internal/check/image_test.go]
after: [T-2164]
---
# T-2171 flai check reports agent.image when the manifest names a layer that is not built or whose Dockerfile is missing, and notes when docker is not on PATH

## Work

Add a check `agent.image` in `flai/internal/check/image.go`, registered in `check.go`: when the manifest names `agent.image.layer`, it is a warning when the Dockerfile does not exist, and a warning naming `make agent-images` when `docker image inspect flai-agent-<project key>-<inner>:local` finds no image for the default inner harness; when `docker` is not on `PATH` it is a note saying the image could not be checked, as `components.examples` notes a missing `node`; when the manifest names no layer it is silent. Run `docker` through the runner the other checks use, so a test fakes it. Test each outcome in `image_test.go`. Waits for T-2164, which reads the key.

## Done when

`image_test.go` covers the four outcomes, `flai check --strict` on this repository with the layer built is silent on `agent.image`, and `flai test flai/internal/check` passes.

## Notes

A check scoped to a story leaves out an `agent.image` outside it, as the other repository-wide checks do; follow the pattern in `check.go`.

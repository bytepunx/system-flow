---
id: T-2157
type: task
nature: feature
title: The protected list gains flai/images, scripts/agent-image.sh, and scripts/agent-container.sh
status: done
parent: S-0508
owner: alex
created: 2026-10-26T04:48:17Z
updated: 2026-11-02T07:09:52Z
transitions:
  - to: ready
    at: 2026-11-02T06:53:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T06:53:45Z
    by: system-flow
  - to: done
    at: 2026-11-02T07:09:52Z
    by: agent-S-0508
stream: S-0508
tags: [cli, docker]
touches: [flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/guard/call.go, flai/internal/guard/call_test.go]
usage:
  source: log
  seconds: 967
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 10381
      cache_read: 2370296
      cache_write: 78774
      cost: 1.2119
---
# T-2157 The protected list gains flai/images, scripts/agent-image.sh, and scripts/agent-container.sh

## Work

Add to the protected list in `flai/internal/protected/protected.go`, under a `container` harness entry beside `claude-code`, the paths that change what a containerised agent may do: the folder `flai/images/` (every Dockerfile under it), `scripts/agent-image.sh`, and `scripts/agent-container.sh`, so that a story whose branch changes one is accepted by the operator alone (ADR-0106) and the guard's hold-and-ask covers a write to them. Test in `protected_test.go` that `Path` matches each and does not match a sibling. This task waits for nothing.

## Done when

`protected_test.go` covers the three paths, and `flai test flai/internal/protected` passes.

## Notes

ADR-0203 names the manifest layer, the start request, and the harness's configuration as protected; the start request is on standard input and the configuration lives on the host, so the files here are what remains in the repository.

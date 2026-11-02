---
id: T-2172
type: task
nature: feature
title: The container adapter runs the project's layer image in place of the harness image when the manifest names one
status: done
parent: S-0510
owner: alex
created: 2026-10-26T04:51:02Z
updated: 2026-11-02T18:46:40Z
transitions:
  - to: ready
    at: 2026-11-02T18:33:48Z
    by: agent-S-0510
  - to: in-progress
    at: 2026-11-02T18:33:49Z
    by: agent-S-0510
  - to: done
    at: 2026-11-02T18:46:40Z
    by: agent-S-0510
stream: S-0510
tags: [cli, docker]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2164]
usage:
  source: log
  seconds: 771
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 6482
      cache_read: 2846025
      cache_write: 34510
      cost: 0.9319
---
# T-2172 The container adapter runs the project's layer image in place of the harness image when the manifest names one

## Work

Give the container adapter the project's layer: `flai serve` reads `agent.image.layer` from the project's manifest and, when set, passes the derived tag `flai-agent-<project key>-<inner>:local` to the adapter in place of the host's image for the inner harness (`Request` or `Host.Args`, whichever T-2155 chose for the image), and `startAllowed`'s image refusal checks that tag. Test in `container_test.go` that the argv names the layer's tag when given and the harness image otherwise, and in `agents_test.go` that a project with the key set starts on the layer. Waits for T-2164 for the manifest key.

## Done when

A story of a project whose manifest names a layer starts on `flai-agent-<project key>-<inner>:local` in the tests, one without starts on the harness image, and `flai test flai/internal/harness flai/internal/serve` passes.

## Notes

The host's `--image` setting is the harness image; the layer is the project's, so the project wins when both exist.

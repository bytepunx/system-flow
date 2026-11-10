---
id: S-0672
type: story
nature: remediation
title: The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there
status: ready
owner: alex
created: 2026-11-10T08:34:49Z
updated: 2026-11-10T09:29:05Z
transitions:
  - to: ready
    at: 2026-11-10T08:50:47Z
    by: alex
tags: [docker, scripts]
topics: [harness, docker, testing]
touches: [flai/images/system-flow/Dockerfile, scripts/agent-image.sh, scripts/README.md, design/tech/docker.md, docs/operators/index.md, docs/users/flai.md, design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 13
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 371213
          cache_write: 5094
          cost: 0.102
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-10T08:34:49Z
  value: 75
  by: planner-S-0672
  at: 2026-11-10T08:47:48Z
forecast:
  duration: 1h
  delivery: 2026-11-10T11:30:00Z
  basis: "Its own forecast of 1h; 4th in the pull order with an in-progress limit of 5, behind S-0670, S-0671, S-0682, S-0684, S-0687, S-0668, S-0669 and S-0673."
  by: flai
  at: 2026-11-10T09:29:05Z
finalized:
  by: alex
  at: 2026-11-10T08:50:41Z
---
# S-0672 The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there

## Goal

This story remediates [I-0281](../../../design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md), "The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0281 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0281 is closed with `flai issue close I-0281 --reason` saying what fixed it

## Tasks
- T-2917 The system-flow layer carries curl, openssl, and GoReleaser at the version scripts/install-tools.sh pins
- T-2919 scripts/agent-image.sh --check runs the smoke tier inside the built layer and fails on an image that lacks a tool it runs
- T-2922 The docs that list the system-flow layer's tools name curl, openssl, GoReleaser, and scripts/agent-image.sh --check
- T-2924 I-0281 is closed with the reason that the layer carries curl, openssl, and GoReleaser and scripts/agent-image.sh --check passes

## Notes

Cost of delay inputs set by flai from I-0281. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T07:12:34Z, 3.1 days before this story; under one cycle counts as one).

### Planning

The fix, from I-0281's one instance (S-0565's T-2452): the layer `flai/images/system-flow/Dockerfile` adds `curl`, `openssl`, and GoReleaser at the version `scripts/install-tools.sh` pins (v2.18.1). `install.sh` stops on "curl is required" first, then on "openssl is required" (it verifies the release's signature, and `scripts/install-test.sh` signs its test releases with `openssl`); Alpine's base has neither binary. `scripts/flai-snapshot.sh` then needs `goreleaser`. GoReleaser v2.18.1 builds only with Go 1.27, while the layer holds `flai/go.mod`'s Go 1.26.0 with `GOTOOLCHAIN=local`, so the layer takes GoReleaser's release archive for the target architecture, checked against its `checksums.txt`, or builds it in a stage of its own on a Go 1.27 image. The test that reproduces the cause is a check of the built layer, `scripts/agent-image.sh --check`, which runs the smoke tier inside the image with the tree mounted, as the close-out's container does (ADR-0274): it fails on the image without the tools and passes on the new one.

Touches, each a file, none a folder:

| Touch | From |
|-------|------|
| `flai/images/system-flow/Dockerfile` | design: I-0281's remedy names the layer |
| `scripts/agent-image.sh` | layout: it passes `install-tools.sh`'s pins as build arguments, so it passes `GORELEASER_VERSION`, and it gains `--check` |
| `scripts/README.md` | layout: its `agent-image.sh` row describes the flags and the pins |
| `design/tech/docker.md` | design: its agent images row lists the layer's tools and versions |
| `docs/operators/index.md` | design: its container harness entry lists what this repository's layer carries |
| `docs/users/flai.md` | design: its `agent.image.layer` paragraph lists the same |
| `design/issues/I-0281-…md` | declared by criterion 2, closed with `flai issue close` |
| `design/issues/summary.md` | co-change: written whenever an issue closes |

`flai touches suggest` lists co-change with these paths led by `design/system/flai-cli.md` (28%), `docs/users/flai-reference.md`, and the dashboard's documents. They co-change through `design/issues/summary.md` and the docs, not through the layer, and none of them describes the layer's tools, so none is taken.

No `after`: the story reads no output another open story writes.

Forecast: `flai forecast` gives 22m, 131 s per unit of size over 50 done remediation stories in the medium band, times size 10. Raised to 1h: the layer's last stage is rebuilt, and the smoke tier runs inside the container at least twice, on the image without the tools and on the new one, each building flai and two GoReleaser snapshots; S-0630, which changed the same Dockerfile, took 30m of agent time. Delivery moved from flai's 12:07Z by the same 38m, to 12:45Z.

Cost of delay: `flai cod` gives 75.00 USD a week, from the 30m per 168h cycle that flai set from I-0281 at 150 USD an hour. It stands: the one instance cost 30m, and a story's close-out passes on the host meanwhile, so nothing argues for more.

The paths under `flai/images/` and `scripts/agent-image.sh` are on the protected list (T-2157), so the story's agent writes them itself, not a sub-agent, and the operator accepts the story (ADR-0106). Building the layer and running its check need Docker; an agent that runs inside the container has none, and then asks the operator on a thread to run `scripts/agent-image.sh --layer` and `--check` on the host.

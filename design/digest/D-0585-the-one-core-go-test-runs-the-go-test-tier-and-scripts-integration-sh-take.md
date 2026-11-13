---
id: D-0585
title: "The one-core Go test runs, the go-test tier and scripts/integration.sh, take -timeout 10m; only the race tier keeps the 20m sized for -race"
kind: decision
sources: [S-0722, I-0227]
topics: [go, testing, flai, code]
paths: [system-flow.yaml, scripts/integration.sh, scripts/test.sh, flai/cmd, flai/internal/mcpserver, flai/internal/guard]
date: 2026-11-11
status: absorbed
absorbed_by: design/system/project-manifest.md
asked_again: []
cited: []
---

# D-0585 The one-core Go test runs, the go-test tier and scripts/integration.sh, take -timeout 10m; only the race tier keeps the 20m sized for -race

S-0548 gave the go-test tier `-timeout 20m` for I-0227, when it ran under `-race` on every core and `flai/internal/mcpserver` alone could pass 10 minutes on a loaded host (D-0416). S-0722 moved the close-out's Go tests to one core without `-race` (D-0584) and set the timeouts again from three one-core runs. From S-0722's Decisions (`wip/archive/agents/S-0722.md`):

```text
Timeouts: the one-core runs, the go-test tier and integration.sh, take -timeout 10m, not the 20m sized for -race: the slowest packages took 130 to 133 s (./cmd) and 121 s (./internal/mcpserver) in all three runs, so 10m holds them on a host loaded three times over and names a hung test sooner. integration.sh --race, the race tier, keeps 20m, its comment restating why with this story's figure: the slowest package under -race took 252 s (./internal/guard) on every core at a load of 69.
```

So a timeout raised for a one-core tier needs a package past about 10 minutes at one core, not the `-race` figures. When S-0722 was accepted, `docs/users/flai.md` § on `flai test`'s tiers still said this repository's Go tiers pass `-timeout 20m`; `system-flow.yaml` and `scripts/integration.sh` are what they run.

Sources: S-0722; I-0227; D-0416.

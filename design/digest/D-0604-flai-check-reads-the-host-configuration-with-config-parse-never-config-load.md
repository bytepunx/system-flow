---
id: D-0604
title: "flai check reads the host configuration with config.Parse, never config.Load, which writes a missing config file, from the file FLAI_CONFIG names, else ~/.flai/config.json, never a --config flag"
kind: decision
sources: [S-0735]
topics: [cli, go, docker, code]
paths: [flai/internal/check/image.go, flai/internal/config/config.go]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0604 flai check reads the host configuration with config.Parse, never config.Load, which writes a missing config file, from the file FLAI_CONFIG names, else ~/.flai/config.json, never a --config flag

S-0735 made `flai check`'s `agent.image` finding judge the same images `flai serve` starts, so the check needs the host's `--image` for each project, which lives in the host configuration. From S-0735's Decisions (`wip/archive/agents/S-0735.md`):

```text
`flai check` reads the host's `--image` per project from the config file `FLAI_CONFIG` names, else `~/.flai/config.json`, with `config.Parse`. `config.Load` would write a missing file. So the check and flai serve judge the same images.
```

The verifier's review added that `flai check` reads `FLAI_CONFIG` or `~/.flai/config.json`, never a `--config` flag. So a later check or read-only command that needs a host setting parses the file with `config.Parse`. `config.Load` would create the file on a host or in a container that has none.

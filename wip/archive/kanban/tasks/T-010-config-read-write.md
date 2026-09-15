---
id: T-010
type: task
nature: feature
title: Config read and write
status: done
parent: S-004
owner: agent
created: 2026-09-15T17:05:00Z
updated: 2026-09-15T17:38:00Z
transitions:
  - to: ready
    at: 2026-09-15T17:05:00Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:24:00Z
    by: agent
  - to: done
    at: 2026-09-15T17:38:00Z
    by: agent
stream: S-004
tags: [cli]
---

# T-010 Config read and write

## Work
internal/config: Load creates ~/.flai/config.json with defaults on first run, Save writes atomically, path resolution via --config, FLAI_CONFIG, then default. flai config get [key], flai config set key value, flai config path.

## Done when
Round-trip test passes; get and set operate on dotted keys for every field in flai-cli.md.

## Notes

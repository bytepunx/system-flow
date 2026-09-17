---
id: T-0010
type: task
nature: feature
title: Config read and write
status: done
parent: S-0004
owner: agent
created: 2026-09-15T16:24:42Z
updated: 2026-09-15T16:31:14Z
transitions:
  - to: ready
    at: 2026-09-15T16:24:42Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:29:56Z
    by: agent
  - to: done
    at: 2026-09-15T16:31:14Z
    by: agent
stream: S-0004
tags: [cli]
---

# T-0010 Config read and write

## Work
internal/config: Load creates ~/.flai/config.json with defaults on first run, Save writes atomically, path resolution via --config, FLAI_CONFIG, then default. flai config get [key], flai config set key value, flai config path.

## Done when
Round-trip test passes; get and set operate on dotted keys for every field in flai-cli.md.

## Notes

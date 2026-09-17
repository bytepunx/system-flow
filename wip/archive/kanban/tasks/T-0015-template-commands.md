---
id: T-0015
type: task
nature: feature
title: flai template show, update, use
status: done
parent: S-0005
owner: agent
created: 2026-09-15T16:36:28Z
updated: 2026-09-15T16:38:47Z
transitions:
  - to: ready
    at: 2026-09-15T16:36:28Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:38:47Z
    by: agent
  - to: done
    at: 2026-09-15T16:38:47Z
    by: agent
stream: S-0005
tags: [cli, template]
---

# T-0015 flai template show, update, use

## Work
template show prints the resolved source, cache path, and manifest summary; template update re-clones; template use <repo> [--ref] writes config.

## Done when
Each command has a test.

## Notes

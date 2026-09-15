---
id: T-009
type: task
nature: feature
title: flai version
status: done
parent: S-004
owner: agent
created: 2026-09-15T17:05:00Z
updated: 2026-09-15T17:24:00Z
transitions:
  - to: ready
    at: 2026-09-15T17:05:00Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:20:00Z
    by: agent
  - to: done
    at: 2026-09-15T17:24:00Z
    by: agent
stream: S-004
tags: [cli]
---

# T-009 flai version

## Work
Version, commit, and build date set via ldflags with dev defaults; plain and --json output.

## Done when
flai version prints all three; test covers both formats.

## Notes

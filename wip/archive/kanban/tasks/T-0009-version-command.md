---
id: T-0009
type: task
nature: feature
title: flai version
status: done
parent: S-0004
owner: agent
created: 2026-09-15T16:24:42Z
updated: 2026-09-15T16:29:56Z
transitions:
  - to: ready
    at: 2026-09-15T16:24:42Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:28:37Z
    by: agent
  - to: done
    at: 2026-09-15T16:29:56Z
    by: agent
stream: S-0004
tags: [cli]
---

# T-0009 flai version

## Work
Version, commit, and build date set via ldflags with dev defaults; plain and --json output.

## Done when
flai version prints all three; test covers both formats.

## Notes

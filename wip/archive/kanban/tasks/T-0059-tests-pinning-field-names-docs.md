---
id: T-0059
type: task
nature: feature
title: Tests pinning field names, docs
status: done
parent: S-0031
owner: alex
created: 2026-09-16T23:45:51Z
updated: 2026-09-16T23:48:50Z
transitions:
  - to: ready
    at: 2026-09-16T23:48:50Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:48:50Z
    by: agent
  - to: done
    at: 2026-09-16T23:48:50Z
    by: agent
stream: S-0031
tags: [cli, docs]
---

# T-0059 Tests pinning field names, docs

## Work
Tests pin the field names of the documented events (config created, template fetch, command failed); docs/users/flai.md documents --verbose, LOG_LEVEL, LOG_FORMAT; flai-cli.md notes the logger; design/tech/observability.md unchanged unless the choice moved.

## Done when
Docs match behaviour; all tiers green.

## Notes

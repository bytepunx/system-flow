---
id: T-0070
type: task
nature: feature
title: flai release and flai accept commands
status: done
parent: S-0029
owner: alex
created: 2026-09-17T02:04:19Z
updated: 2026-09-17T02:09:09Z
transitions:
  - to: ready
    at: 2026-09-17T02:09:09Z
    by: agent
  - to: in-progress
    at: 2026-09-17T02:09:09Z
    by: agent
  - to: done
    at: 2026-09-17T02:09:09Z
    by: agent
stream: S-0029
tags: [cli, release]
---

# T-0070 flai release and flai accept commands

## Work
cmd/release.go: flai release <id> [--dry-run] [--deliver] [--apply]; cmd/accept.go: flai accept <id> [--by] [--deliver] [--no-release] [--no-push] [--dry-run] [--trailer] does move to done (epics via review), archive, template version bump, commit, tags, push; refuses a dirty tree unless --yes; flai move ... done logs a hint.

## Done when
In-process test accepts a story end to end on a synthetic git repo.

## Notes

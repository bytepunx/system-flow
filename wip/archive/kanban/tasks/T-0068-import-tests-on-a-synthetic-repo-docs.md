---
id: T-0068
type: task
nature: feature
title: Import tests on a synthetic repo, docs
status: done
parent: S-0006
owner: alex
created: 2026-09-17T00:15:09Z
updated: 2026-09-17T00:19:51Z
transitions:
  - to: ready
    at: 2026-09-17T00:19:50Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:19:50Z
    by: agent
  - to: done
    at: 2026-09-17T00:19:51Z
    by: agent
stream: S-0006
tags: [cli, docs]
---

# T-0068 Import tests on a synthetic repo, docs

## Work
Tests build a synthetic repo with README, docs/, adr/, a loose note, a Go service and a SvelteKit app; docs/users/flai.md gains a Convert an existing repository section; flai-cli.md import flow matches the behaviour.

## Done when
All tiers green; docs match.

## Notes

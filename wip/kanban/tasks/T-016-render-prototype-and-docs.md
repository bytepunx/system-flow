---
id: T-016
type: task
nature: feature
title: Render ./template end to end and update docs
status: done
parent: S-005
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
stream: S-005
tags: [cli, template]
---

# T-016 Render ./template end to end and update docs

## Work
Integration test renders the real ./template with defaults and checks the expected tree, front matter validity, and no leftover template syntax; add internal/manifest for reading system-flow.yaml; update docs/users/flai.md and design docs with anything learned.

## Done when
Test passes against ./template; docs updated; S-002's deferred criterion confirmed.

## Notes

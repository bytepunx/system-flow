---
id: T-0670
type: task
nature: feature
title: "Design, operator guide, and release rule describe the version check and the minimum; I-0049 and I-0051 closed"
status: done
parent: S-0181
owner: arobson
created: 2026-10-01T10:43:14Z
updated: 2026-10-01T11:02:59Z
transitions:
  - to: ready
    at: 2026-10-01T10:53:42Z
    by: agent-S-0181
  - to: in-progress
    at: 2026-10-01T10:53:42Z
    by: agent-S-0181
  - to: done
    at: 2026-10-01T11:02:59Z
    by: agent-S-0181
stream: S-0181
tags: []
usage:
  source: log
  seconds: 557
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 185
      output: 52551
      cache_read: 12879167
      cache_write: 219744
      cost: 5.0746
---

# T-0670 Design, operator guide, and release rule describe the version check and the minimum; I-0049 and I-0051 closed

## Work

Describe the version check, the minimum, and the tolerant reading in `design/system/flai-cli.md`, `design/system/work-hierarchy.md`, `design/system/project-manifest.md`, the operator settings index and guide, and `docs/users/flai.md`. Add the rule that a release adding a front-matter field raises the minimum. Close I-0049 and I-0051 with what fixed them.

## Done when

- The documents above describe the behaviour, and `flai check --strict` and the docs lint pass
- I-0049 and I-0051 are closed with the fix named

## Notes

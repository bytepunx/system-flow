---
id: T-0618
type: task
nature: remediation
title: Design and user guide describe the remote release-tag check
status: done
parent: S-0174
owner: arobson
created: 2026-10-01T07:51:04Z
updated: 2026-10-01T08:10:14Z
transitions:
  - to: ready
    at: 2026-10-01T08:08:01Z
    by: agent-S-0174
  - to: in-progress
    at: 2026-10-01T08:08:01Z
    by: agent-S-0174
  - to: done
    at: 2026-10-01T08:10:14Z
    by: agent-S-0174
stream: S-0174
tags: []
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/]
usage:
  source: log
  seconds: 133
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 11272
      cache_read: 2982787
      cache_write: 32855
      cost: 1.085
---
# T-0618 Design and user guide describe the remote release-tag check

## Work

- Describe the check in `design/system/flai-cli.md` (release --pending, publish.preview, publish.run) and `design/system/flaiover-dashboard.md` (banner, done lane), and in the user guide for flai and the dashboard.

## Done when

- The documents say what is checked, when, and what the operator does; `flai check --strict` passes.

## Notes

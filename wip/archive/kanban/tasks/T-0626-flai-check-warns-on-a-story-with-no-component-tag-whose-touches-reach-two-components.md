---
id: T-0626
type: task
nature: remediation
title: flai check warns on a story with no component tag whose touches reach two components
status: done
parent: S-0174
owner: arobson
created: 2026-10-01T08:00:54Z
updated: 2026-10-01T08:04:56Z
transitions:
  - to: ready
    at: 2026-10-01T08:03:31Z
    by: agent-S-0174
  - to: in-progress
    at: 2026-10-01T08:03:31Z
    by: agent-S-0174
  - to: done
    at: 2026-10-01T08:04:56Z
    by: agent-S-0174
stream: S-0174
tags: []
touches: [flai/internal/check, design/issues]
usage:
  source: log
  seconds: 85
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 10898
      cache_read: 2883797
      cache_write: 31764
      cost: 1.049
---
# T-0626 flai check warns on a story with no component tag whose touches reach two components

## Work

- A `story.component-tag` warning on an open feature, improvement, or remediation story whose touches (and its open tasks') reach two or more components while no tag of its own or its epic's names one of them, saying which tags would do.
- Close I-0024 with what fixed it.

## Done when

- Check tests cover a story warned, one tagged, one whose epic is tagged, and one touching a single component; `flai check --strict` on this repository passes or its new findings are dealt with.

## Notes

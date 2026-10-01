---
id: T-0647
type: task
nature: remediation
title: The design and guides describe recorded variables, and I-0040 and I-0041 are closed
status: done
parent: S-0185
owner: arobson
created: 2026-10-01T08:49:42Z
updated: 2026-10-01T08:57:22Z
transitions:
  - to: ready
    at: 2026-10-01T08:49:50Z
    by: agent-S-0185
  - to: in-progress
    at: 2026-10-01T08:57:04Z
    by: agent-S-0185
  - to: done
    at: 2026-10-01T08:57:22Z
    by: agent-S-0185
stream: S-0185
tags: []
touches: [design/system/template.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/contributors/template.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues]
usage:
  source: log
  seconds: 18
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 74
      cache_read: 1013937
      cache_write: 8267
      cost: 0.4043
---
# T-0647 The design and guides describe recorded variables, and I-0040 and I-0041 are closed

## Work

- `design/system/template.md` and `design/system/project-manifest.md` say what the lock records and how upgrade resolves variables; `design/system/flai-cli.md` names `flai upgrade --var`.
- `docs/contributors/template.md` loses the I-0040 and I-0041 limitations and says how a fork's variable survives an upgrade; `docs/users/flai.md` shows `--var`; `docs/users/flai-reference.md` regenerated.
- Close I-0040 and I-0041 with `flai issue close`, saying what fixed them.

## Done when

- No document says upgrade renders with five variables only, or that an empty `--var` passes.
- `flai check --strict` passes.

## Notes

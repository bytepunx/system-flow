---
id: T-0644
type: task
nature: remediation
title: "Design and contributor guide say what an agent's environment holds; close I-0001, I-0039, I-0044"
status: done
parent: S-0183
owner: arobson
created: 2026-10-01T08:34:18Z
updated: 2026-10-01T08:40:28Z
transitions:
  - to: ready
    at: 2026-10-01T08:34:25Z
    by: agent-S-0183
  - to: in-progress
    at: 2026-10-01T08:38:42Z
    by: agent-S-0183
  - to: done
    at: 2026-10-01T08:40:28Z
    by: agent-S-0183
stream: S-0183
tags: []
touches: [design/system/flai-cli.md, design/system/devex.md, docs/contributors, design/issues]
usage:
  source: log
  seconds: 106
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 49
      output: 12207
      cache_read: 2546798
      cache_write: 45344
      cost: 1.0791
---
# T-0644 Design and contributor guide say what an agent's environment holds; close I-0001, I-0039, I-0044

## Work
Describe in design/system/flai-cli.md, design/system/devex.md, and the contributor guide what flai serve gives an agent, and what the scripts put on PATH. Close the three issues with what fixed them.

## Done when
The documents say it; flai check --strict passes; the issues are closed.

## Notes

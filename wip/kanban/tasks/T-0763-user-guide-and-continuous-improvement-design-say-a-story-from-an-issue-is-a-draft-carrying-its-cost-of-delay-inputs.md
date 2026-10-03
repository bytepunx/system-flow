---
id: T-0763
type: task
nature: improvement
title: User guide and continuous-improvement design say a story from an issue is a draft carrying its cost of delay inputs
status: done
parent: S-0203
owner: alex
created: 2026-10-03T18:07:46Z
updated: 2026-10-03T18:24:18Z
transitions:
  - to: ready
    at: 2026-10-03T18:08:34Z
    by: agent-S-0203
  - to: in-progress
    at: 2026-10-03T18:21:22Z
    by: agent-S-0203
  - to: done
    at: 2026-10-03T18:24:18Z
    by: agent-S-0203
stream: S-0203
tags: []
touches: [docs/users, design/system/continuous-improvement.md, design/system/flai-cli.md]
after: [T-0762]
usage:
  source: log
  seconds: 176
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 67
      output: 25591
      cache_read: 2772755
      cache_write: 99837
      cost: 1.6792
---

# T-0763 User guide and continuous-improvement design say a story from an issue is a draft carrying its cost of delay inputs

## Work

`docs/users/flai.md` (Record recurring friction, and Drafts, cost of delay, and forecasts), `design/system/continuous-improvement.md`, and `design/system/flai-cli.md` where it describes `flai issue story` say the story is a draft, carries the issue's cost of delay inputs as derived, reads an `## Impact` section's figures, and is named in the issue's Remediation section. Regenerate `docs/users/flai-reference.md` with `make flai-reference`. Waits for T-0762, whose help text it reflects.

## Done when

- The documents say so, with `updated` bumped, and the generated reference matches the help text.

## Notes

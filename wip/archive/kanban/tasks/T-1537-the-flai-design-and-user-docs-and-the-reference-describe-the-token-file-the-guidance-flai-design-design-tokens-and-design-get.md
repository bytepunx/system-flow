---
id: T-1537
type: task
nature: feature
title: The flai design and user docs and the reference describe the token file, the guidance, flai design, design.tokens, and design.get
status: done
parent: S-0381
owner: alex
created: 2026-10-13T17:46:50Z
updated: 2026-10-15T00:59:04Z
transitions:
  - to: ready
    at: 2026-10-15T00:55:16Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-15T00:55:17Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:59:04Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1534, T-1535, T-1536, T-1591]
usage:
  source: log
  seconds: 227
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 56
      output: 29389
      cache_read: 2760584
      cache_write: 133460
      cost: 1.8731
---
# T-1537 The flai design and user docs and the reference describe the token file, the guidance, flai design, design.tokens, and design.get

## Work

Describe in `design/system/flai-cli.md` the token file: DTCG 2025.10, the subset flai reads, its layout with the dark values where T-1592 put them, and the guidance beside it. Describe there too the `designsystem` package, `flai design init`, `show`, and `css`, the contrast pairs and thresholds, the rule `design.tokens`, and the host method `design.get`, linking ADR-0146 and ADR-0148. In `docs/users/flai.md`, say how a user starts a design system, checks it, and reads its contrast. Regenerate `docs/users/flai-reference.md` with `make flai-reference`. It waits for T-1534, T-1535, and T-1536, whose behaviour it describes, and T-1591, whose exception it links.

## Done when

- The three documents describe what the story built, with the tokens in the token file, not front matter.
- The markdown lint passes on them.

## Notes

- Layer 5. Rewritten by planner-S-0381 from the front matter plan, as TH-0453 approved.

---
id: T-0623
type: task
nature: improvement
title: The design and user guide describe what agents delegate and what sub-agents may do
status: done
parent: S-0175
owner: arobson
created: 2026-10-01T07:53:02Z
updated: 2026-10-01T08:12:57Z
transitions:
  - to: ready
    at: 2026-10-01T07:53:25Z
    by: agent-S-0175
  - to: in-progress
    at: 2026-10-01T08:09:06Z
    by: agent-S-0175
  - to: done
    at: 2026-10-01T08:12:57Z
    by: agent-S-0175
stream: S-0175
tags: []
touches: [design/system/flai-cli.md, docs/users]
usage:
  source: log
  seconds: 231
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 63
      output: 24601
      cache_read: 5866823
      cache_write: 74516
      cost: 2.1789
---
# T-0623 The design and user guide describe what agents delegate and what sub-agents may do

## Work

- `design/system/flai-cli.md` (prime, serve's prompt), `design/system/conventions.md` (roles, the new file), `docs/users/flai.md`, `docs/users/conventions.md`, and the generated reference.

## Done when

- Each describes delegation and the sub-agent roles, and `make flai-reference` leaves no diff.

## Notes

---
id: T-0716
type: task
nature: improvement
title: The baseline's roles name every agent that reads each convention
status: done
parent: S-0196
owner: arobson
created: 2026-10-02T23:29:26Z
updated: 2026-10-02T23:42:08Z
transitions:
  - to: ready
    at: 2026-10-02T23:29:45Z
    by: agent-S-0196
  - to: in-progress
    at: 2026-10-02T23:41:37Z
    by: agent-S-0196
  - to: done
    at: 2026-10-02T23:42:08Z
    by: agent-S-0196
stream: S-0196
tags: []
touches: [design/conventions/, template/root/design/conventions, template/CHANGELOG.md]
after: [T-0714]
usage:
  source: log
  seconds: 31
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2873
      cache_read: 532649
      cache_write: 12151
      cost: 0.251
---
# T-0716 The baseline's roles name every agent that reads each convention

## Work

Set `roles` in each file of `design/conventions/` and `template/root/design/conventions/` as the designer chose on 2026-10-02 (S-0196's third criterion, with TH-0071 settling `communication` and `tooling`). `code-quality` gets [story, verify]. `continuous-improvement`, `decisions`, `git`, `session-start`, and `work-management` get [story]. `delegation` gets [story, explore, verify] and `logging` [story, verify]. `documentation`, `safety`, `telemetry`, and `README` carry no roles. Bump each changed file's `updated`, and add a `template/CHANGELOG.md` entry. Waits for T-0714: until it lands, `flai check --strict` warns on `story`.

## Done when

- Both folders carry the same roles, matching the criterion and TH-0071
- `flai check --strict` reports nothing about `conventions.roles`

## Notes

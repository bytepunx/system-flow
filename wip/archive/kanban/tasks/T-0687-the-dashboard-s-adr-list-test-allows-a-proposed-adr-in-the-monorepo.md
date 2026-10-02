---
id: T-0687
type: task
nature: feature
title: The dashboard's ADR list test allows a proposed ADR in the monorepo
status: done
parent: S-0194
owner: arobson
created: 2026-10-02T10:21:02Z
updated: 2026-10-02T10:21:17Z
transitions:
  - to: ready
    at: 2026-10-02T10:21:09Z
    by: agent-S-0194
  - to: in-progress
    at: 2026-10-02T10:21:09Z
    by: agent-S-0194
  - to: done
    at: 2026-10-02T10:21:17Z
    by: agent-S-0194
stream: S-0194
tags: []
usage:
  source: log
  seconds: 8
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1385
      cache_read: 400701
      cache_write: 4636
      cost: 0.1449
---

# T-0687 The dashboard's ADR list test allows a proposed ADR in the monorepo

## Work

`flaiover/src/lib/server/search.test.ts` "lists ADRs with the supersession chain" reads this repository's ADRs and requires every one to be accepted, so `main` went red when 2d4e090 indexed ADR-0067 as `proposed`, and S-0194's close-out stops on it. The test asks instead that every status is one the ADR template allows (proposed, accepted, superseded, deprecated) and that ADR-0001 is accepted.

## Done when

- The test passes with a proposed ADR in the monorepo and still fails on a status outside the template's list.

## Notes

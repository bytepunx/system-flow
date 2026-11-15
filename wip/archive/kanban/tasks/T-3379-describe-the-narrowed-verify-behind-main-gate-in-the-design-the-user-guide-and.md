---
id: T-3379
type: task
nature: remediation
title: Describe the narrowed verify_behind_main gate in the design, the user guide, and flai accept's help
status: done
parent: S-0792
owner: alex
created: 2026-11-15T07:43:41Z
updated: 2026-11-15T17:10:48Z
transitions:
  - to: ready
    at: 2026-11-15T17:05:17Z
    by: agent-S-0792
  - to: in-progress
    at: 2026-11-15T17:05:17Z
    by: agent-S-0792
  - to: done
    at: 2026-11-15T17:10:48Z
    by: agent-S-0792
stream: S-0792
tags: [docs, acceptance]
touches: [design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, flai/cmd/accept.go, docs/users/flai-reference.md]
after: [T-3378]
usage:
  source: log
  seconds: 331
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 6452
      cache_read: 1558504
      cache_write: 70389
      cost: 0.9445
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 21939
          cache_write: 602
          cost: 0.0068
---
# T-3379 Describe the narrowed verify_behind_main gate in the design, the user guide, and flai accept's help

## Work

Rewrite the gate's paragraph in `design/system/workflow.md` § What waits on an acceptance (S-0296), its lines in `design/system/flai-cli.md` (the acceptance gate and the `stale_tiers` output), and `docs/users/flai.md`'s description of `verify_behind_main` to say that only main's paths overlapping the story's change set send back tiers, linking T-3376's ADR. Change the help text in `flai/cmd/accept.go` that names the gate, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Record on I-0378, with `flai issue`, that S-0792 remediates it: a main commit outside the story's change set, such as the digest fix it records, no longer sends the story back.

It waits for T-3378, so the docs describe the gate and the message as built.

## Done when

- The three documents and the help describe the overlap rule and link the ADR, and none still says that any main path selecting a tier sends it back.
- `TestReferenceIsCurrent` passes, and `flai check --strict` and the markdown lint are clean.

## Notes

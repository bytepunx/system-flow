---
id: T-1920
type: task
nature: improvement
title: "An ADR refines ADR-0142: task done refuses only when the newest flai release's task-done rules number is above the running flai's"
status: done
parent: S-0445
owner: alex
created: 2026-10-20T20:41:42Z
updated: 2026-10-20T20:45:34Z
transitions:
  - to: ready
    at: 2026-10-20T20:44:55Z
    by: agent-S-0445
  - to: in-progress
    at: 2026-10-20T20:44:56Z
    by: agent-S-0445
  - to: done
    at: 2026-10-20T20:45:34Z
    by: agent-S-0445
stream: S-0445
tags: [flai, adr]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 38
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2449
      cache_read: 679644
      cache_write: 17845
      cost: 0.3154
---
# T-1920 An ADR refines ADR-0142: task done refuses only when the newest flai release's task-done rules number is above the running flai's

## Work

Write the decision TH-0506 proposed (option B) as a new ADR with `flai adr new`, refining ADR-0142 and ADR-0160:

- flai carries a task-done rules number, `buildinfo.TaskDoneRules`, bumped by any change to task done's commit, sync, move, log, touches, or check steps; `flai version --json` prints it.
- `flai release` writes `task-done-rules: N` into the annotation of each `flai/vX.Y.Z` tag it makes.
- `flai task done`, and ADR-0160's caught-up check of the installed flai, refuse only when the newest release's rules number is above the running (or installed) flai's. A tag without the line, or a flai that prints no number, keeps ADR-0142's minor-version rule.
- Context: I-0157's instances, ten minor releases in about 28 hours with `auto-upgrade` off, and 14 open upgrade threads.
- Alternatives: auto-upgrade on (A), resolving superseded upgrade threads (C), and the orchestrator batching releases, with why each leaves the cause in place.

This task waits for nothing; every other task builds what it decides. If the operator answers TH-0506 with A or C instead, rewrite this task and drop the code tasks.

## Done when

- The ADR exists, is linked from `design/adrs/README.md`, and names ADR-0142 and ADR-0160 as refined.
- It states the rules number, where a release records it, the comparison, and the fallback for an untagged release.
- The operator has accepted it, or the story's narrative says it waits for them.

## Notes

The ADR's own file is named by `flai adr new`; `flai task done` widens the touches with it.

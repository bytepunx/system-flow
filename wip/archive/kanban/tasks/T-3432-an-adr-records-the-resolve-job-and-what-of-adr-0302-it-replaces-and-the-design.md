---
id: T-3432
type: task
nature: feature
title: An ADR records the resolve job and what of ADR-0302 it replaces, and the design and the user guide describe it
status: done
parent: S-0804
owner: alex
created: 2026-11-16T23:19:55Z
updated: 2026-11-17T01:05:54Z
transitions:
  - to: ready
    at: 2026-11-16T23:20:39Z
    by: agent-S-0804
  - to: in-progress
    at: 2026-11-17T00:52:30Z
    by: agent-S-0804
  - to: done
    at: 2026-11-17T01:05:54Z
    by: agent-S-0804
stream: S-0804
tags: []
touches: [design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flaiover.md, design/adrs/README.md, design/adrs/0302-the-operator-sends-a-review-story-that-cannot-be-accepted-back-to-its-agent-to.md, design/adrs/0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md]
after: [T-3430]
usage:
  source: log
  seconds: 804
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 77
      output: 32435
      cache_read: 5851566
      cache_write: 219878
      cost: 3.046
---
# T-3432 An ADR records the resolve job and what of ADR-0302 it replaces, and the design and the user guide describe it

## Work

Waits for T-3430, so that what it describes is built; runs beside the dashboard task, with no path in common.

1. A new ADR with `flai adr new` (no `--commit`: the story's agent commits it), superseding ADR-0302: `flai accept --resolve` and `accept.resolve` start a resolve job in flai serve, answered at once; the job rebases with `flai stream sync`, closes out at the rebased head with its output in the job's log and its steps and tiers reported as they start and end (`FLAI_VERIFY_PROGRESS`), and sends the story back through ADR-0248's send-back thread only on conflicts, which it aborts, or a failed step; a clean pass leaves the story in review; a second request answers the job under way; an interrupted job sends nothing back; the dashboard follows the job through `resolve.status` and `log.read`. Say what of ADR-0302 stays (the refusals, the operator's alone, no host action) and what goes (the send-back on every request). Topics: flai, cli, dashboard, acceptance, verify.
2. `design/system/flai-cli.md`: the `flai accept <id> --resolve` row, and the close-out job section (the kind, the progress lines, `Status`), linking the ADR.
3. `design/system/workflow.md`: where it names Rebase and resolve and the send back from review.
4. `docs/users/flai.md` (Accept and release) and `docs/users/flaiover.md` (Reviewing a story): what the operator sees and what happens next.

## Done when

- The markdown lint passes on the changed files (`flai test` on them).
- No document still says Rebase and resolve sends every story back.

## Notes

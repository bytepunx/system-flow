---
id: T-3051
type: task
nature: remediation
title: "ADR: an acceptance waits for a verify over main's new paths"
status: done
parent: S-0683
owner: alex
created: 2026-11-10T10:35:04Z
updated: 2026-11-12T15:53:59Z
transitions:
  - to: ready
    at: 2026-11-12T15:52:56Z
    by: agent-S-0683
  - to: in-progress
    at: 2026-11-12T15:52:57Z
    by: agent-S-0683
  - to: done
    at: 2026-11-12T15:53:59Z
    by: agent-S-0683
stream: S-0683
tags: [adr, acceptance]
touches: [design/adrs/0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md, design/adrs/README.md]
usage:
  source: log
  seconds: 61
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 46
      cache_read: 516756
      cache_write: 2168
      cost: 0.193
---
# T-3051 ADR: an acceptance waits for a verify over main's new paths

## Work

Write the decision that remediates I-0229 with `flai adr new "An acceptance waits for a verify over main's new paths"`, and accept it once the plan thread on S-0683 agrees to it. The number is the next free one; 0279 is a guess, and `flai task done` widens the touches to the real file.

Context: S-0416's close-out passed at 6c25aad7 over base 1dfa57c8. S-0535 then landed and reworded the to-backlog condition's `Allows` in `flai/internal/roles/conditions.go`. The operator accepted S-0416 by hand, and `flai accept` rebased it onto that main and fast-forwarded. No tier ran at the new head, and main went red at `flai/internal/guard/roles_test.go`. The merge queue then failed S-0519, which had not caused the break. Acceptance runs no tier. `queueTiersGate` (ADR-0221) compares the last report only to the branch's own head, so it cannot see how far main has moved.

Decision to record: every acceptance but the queue's landing (`--batch`, whose batch verified the rebased head, ADR-0176) is refused before anything is merged when both hold. Main has brought in commits since the base the story's last passing `flai verify` ran over. Those commits change paths outside `wip/` that send back a tier that verify ran, as ADR-0156's `sentBack` selects them, `{importers}` and covering tiers included. The blocker names the tiers and main's paths that reach them. It also names the close-out to ask for again (`flai verify <id> --job`, or the MCP tool `verify`), which re-runs only those tiers and reuses the rest. This covers the operator's, the board's, and the orchestrator's acceptance alike, with the queue on or off. Name the alternatives you weighed: running the tiers inside `flai accept`, an override flag, and only fixing the test. Say why each was rejected or kept. This task waits for no task.

## Done when

- The ADR is written, its status matches the plan thread's answer, and `design/adrs/README.md` lists it (`flai check` reports no `adr.index`).
- The ADR names I-0229, ADR-0156, ADR-0176, and ADR-0221.

## Notes

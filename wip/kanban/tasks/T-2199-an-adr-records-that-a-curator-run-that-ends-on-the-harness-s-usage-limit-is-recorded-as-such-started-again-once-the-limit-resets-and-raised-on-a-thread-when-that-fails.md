---
id: T-2199
type: task
nature: remediation
title: An ADR records that a curator run that ends on the harness's usage limit is recorded as such, started again once the limit resets, and raised on a thread when that fails
status: done
parent: S-0514
owner: alex
created: 2026-10-28T12:52:56Z
updated: 2026-10-29T01:02:14Z
transitions:
  - to: ready
    at: 2026-10-29T00:58:07Z
    by: agent-S-0514
  - to: in-progress
    at: 2026-10-29T00:58:09Z
    by: agent-S-0514
  - to: done
    at: 2026-10-29T01:02:14Z
    by: agent-S-0514
stream: S-0514
tags: [design, adr, curator]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 232
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 36
      output: 197
      cache_read: 1136698
      cache_write: 73381
      cost: 0.4119
---
# T-2199 An ADR records that a curator run that ends on the harness's usage limit is recorded as such, started again once the limit resets, and raised on a thread when that fails

## Work

Write the ADR with `flai adr new` (or the MCP tool `adr_new`), refining ADR-0199, which says a start the checks refuse is not tried again, and ADR-0196. Its decision, as the plan on S-0514's thread proposes it:

- A curator run whose log ends on the harness's usage limit (the usage reader's mark, from Claude Code's `rate_limit_event` and `result` records, with the reset time) is recorded in `serve/agents.json` as `failed` with `why` saying it ended on the usage limit and when the limit resets, and with the reset time as a field of its own, whether or not its exit was seen. It logs no activity: it spent nothing, and `since` keeps pointing at the last activity.
- While the curator's last run ended on the limit and the reset time has not passed, the curate scheduler starts nothing: the acceptances and the schedule times that come meanwhile are held. At the first look once the limit has reset, it starts one run, with the trigger of the run it retries and the held acceptances joined, marked as a retry on its record, through the same checks as any start. A `flai curate` the operator asked for is retried the same way. The state is on the run in `serve/agents.json`, so a `flai serve` restarted meanwhile still retries.
- A retry that ends on the limit again, and a run whose log names the limit but no reset time, are raised to the operator once on a thread on `system-flow.yaml`, by `flai serve`, naming the run, the reset time, and what starts the curator again (`flai curate`, the next acceptance, or the schedule); the thread is recorded on the run so it is opened once. No further retry is scheduled for that limit.
- The analyzer, the planner, and the orchestrator keep their behaviour; extending the mark to them is a story of its own.

Record the alternatives considered: a thread alone with no retry, and a retry for every strategic agent. Set the ADR's topics with `flai adr topics` (planning, digest). `flai adr new` adds the row to `design/adrs/README.md`.

## Done when

- The ADR exists under `design/adrs/`, accepted, with context, decision, alternatives, and consequences, and the index row in `design/adrs/README.md` names it.
- `flai check --strict` and the markdown lint pass on the ADR and the index.

## Notes

Layer 1 of S-0514's plan; waits for nothing. The ADR's file is not in the touches, since `flai adr new` allocates its number: `flai task done` widens the touches with it once committed (ADR-0144). `design/adrs` is in `claims.shared`, so the overlap with S-0502's ADR holds nothing; if S-0502 lands first, the number this story gets is the next free one.

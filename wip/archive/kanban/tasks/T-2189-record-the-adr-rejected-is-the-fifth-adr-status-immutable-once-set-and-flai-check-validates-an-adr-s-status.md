---
id: T-2189
type: task
nature: remediation
title: "Record the ADR: rejected is the fifth ADR status, immutable once set, and flai check validates an ADR's status"
status: done
parent: S-0512
owner: alex
created: 2026-10-28T08:27:31Z
updated: 2026-10-30T14:10:47Z
transitions:
  - to: ready
    at: 2026-10-30T14:10:24Z
    by: agent-S-0512
  - to: in-progress
    at: 2026-10-30T14:10:25Z
    by: agent-S-0512
  - to: done
    at: 2026-10-30T14:10:47Z
    by: agent-S-0512
stream: S-0512
tags: [cli, template]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 22
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1205
      cache_read: 270461
      cache_write: 11882
      cost: 0.1584
---
# T-2189 Record the ADR: rejected is the fifth ADR status, immutable once set, and flai check validates an ADR's status

## Work

Record the decision with `flai adr new --commit` from the story's worktree, status `accepted`, refining ADR-0001 (markdown and ADRs). The decision sentence: an ADR's status is one of `proposed`, `accepted`, `superseded`, `deprecated`, or `rejected`; a rejected ADR is a proposed decision the operator declined, recorded in favour of another, and is immutable once set as an accepted one is; and `flai check` reports `adr.status` on an ADR whose status is outside the five.

Context: ADR-0202 was set `rejected` by hand on main (commit a554b18ff, in favour of ADR-0203), a status the template's comment, `documentation-standard.md`, `flai doc`'s refusal, and `search.test.ts` did not name; `flai check` accepted it without a finding, and every story's vitest tier failed until S-0501 added it to the test. The alternatives to name: keep four statuses and record a rejection as a superseding ADR (loses the record that nothing replaced the decision's shape, and ADR-0202 already carries `rejected`); validate nothing and let the test keep its own list (the cause stays). Consequences: `flai adr new --status` still takes only `proposed` and `accepted`; a rejection is written on a proposed ADR, which is editable, and the ADR is read-only from then on.

This task waits for nothing: the decision is settled in the plan, and the code tasks of the same layer build it.

## Done when

- The ADR exists under `design/adrs/`, status `accepted`, with its row in `design/adrs/README.md`, committed on the story branch.
- Its `## Decision` opens with the decision sentence, so that `flai check` reports no `adr.decision`.
- `flai check --strict` reports nothing new.

## Notes

Drafted by planner-S-0512. The ADR's file is named by `flai adr new`; `flai task done` widens the touches with it.

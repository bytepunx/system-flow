---
id: T-0713
type: task
nature: improvement
title: The design records what S-0230 adopted and how a task's usage is measured
status: done
parent: S-0230
owner: arobson
created: 2026-10-02T17:15:49Z
updated: 2026-10-02T17:26:29Z
transitions:
  - to: ready
    at: 2026-10-02T17:16:14Z
    by: agent-S-0230
  - to: in-progress
    at: 2026-10-02T17:22:41Z
    by: agent-S-0230
  - to: done
    at: 2026-10-02T17:26:29Z
    by: agent-S-0230
stream: S-0230
tags: []
touches: [design/adrs, design/system, design/issues, docs, flai/cmd]
after: [T-0711, T-0712]
usage:
  source: log
  seconds: 228
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 83
      output: 27535
      cache_read: 3298099
      cache_write: 107130
      cost: 1.8951
---
# T-0713 The design records what S-0230 adopted and how a task's usage is measured

## Work

Record what S-0230 adopts, in the living design and its ADRs:

- An ADR, `flai adr new ... --status accepted --refines 0051` (check the flag for refining), that a task's usage is the calls of the sub-agents started for it, named by its ID in the `Agent` call's description or prompt (`parent_tool_use_id`), plus its share of the story's agent's own calls while it was in progress, split evenly among the tasks in progress at once; as T-0711 built it. Change the design that describes ADR-0051's apportioning (search `design/system`, `design/tech`, and `docs/` for it: `metrics.md`, `flai-cli.md`, `flaiover-dashboard.md`, `docs/users/flai.md` are candidates) to say so, with a link to the new ADR.
- An ADR that a story's agent plans its tasks with `after` and hands each to a task sub-agent, and that running a layer's tasks at once is optional and worth it only where the tasks are long beside the story's fixed costs, following S-0176's results document (`design/experiments/S-0176-...md`) and ADR-0066.
- `design/system/agent-context.md` § Tasks in parallel: a subsection saying what S-0230 adopted (task `after` and the plan's display, the hand-off of tasks to sub-agents named by their task's ID, the layers optional), linking the results document and both ADRs; that the host's flai 1.28.0 carries `after` and `flai.minimum` is 1.28.0 (published in `e7e9709`); that I-0054 is fixed; and that on 2026-10-02 a `claude -p` session flai serve started (Claude Code 2.1.286) still refused `subagent_type: fork` ("Agent type 'fork' not found"), so the repeat measurement is S-0241 (TH-0069). Update the section's opening, which still calls S-0176 an experiment in the present, and the Recommendation's "Before any of it is adopted" paragraph only so far as they would now be untrue.
- `design/system/work-hierarchy.md`: where it says a host with an older flai refuses a task's `after`, say the minimum is 1.28.0.
- `design/issues/I-0054-*.md`: record the remediation and close it, and update `design/issues/summary.md`.

Bump `updated:` on every design file changed.

## Done when

- [x] Both ADRs exist, accepted, indexed, and the apportioning one refines ADR-0051
- [x] `agent-context.md` records what was adopted, links the results document, and records the fork probe and S-0241
- [x] No design or docs file still describes task usage as the window alone
- [x] I-0054 is closed with its remediation

## Notes

Waits for T-0711 and T-0712: it describes what they built, in their words.

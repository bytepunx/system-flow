---
id: T-0747
type: task
nature: improvement
title: A baseline strategic-agents.md convention, and the baseline's roles name the new roles
status: done
parent: S-0207
owner: alex
created: 2026-10-03T07:11:19Z
updated: 2026-10-03T07:20:59Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:43Z
    by: agent-S-0207
  - to: in-progress
    at: 2026-10-03T07:12:32Z
    by: agent-S-0207
  - to: done
    at: 2026-10-03T07:20:59Z
    by: agent-S-0207
stream: S-0207
tags: [template]
touches: [design/conventions, template/root/design/conventions, template/CHANGELOG.md, template/template.yaml]
usage:
  source: log
  seconds: 507
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 11923
      cache_read: 2326199
      cache_write: 47322
      cost: 1.0018
---
# T-0747 A baseline strategic-agents.md convention, and the baseline's roles name the new roles

## Work

- Write `strategic-agents.md`, here and in the template, under 120 lines, `roles: [plan, orchestrate, analyze]`, listed in both `README.md` indexes in its order: what the planner, the orchestrator, and the analyzer each do, what each never does (the planner never moves an item past backlog; the orchestrator never edits code; the analyzer never authors stories), how each logs its activity (its activity document, S-0206), and how each asks the operator (`thread_open` on the item, then `wait_for_events`).
- Set the baseline's `roles` to name the new roles where the designer confirms on TH-0081, here and in the template, and record the template change in `template/CHANGELOG.md` and `template/template.yaml`.
- Waits for nothing in code: the roles are only warnings in `flai check` until T-0746 lands, and the paths are disjoint, so it runs in the first layer with T-0746; it waits on TH-0081's answer for the roles.

## Done when

- `strategic-agents.md` exists in both places, identical above the marker, indexed, and lint clean.
- Every baseline file's `roles` matches TH-0081's answer in both places.

## Notes

---
id: T-2526
type: task
nature: improvement
title: Follow the analysis-cadence instance to review and acceptance
status: done
parent: S-0564
owner: alex
created: 2026-11-07T03:47:38Z
updated: 2026-11-07T08:37:30Z
transitions:
  - to: ready
    at: 2026-11-07T08:04:27Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:04:28Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:37:30Z
    by: agent-S-0564
stream: S-0564
tags: [template]
after: [T-2520]
usage:
  source: log
  seconds: 1672
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 117
      output: 13686
      cache_read: 6488662
      cache_write: 72551
      cost: 2.1523
---
# T-2526 Follow the analysis-cadence instance to review and acceptance

## Work

`flai serve` starts the instance's agent. That agent starts the analyzer step with `role_run` and holds `wait_for_events` until the analyzer's `agent` event on the instance (ADR-0245). Then it runs the listing step, `flai issue list --report latest --json`, in the main checkout. It records in the instance's Notes the report and the issues it names, each with its title and class. Follow it with `wait_for_events`, `flai board`, and `flai show <instance>`.

When a step fails, the instance's agent blocks its task and asks the operator on a thread on the instance: the analyzer's run failed or asked, or the listing found no report. Answering it is the operator's.

The instance is accepted by the operator, or by the orchestrator while `orchestration.permissions.accept_reviews` is on. Its acceptance adds no bump, since its branch changes nothing. Write nothing in S-0564's Notes yet: the record of both runs is a later task.

This task waits for T-2520, which makes the instance. It runs beside T-2523, with which it shares no path.

## Done when

- The instance is done, and its Notes name the analyzer's report and the issues the report names.

## Notes

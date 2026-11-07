---
id: T-2520
type: task
nature: improvement
title: Check the cadence run's preconditions and start one instance of the analysis-cadence formula from the main checkout
status: done
parent: S-0564
owner: alex
created: 2026-11-07T03:47:15Z
updated: 2026-11-07T08:03:54Z
transitions:
  - to: ready
    at: 2026-11-07T08:03:14Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:03:15Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:03:54Z
    by: agent-S-0564
stream: S-0564
tags: [template]
usage:
  source: log
  seconds: 39
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 340
      cache_read: 162644
      cache_write: 1641
      cost: 0.0525
---
# T-2520 Check the cadence run's preconditions and start one instance of the analysis-cadence formula from the main checkout

## Work

Check these on the host before anything runs. If one does not hold, ask the operator on a thread on S-0564, and run nothing until it is answered:

- The `analyze` host action is on for this project. The analyzer runs under it, from a formula's role step as from the schedule (ADR-0245).
- No analyzer runs now. The analyzer runs once per project at a time, and a role step asked while another runs is refused, so the step fails (ADR-0245).
- TH-0845, the plan's thread, names the focus on question 1: `bottlenecks`, `intent`, `risk`, or none for all three.

Then start one instance in the main checkout:

```bash
flai formula run analysis-cadence --input focus=<focus> --ready --autocommit
```

Leave out `--input` when the focus is none: the input's default is empty, which looks for all three. Write the instance's ID, the focus, and how it was started in this task's Notes with `flai edit`. Run no second instance.

This task waits for no other task. It is in the plan's first layer, beside the import's start, with which it shares no path.

## Done when

- Each precondition above is checked, and the result is written in this task's Notes.
- One instance of the `analysis-cadence` formula exists in this project, in ready or later, and this task's Notes name its ID and its focus.

## Notes

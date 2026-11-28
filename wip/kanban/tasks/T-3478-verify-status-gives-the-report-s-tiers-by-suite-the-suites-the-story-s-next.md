---
id: T-3478
type: task
nature: feature
title: verify.status gives the report's tiers by suite, the suites the story's next transitions need, and which of them have not passed at its head
status: in-progress
parent: S-0814
owner: alex
created: 2026-11-16T23:50:58Z
updated: 2026-11-28T00:43:59Z
transitions:
  - to: ready
    at: 2026-11-28T00:43:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-28T00:43:59Z
    by: system-flow
stream: S-0814
tags: [flai]
touches: [flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/verify_suites.go, flai/internal/hostapi/verify_suites_test.go]
---
# T-3478 verify.status gives the report's tiers by suite, the suites the story's next transitions need, and which of them have not passed at its head

## Work

Criterion 1, the host API's part. `verifyStatus` in `writes_host.go` answers the stored report. Add to its answer, in a new `verify_suites.go`:

- the report's tiers grouped by suite, as the report S-0807 writes records them;
- the suites the transitions out of the story's current state need, from their `verify:`, with each transition's to and severity, read through the function S-0807 adds for the close-out;
- for each needed suite, whether every tier of it the story's paths select passed at the branch's head, and, when not, why: not run, failed, or passed at another head.

A report from before S-0807, with no suites, groups every tier under its own name, as a tier with no `suite` is the suite of its own name.

Waits for no task: it reads S-0807's report and S-0806's workflow, both on main before this story starts, and touches no file the other first-layer tasks touch.

## Done when

- `verify.status` answers the grouped tiers, the needed suites, and their state at the head, for a story whose report passed every suite, one with a suite not run, one with a suite failed, and one whose head moved.
- A report with no suite field answers each tier as its own suite.
- `verify_suites_test.go` and `writes_test.go` cover these, and `flai test` on the changed paths passes.

## Notes

- Drafted by planner-S-0814. The report type is `verify.Report` in `flai/internal/verify/story.go`, whose tiers are steps with `Tier` set; S-0807 adds the suites.

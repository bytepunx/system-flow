---
id: S-0824
type: story
nature: improvement
title: A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted
status: backlog
owner: alex
created: 2026-11-17T07:13:20Z
updated: 2026-11-17T07:13:20Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: flai
    at: 2026-11-17T07:13:20Z
---
# S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted

## Goal

This story remediates [I-0398](../../../design/issues/I-0398-a-rebase-and-resolve-job-runs-the-whole-close-out-again-each-time-main-moves.md), "A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted". The issue recommends this solution:

Decide, in an ADR refining ADR-0313 and ADR-0156, what a rebase and resolve job owes when main moves under it, and build it. Options to weigh: the job passes at the base its last round ran on and leaves the rest to the acceptance's own `verify_behind_main` gate (ADR-0291, ADR-0303), which sends back only the tiers the new paths overlap; the later rounds apply ADR-0303's overlap rule in place of every tier main's paths select; or flai serve holds acceptances and resolve jobs in one order, so that main does not move under a job it is about to accept. ADR-0312's acceptance suites (S-0808) bear on the same cost. The job's record and the review page say which round the job is on and why it started another.

## Acceptance criteria
- [ ] The cause I-0398 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0398 is closed with `flai issue close I-0398 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0398. time_lost_per_cycle 1h carried over from I-0398's Impact section. Its Impact time_lost_per_cycle was taken rather than the 22m derived from its cost and count.

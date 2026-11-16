---
id: T-3391
type: task
nature: feature
title: flai guard and the role files' check judge a to condition against the project's workflow file
status: done
parent: S-0755
owner: alex
created: 2026-11-15T23:06:47Z
updated: 2026-11-15T23:06:57Z
transitions:
  - to: ready
    at: 2026-11-15T23:06:50Z
    by: agent-S-0755
  - to: in-progress
    at: 2026-11-15T23:06:50Z
    by: agent-S-0755
  - to: done
    at: 2026-11-15T23:06:57Z
    by: agent-S-0755
stream: S-0755
tags: []
touches: [flai/cmd/guard.go, flai/internal/guard/roles.go, flai/internal/check/check_roles.go, flai/internal/check/roles_test.go]
after: [T-3390]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3735
      cache_read: 770141
      cache_write: 22597
      cost: 0.3618
---
# T-3391 flai guard and the role files' check judge a to condition against the project's workflow file

## Work

Waits for T-3390, whose review found it. `flai guard` held a role to the baseline (`Guard.Workflow` was never set outside tests), and `flai check` validated a role file's `to` conditions against the baseline, so a role file naming a state of the project's own file was refused. The guard now takes `repo.Flow()` where it opens the project, and so does the check of role files (ADR-0307 §1).

## Done when

A check test shows a `to` condition on a state that only the project's file adds is a finding on the baseline and passes with the file; `flai test` passes on the changed paths.

## Notes
